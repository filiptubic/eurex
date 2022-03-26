package ecb

import "fmt"

type InvalidCurrency struct {
	currency string
}

func (e InvalidCurrency) Error() string {
	return fmt.Sprintf("invalid currency: %s", e.currency)
}

type InvalidDateFormat struct {
	date   string
	layout string
}

func (e InvalidDateFormat) Error() string {
	return fmt.Sprintf("expected date layout: %s, got %s", e.layout, e.date)
}

type DateParseError struct {
	msg string
}

func (e DateParseError) Error() string {
	return fmt.Sprintf("%s", e.msg)
}
