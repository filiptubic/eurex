/*
	Use this package to access complete list of currencies. Various converters use various subsets of this list.
*/
package currency

// Currency is string type alias which represends all possible currencies.
// Its designed to be used in ConverterInterface API.
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

	Currencies = map[string]Currency{
		"AUD": AUD,
		"BGN": BGN,
		"BRL": BRL,
		"CAD": CAD,
		"CHF": CHF,
		"CNY": CNY,
		"CZK": CZK,
		"DKK": DKK,
		"EUR": EUR,
		"GBP": GBP,
		"HKD": HKD,
		"HRK": HRK,
		"ISK": ISK,
		"HUF": HUF,
		"IDR": IDR,
		"ILS": ILS,
		"INR": INR,
		"JPY": JPY,
		"KRW": KRW,
		"MXN": MXN,
		"MYR": MYR,
		"NOK": NOK,
		"NZD": NZD,
		"PHP": PHP,
		"PLN": PLN,
		"RON": RON,
		"RUB": RUB,
		"SEK": SEK,
		"SGD": SGD,
		"THB": THB,
		"TRY": TRY,
		"USD": USD,
		"ZAR": ZAR,
	}
)
