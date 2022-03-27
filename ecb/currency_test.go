package ecb

import (
	"fmt"
	"testing"
	"time"

	"github.com/filiptubic/eurex/currency"
)

func ExampleIsValidCurrency() {
	firstOfMarch := time.Date(2022, time.March, 1, 0, 0, 0, 0, time.Local)
	secondOfMarch := time.Date(2022, time.March, 2, 0, 0, 0, 0, time.Local)
	fmt.Println(IsValidCurrency(currency.RUB, firstOfMarch))
	fmt.Println(IsValidCurrency(currency.RUB, secondOfMarch))
	// Output:
	// true
	// false
}

func TestIsValidCurrency(t *testing.T) {
	tt := []struct {
		name     string
		date     time.Time
		currency currency.Currency
		expected bool
	}{
		{
			name:     "CZK is valid",
			date:     time.Now(),
			currency: currency.CZK,
			expected: true,
		},
		{
			name:     "invalid currency",
			date:     time.Now(),
			currency: currency.Currency("foo"),
			expected: false,
		},
		{
			name:     "RUB before 1. March",
			date:     time.Date(2022, time.February, 23, 0, 0, 0, 0, time.Local),
			currency: currency.RUB,
			expected: true,
		},
		{
			name:     "RUB after 1. March",
			date:     time.Date(2022, time.March, 2, 0, 0, 0, 0, time.Local),
			currency: currency.RUB,
			expected: false,
		},
	}

	for _, test := range tt {
		t.Run(t.Name(), func(t *testing.T) {
			if ok := IsValidCurrency(test.currency, test.date); ok != test.expected {
				t.Errorf("expected %v, got %v", test.expected, ok)
			}
		})
	}
}
