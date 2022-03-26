package ecb

import (
	"errors"
	"testing"
	"time"

	"github.com/filiptubic/eurex/currency"
)

func TestECBConverter_makeRatesMap(t *testing.T) {
	tt := []struct {
		name   string
		data   *ECBResponseData
		verify func(m map[time.Time]currencyMap, err error)
	}{
		{
			name: "ok map",
			data: &ECBResponseData{
				Data: []Data{
					{Date: Date("2022-1-4"), Rates: []Rate{{Currency: "USD", Rate: 1.5}}},
					{Date: Date("2022-1-3"), Rates: []Rate{{Currency: "TRY", Rate: 2}}},
				},
			},
			verify: func(m map[time.Time]currencyMap, err error) {
				if err != nil {
					t.Error(err)
				}

				if len(m) != 2 {
					t.Errorf("expected only two entry, got: %d", len(m))
				}
				t1 := time.Date(2022, 1, 4, 0, 0, 0, 0, time.Local)
				if _, ok := m[t1]; !ok {
					t.Errorf("missing %v time in map", t1)
				}
				t2 := time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local)
				if _, ok := m[t1]; !ok {
					t.Errorf("missing %v time in map", t2)
				}
			},
		},
		{
			name: "invalid date layout",
			data: &ECBResponseData{
				Data: []Data{
					{Date: Date("invalid date")},
				},
			},
			verify: func(m map[time.Time]currencyMap, err error) {
				if _, ok := err.(InvalidDateFormat); !ok {
					t.Errorf("expecting InvalidDateFormat got: %v", err)
				}
			},
		},
		{
			name: "invalid currency",
			data: &ECBResponseData{
				Data: []Data{
					{Date: Date("2022-12-12"), Rates: []Rate{{Currency: "UNKNOWN"}}},
				},
			},
			verify: func(m map[time.Time]currencyMap, err error) {
				if _, ok := err.(InvalidCurrency); !ok {
					t.Errorf("expecting InvalidCurrency got: %v", err)
				}
			},
		},
	}

	converter := New(&ECBClientMock{})
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			test.verify(converter.makeRatesMap(test.data))
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
					Data: []Data{
						{
							Date:  Date("2022-1-4"),
							Rates: []Rate{{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3}},
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
					Data: []Data{
						{
							Date: Date("2022-1-4"),
							Rates: []Rate{
								{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3},
							},
						},
						{
							Date: Date("2022-1-3"),
							Rates: []Rate{
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
					Data: []Data{
						{
							Date:  Date("2022-1-4"),
							Rates: []Rate{{Currency: "USD", Rate: 2}, {Currency: "JPY", Rate: 3}},
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
					Data: []Data{
						{
							Date: Date("INVALID"),
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
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			converter := New(&ECBClientMock{GetRatesMock: test.GetRatesMock})
			test.verify(converter.Convert(test.date, test.value, test.from, test.to))
		})
	}
}
