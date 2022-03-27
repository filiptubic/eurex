package ecb

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RateXML is string type alias used for unmarshaling date from
// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml.
type DateXML string

func (d DateXML) toTime() (time.Time, error) {
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

// RateXML is type used for unmarshaling rate XML data from
// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml.
type RateXML struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

// DataXML is type used for unmarshaling XML data from
// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml.
type DataXML struct {
	Date  DateXML   `xml:"time,attr"`
	Rates []RateXML `xml:"Cube"`
}

// ECBResponseData is type used for unmarshaling XML data from
// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml.
type ECBResponseData struct {
	Data []DataXML `xml:"Cube>Cube"`
}
