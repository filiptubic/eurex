package ecb

import (
	"fmt"
	"time"
)

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

type ECBClientError struct {
	statusCode int
}

func (e ECBClientError) Error() string {
	return fmt.Sprintf("http error: code=%v", e.statusCode)
}

type DateOutOfBound struct {
	date        time.Time
	first, last time.Time
}

func (e DateOutOfBound) Error() string {
	return fmt.Sprintf("%v out of date scope: [%v, %v]", e.date, e.first, e.last)
}
