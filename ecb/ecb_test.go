package ecb

import (
	"errors"
	"testing"
	"time"

	"github.com/filiptubic/eurex/currency"
	log "github.com/sirupsen/logrus"
)

func TestECBConverter_newRates(t *testing.T) {
	tt := []struct {
		name   string
		data   *ECBResponseData
		verify func(m *Rates, err error)
	}{
		{
			name: "ok map",
			data: &ECBResponseData{
				Data: []DataXML{
					{Date: DateXML("2022-1-4"), Rates: []RateXML{{Currency: "USD", Rate: 1.5}}},
					{Date: DateXML("2022-1-3"), Rates: []RateXML{{Currency: "TRY", Rate: 2}}},
				},
			},
			verify: func(rates *Rates, err error) {
				if err != nil {
					t.Error(err)
				}

				if len(rates.rates) != 2 {
					t.Errorf("expected only two entry, got: %d", len(rates.rates))
				}
				t1 := time.Date(2022, 1, 4, 0, 0, 0, 0, time.Local)
				if _, ok := rates.rates[t1]; !ok {
					t.Errorf("missing %v time in map", t1)
				}
				t2 := time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local)
				if _, ok := rates.rates[t1]; !ok {
					t.Errorf("missing %v time in map", t2)
				}
			},
		},
		{
			name: "invalid date layout",
			data: &ECBResponseData{
				Data: []DataXML{
					{Date: DateXML("invalid date")},
				},
			},
			verify: func(m *Rates, err error) {
				if _, ok := err.(InvalidDateFormat); !ok {
					t.Errorf("expecting InvalidDateFormat got: %v", err)
				}
			},
		},
		{
			name: "invalid currency",
			data: &ECBResponseData{
				Data: []DataXML{
					{Date: DateXML("2022-12-12"), Rates: []RateXML{{Currency: "UNKNOWN"}}},
				},
			},
			verify: func(m *Rates, err error) {
				if _, ok := err.(InvalidCurrency); !ok {
					t.Errorf("expecting InvalidCurrency got: %v", err)
				}
			},
		},
		{
			name: "validate first and last date",
			data: &ECBResponseData{
				Data: []DataXML{
					{Date: DateXML("2022-5-5"), Rates: []RateXML{{Currency: "USD"}}},
					{Date: DateXML("2023-1-1"), Rates: []RateXML{{Currency: "USD"}}},
					{Date: DateXML("2022-4-3"), Rates: []RateXML{{Currency: "USD"}}},
					{Date: DateXML("2022-1-1"), Rates: []RateXML{{Currency: "USD"}}},
					{Date: DateXML("2022-5-3"), Rates: []RateXML{{Currency: "USD"}}},
				},
			},
			verify: func(m *Rates, err error) {
				if err != nil {
					t.Error(err)
				}
				if m.first != time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local) {
					t.Errorf("invalid first date in rates: %v", m.first)
				}
				if m.last != time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local) {
					t.Errorf("invalid last date in rates: %v", m.last)
				}
			},
		},
	}

	converter := New(&ECBClientMock{}, false, nil)
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			test.verify(converter.newRates(test.data))
		})
	}
}

