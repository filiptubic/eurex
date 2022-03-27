package eurex

import (
	"testing"
	"time"

	"github.com/filiptubic/eurex/currency"
	"github.com/filiptubic/eurex/ecb"
)

func TestConverter_Convert(t *testing.T) {
	date := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
	converter := New(ecb.New(&ecb.ECBClientMock{GetRatesMock: func() (*ecb.ECBResponseData, error) {
		return &ecb.ECBResponseData{
			Data: []ecb.Data{
				{
					Date: "2022-1-1",
					Rates: []ecb.Rate{
						{
							Currency: "USD",
							Rate:     2,
						},
					},
				},
			},
		}, nil
	}}, true, DefaultLogger))

	value, err := converter.Convert(date, 100, currency.USD, currency.EUR)
	if err != nil {
		t.Error(err)
	}
	if value != 50 {
		t.Errorf("expecting 100 USD to be 50 EUR, got %v", value)
	}

}
