package ecb

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Date string

func (d Date) toTime() (time.Time, error) {
	parts := strings.Split(string(d), "-")
	if len(parts) != 3 {
		return time.Time{}, InvalidDateFormat{layout: "yyyy-dd-mm", date: string(d)}
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, DateParseError{
			msg: fmt.Sprintf("failed to parse year from date: %s, err: %v", d, err),
		}
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, DateParseError{
			msg: fmt.Sprintf("failed to parse month from date: %s, err: %v", d, err),
		}
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, DateParseError{
			msg: fmt.Sprintf("failed to parse day from date: %s, err: %v", d, err),
		}
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

type Rate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

type ECBResponseData struct {
	Data []struct {
		Date  Date   `xml:"time,attr"`
		Rates []Rate `xml:"Cube"`
	} `xml:"Cube>Cube"`
}
