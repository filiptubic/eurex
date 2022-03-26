package ecb

import "time"

type Currency string

var (
	AUD = Currency("AUD")
	BGN = Currency("BGN")
	BRL = Currency("BRL")
	CAD = Currency("CAD")
	CHF = Currency("CHF")
	CNY = Currency("CNY")
	CZK = Currency("CZK")
	DKK = Currency("DKK")
	EUR = Currency("EUR")
	GBP = Currency("GBP")
	HKD = Currency("HKD")
	HRK = Currency("HRK")
	ISK = Currency("ISK")
	HUF = Currency("HUF")
	IDR = Currency("IDR")
	ILS = Currency("ILS")
	INR = Currency("INR")
	JPY = Currency("JPY")
	KRW = Currency("KRW")
	MXN = Currency("MXN")
	MYR = Currency("MYR")
	NOK = Currency("NOK")
	NZD = Currency("NZD")
	PHP = Currency("PHP")
	PLN = Currency("PLN")
	RON = Currency("RON")
	RUB = Currency("RUB")
	SEK = Currency("SEK")
	SGD = Currency("SGD")
	THB = Currency("THB")
	TRY = Currency("TRY")
	USD = Currency("USD")
	ZAR = Currency("ZAR")

	validCurrencies = map[string]struct{}{
		string(AUD): {},
		string(BGN): {},
		string(BRL): {},
		string(CAD): {},
		string(CHF): {},
		string(CNY): {},
		string(CZK): {},
		string(DKK): {},
		string(EUR): {},
		string(GBP): {},
		string(HKD): {},
		string(HRK): {},
		string(ISK): {},
		string(HUF): {},
		string(IDR): {},
		string(ILS): {},
		string(INR): {},
		string(JPY): {},
		string(KRW): {},
		string(MXN): {},
		string(MYR): {},
		string(NOK): {},
		string(NZD): {},
		string(PHP): {},
		string(PLN): {},
		string(RON): {},
		string(RUB): {},
		string(SEK): {},
		string(SGD): {},
		string(THB): {},
		string(TRY): {},
		string(USD): {},
		string(ZAR): {},
	}
)

func IsValid(currency string, date time.Time) bool {
	_, ok := validCurrencies[currency]
	if !ok {
		return false
	}

	rubValidTo := time.Date(2022, time.March, 1, 0, 0, 0, 0, time.Local)
	if currency == string(RUB) && date.After(rubValidTo) {
		return false
	}

	return true
}
