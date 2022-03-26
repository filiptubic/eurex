package ecb

import (
	"testing"
	"time"
)

func TestIsValid(t *testing.T) {
	tt := []struct {
		name     string
		date     time.Time
		currency Currency
		expected bool
	}{
		{
			name:     "CZK is valid",
			date:     time.Now(),
			currency: CZK,
			expected: true,
		},
		{
			name:     "invalid currency",
			date:     time.Now(),
			currency: Currency("foo"),
			expected: false,
		},
		{
			name:     "RUB before 1. March",
			date:     time.Date(2022, time.February, 23, 0, 0, 0, 0, time.Local),
			currency: RUB,
			expected: true,
		},
		{
			name:     "RUB after 1. March",
			date:     time.Date(2022, time.March, 2, 0, 0, 0, 0, time.Local),
			currency: RUB,
			expected: false,
		},
	}

	for _, test := range tt {
		t.Run(t.Name(), func(t *testing.T) {
			if ok := IsValid(string(test.currency), test.date); ok != test.expected {
				t.Errorf("expected %v, got %v", test.expected, ok)
			}
		})
	}
}
