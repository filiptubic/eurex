package ecb

import (
	"testing"
	"time"
)

func TestDate_toTime(t *testing.T) {
	tt := []struct {
		name   string
		date   DateXML
		verify func(value time.Time, err error)
	}{
		{
			name: "invalid date format",
			date: DateXML("12/11"),
			verify: func(value time.Time, err error) {
				if _, ok := err.(InvalidDateFormat); !ok {
					t.Errorf("expected InvalidDateFormat got: %v", err)
				}
			},
		},
		{
			name: "invalid year date",
			date: "asdf-12-22",
			verify: func(value time.Time, err error) {
				if _, ok := err.(DateParseError); !ok {
					t.Errorf("expected DateParseError, got: %v", err)
				}
			},
		},
		{
			name: "invalid month date",
			date: "2002-d-22",
			verify: func(value time.Time, err error) {
				if _, ok := err.(DateParseError); !ok {
					t.Errorf("expected DateParseError, got: %v", err)
				}
			},
		},
		{
			name: "invalid day date",
			date: "2022-21-a",
			verify: func(value time.Time, err error) {
				if _, ok := err.(DateParseError); !ok {
					t.Errorf("expected DateParseError, got: %v", err)
				}
			},
		},
		{
			name: "ok date",
			date: "2002-1-1",
			verify: func(value time.Time, err error) {
				if err != nil {
					t.Errorf("got error: %v", err)
				}
				if value != time.Date(2002, 1, 1, 0, 0, 0, 0, time.Local) {
					t.Errorf("invalid date: %v", value)
				}
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			test.verify(test.date.toTime())
		})
	}
}
