package ecb

import (
	"fmt"
	"time"
)

// InvalidCurrency is used when currency is invalid on not registered for specific converter.
type InvalidCurrency struct {
	currency string
}

func (e InvalidCurrency) Error() string {
	return fmt.Sprintf("invalid currency: %s", e.currency)
}

// InvalidDateFormat is used when date is invalid format (eg. not in "yyyy-dd-mm" layout for ECB).
type InvalidDateFormat struct {
	date   string
	layout string
}

func (e InvalidDateFormat) Error() string {
	return fmt.Sprintf("expected date layout: %s, got %s", e.layout, e.date)
}

// DateParseError is used when date parsing failed due malformed year/month/day.
type DateParseError struct {
	msg string
}

func (e DateParseError) Error() string {
	return fmt.Sprintf("%s", e.msg)
}

// ECBClientError represents HTTP related errors (eg. 4xx status codes)
type ECBClientError struct {
	statusCode int
}

func (e ECBClientError) Error() string {
	return fmt.Sprintf("http error: code=%v", e.statusCode)
}

// DateOutOfBound is used when querying date is out of possible dates of conversion.
type DateOutOfBound struct {
	date        time.Time
	first, last time.Time
}

func (e DateOutOfBound) Error() string {
	return fmt.Sprintf("%v out of date scope: [%v, %v]", e.date, e.first, e.last)
}
