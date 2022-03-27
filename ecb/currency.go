package ecb

import (
	"time"

	"github.com/filiptubic/eurex/currency"
)

var (
	validCurrencies = map[currency.Currency]struct{}{
		currency.AUD: {},
		currency.BGN: {},
		currency.BRL: {},
		currency.CAD: {},
		currency.CHF: {},
		currency.CNY: {},
		currency.CZK: {},
		currency.DKK: {},
		currency.EUR: {},
		currency.GBP: {},
		currency.HKD: {},
		currency.HRK: {},
		currency.ISK: {},
		currency.HUF: {},
		currency.IDR: {},
		currency.ILS: {},
		currency.INR: {},
		currency.JPY: {},
		currency.KRW: {},
		currency.MXN: {},
		currency.MYR: {},
		currency.NOK: {},
		currency.NZD: {},
		currency.PHP: {},
		currency.PLN: {},
		currency.RON: {},
		currency.RUB: {},
		currency.SEK: {},
		currency.SGD: {},
		currency.THB: {},
		currency.TRY: {},
		currency.USD: {},
		currency.ZAR: {},
	}
)

// IsValidCurrency check whether currency is valid for certain date.
func IsValidCurrency(c currency.Currency, date time.Time) bool {
	_, ok := validCurrencies[c]
	if !ok {
		return false
	}

	// Russian RUB is invalid on ECB after 1. March 2022
	rubValidTo := time.Date(2022, time.March, 1, 0, 0, 0, 0, time.Local)
	if c == currency.RUB && date.After(rubValidTo) {
		return false
	}

	return true
}