func TestEcbConverter_GetRates(t *testing.T) {
	tt := []struct {
		name   string
		date   time.Time
		c      ECBConverter
		verify func(rates *Rates, err error)
	}{
		{
			name: "using cache without cached data",
			date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			c: ECBConverter{cache: true, cached: nil, client: &ECBClientMock{
				GetRatesMock: func() (*ECBResponseData, error) {
					return &ECBResponseData{
						Data: []DataXML{
							{
								Date:  DateXML("2022-1-1"),
								Rates: []RateXML{{Currency: "USD", Rate: 2}},
							},
						},
					}, nil
				},
			}},
			verify: func(rates *Rates, err error) {
				if err != nil {
					t.Error(err)
				}
				if len(rates.rates) != 1 {
					t.Errorf("expecting one rate got: %d", len(rates.rates))
				}
				if _, ok := rates.rates[time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)]; !ok {
					t.Errorf("missing rate at date '2022-1-1'")
				}
			},
		},
		{
			name: "using cache with cached data",
			date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			c: ECBConverter{
				cache:  true,
				logger: log.New(),
				cached: &Rates{
					rates: map[time.Time]currencyMap{
						time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local): {
							currency.USD: 1.5,
						},
					}},
			},
			verify: func(rates *Rates, err error) {
				if err != nil {
					t.Error(err)
				}
				if len(rates.rates) != 1 {
					t.Errorf("expecting one rate got: %d", len(rates.rates))
				}
				if _, ok := rates.rates[time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)]; !ok {
					t.Errorf("missing rate at date '2022-1-1'")
				}
			},
		},
		{
			name: "using cache with missing new data",
			date: time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
			c: ECBConverter{
				cache:  true,
				logger: log.New(),
				cached: &Rates{
					rates: map[time.Time]currencyMap{
						time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local): {
							currency.USD: 1.5,
						},
					},
				},
				client: &ECBClientMock{GetRatesMock: func() (*ECBResponseData, error) {
					return &ECBResponseData{
						Data: []DataXML{
							{
								Date:  DateXML("2022-1-1"),
								Rates: []RateXML{{Currency: "USD", Rate: 2}},
							},
							{
								Date:  DateXML("2022-1-2"),
								Rates: []RateXML{{Currency: "USD", Rate: 3}},
							},
						},
					}, nil
				}},
			},
			verify: func(rates *Rates, err error) {
				if err != nil {
					t.Error(err)
				}
				if len(rates.rates) != 2 {
					t.Errorf("expecting two rates got: %d", len(rates.rates))
				}
				if _, ok := rates.rates[time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local)]; !ok {
					t.Errorf("missing rate at date '2022-1-1'")
				}
			},
		},
		{
			name: "getting rates without caching",
			date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			c: ECBConverter{
				cache: false,
				client: &ECBClientMock{
					GetRatesMock: func() (*ECBResponseData, error) {
						return &ECBResponseData{
							Data: []DataXML{
								{
									Date:  DateXML("2022-1-1"),
									Rates: []RateXML{{Currency: "USD", Rate: 2}},
								},
							},
						}, nil
					},
				},
			},
			verify: func(rates *Rates, err error) {
				if err != nil {
					t.Error(err)
				}
				if len(rates.rates) != 1 {
					t.Errorf("expecting one rate got: %d", len(rates.rates))
				}
				if _, ok := rates.rates[time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)]; !ok {
					t.Errorf("missing rate at date '2022-1-1'")
				}
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			test.verify(test.c.GetRates(test.date))
		})
	}
}

