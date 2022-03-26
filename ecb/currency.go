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

func IsValid(c currency.Currency, date time.Time) bool {
	_, ok := validCurrencies[c]
	if !ok {
		return false
	}

	rubValidTo := time.Date(2022, time.March, 1, 0, 0, 0, 0, time.Local)
	if c == currency.RUB && date.After(rubValidTo) {
		return false
	}

	return true
}
