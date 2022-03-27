package ecb

import (
	"fmt"
	"net/http"
	"time"
)

func ExampleInvalidCurrency_Error() {
	fmt.Println(InvalidCurrency{currency: "INVALID"}.Error())
	// Outputs:
	// invalid currency: INVALID
}

func ExampleInvalidDateFormat_Error() {
	fmt.Println(InvalidDateFormat{layout: "yyyy-dd-mm", date: "1-1-1993"}.Error())
	// Outputs:
	// expected date layout: yyyy-dd-mm, got 1-1-1993
}

func ExampleDateParseError_Error() {
	fmt.Println(DateParseError{msg: "failed to parse year"}.Error())
	// Outputs:
	// failed to parse year
}

func ExampleECBClientError_Error() {
	fmt.Println(ECBClientError{statusCode: http.StatusForbidden}.Error())
	// Outputs:
	// http error: code=403
}

func ExampleDateOutOfBound_Error() {
	date := time.Date(1993, time.January, 1, 0, 0, 0, 0, time.Local)
	first := time.Date(2002, time.January, 1, 0, 0, 0, 0, time.Local)
	last := time.Date(2003, time.January, 1, 0, 0, 0, 0, time.Local)
	fmt.Println(DateOutOfBound{date: date, first: first, last: last}.Error())
	// Outputs:
	// 1993-01-01 00:00:00 +0100 CET out of date scope: [2002-01-01 00:00:00 +0100 CET, 2003-01-01 00:00:00 +0100 CET]
}