func TestECBConverter_Convert(t *testing.T) {
	tt := []struct {
		name         string
		date         time.Time
		from, to     currency.Currency
		value        float64
		GetRatesMock func() (*ECBResponseData, error)
		verify       func(value float64, err error)
	}{
		{
			name: "invalid from currency",
			date: time.Now(),
			from: "UNKNOWN",
			verify: func(value float64, err error) {
				if _, ok := err.(InvalidCurrency); !ok {
					t.Errorf("expecting InvalidCurrency, got %v", err)
				}
			},
		},
		{
			name: "invalid to currency",
			date: time.Now(),
			from: currency.USD,
			to:   "UNKNOWN",
			verify: func(value float64, err error) {
				if _, ok := err.(InvalidCurrency); !ok {
					t.Errorf("expecting InvalidCurrency, got %v", err)
				}
			},
		},
		{
			name:  "converting to same currency",
			date:  time.Now(),
			from:  currency.USD,
			to:    currency.USD,
			value: 10,
			verify: func(value float64, err error) {
				if err != nil {
					t.Error(err)
				}
				if value != 10 {
					t.Errorf("expecting same value 10, got: %v", value)
				}
			},
		},
		{
			name:  "converting to EUR",
			date:  time.Date(2022, 1, 4, 0, 0, 0, 0, time.Local),
			from:  currency.USD,
			to:    currency.EUR,
			value: 10,
			GetRatesMock: func() (*ECBResponseData, error) {
				return &ECBResponseData{
					Data: []DataXML{
						{
							Date:  DateXML("2022-1-4"),
							Rates: []RateXML{{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3}},
						},
					},
				}, nil
			},
			verify: func(value float64, err error) {
				if err != nil {
					t.Error(err)
				}
				if value != 5 {
					t.Errorf("expecting value 5, got: %v", value)
				}
			},
		},
		{
			name:  "converting from EUR",
			date:  time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local),
			from:  currency.EUR,
			to:    currency.USD,
			value: 10,
			GetRatesMock: func() (*ECBResponseData, error) {
				return &ECBResponseData{
					Data: []DataXML{
						{
							Date: DateXML("2022-1-4"),
							Rates: []RateXML{
								{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3},
							},
						},
						{
							Date: DateXML("2022-1-3"),
							Rates: []RateXML{
								{Currency: "USD", Rate: 3}, {Currency: "JPY", Rate: 3},
							},
						},
					},
				}, nil
			},
			verify: func(value float64, err error) {
				if err != nil {
					t.Error(err)
				}
				if value != 30 {
					t.Errorf("expecting value 30, got: %v", value)
				}
			},
		},
		{
			name:  "converting from USD to JPY",
			date:  time.Date(2022, 1, 4, 0, 0, 0, 0, time.Local),
			from:  currency.USD,
			to:    currency.JPY,
			value: 14,
			GetRatesMock: func() (*ECBResponseData, error) {
				return &ECBResponseData{
					Data: []DataXML{
						{
							Date:  DateXML("2022-1-4"),
							Rates: []RateXML{{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3}},
						},
					},
				}, nil
			},
			verify: func(value float64, err error) {
				if err != nil {
					t.Error(err)
				}
				if value != 21 {
					t.Errorf("expecting value 21, got: %v", value)
				}
			},
		},
		{
			name: "get rates return error",
			date: time.Date(2022, 1, 4, 0, 0, 0, 0, time.Local),
			from: currency.USD,
			to:   currency.JPY,
			GetRatesMock: func() (*ECBResponseData, error) {
				return nil, errors.New("some error")
			},
			verify: func(value float64, err error) {
				if err == nil {
					t.Errorf("expecting error")
				}
			},
		},
		{
			name: "making rates hash map return error",
			date: time.Date(2022, 1, 4, 0, 0, 0, 0, time.Local),
			from: currency.USD,
			to:   currency.JPY,
			GetRatesMock: func() (*ECBResponseData, error) {
				return &ECBResponseData{
					Data: []DataXML{
						{
							Date: DateXML("INVALID"),
						},
					},
				}, nil
			},
			verify: func(value float64, err error) {
				if _, ok := err.(InvalidDateFormat); !ok {
					t.Errorf("expecting InvalidDateFormat got: %v", err)
				}
			},
		},
		{
			name: "query date too old",
			date: time.Date(2020, 5, 5, 0, 0, 0, 0, time.Local),
			from: currency.USD,
			to:   currency.JPY,
			GetRatesMock: func() (*ECBResponseData, error) {
				return &ECBResponseData{
					Data: []DataXML{
						{
							Date:  DateXML("2022-1-4"),
							Rates: []RateXML{{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3}},
						},
						{
							Date:  DateXML("2021-4-4"),
							Rates: []RateXML{{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3}},
						},
					},
				}, nil
			},
			verify: func(value float64, err error) {
				if _, ok := err.(DateOutOfBound); !ok {
					t.Errorf("expecting DateOutOfBound, got: %v", err)
				}
			},
		},
		{
			name: "query date too new",
			date: time.Date(2022, 1, 5, 0, 0, 0, 0, time.Local),
			from: currency.USD,
			to:   currency.JPY,
			GetRatesMock: func() (*ECBResponseData, error) {
				return &ECBResponseData{
					Data: []DataXML{
						{
							Date:  DateXML("2022-1-4"),
							Rates: []RateXML{{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3}},
						},
						{
							Date:  DateXML("2021-4-4"),
							Rates: []RateXML{{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3}},
						},
					},
				}, nil
			},
			verify: func(value float64, err error) {
				if _, ok := err.(DateOutOfBound); !ok {
					t.Errorf("expecting DateOutOfBound, got: %v", err)
				}
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			logger := log.New()
			converter := New(&ECBClientMock{GetRatesMock: test.GetRatesMock}, false, logger)
			test.verify(converter.Convert(test.date, test.value, test.from, test.to))
		})
	}
}
