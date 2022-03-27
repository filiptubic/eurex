# eurex 

Eurex is a simple money conversion library which converts value from one currency to another. 
Rates which are applied are configurable and currently only supported rate source is from ECB ([European Central Bank](https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml)).


## Quickstart
### Install 
```
go get -u github.com/filiptubic/eurex
```
### Example

Convert 10 USD to CHF on 25. March 2022:
```
package main

import (
	"time"

	"github.com/filiptubic/eurex"
	"github.com/filiptubic/eurex/currency"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := eurex.DefaultLogger
	logger.SetLevel(logrus.DebugLevel)
	converter := eurex.DefaultConverter
	date := time.Date(2022, time.March, 25, 0, 0, 0, 0, time.Local)

	value, err := converter.Convert(date, 10, currency.USD, currency.CHF)
	if err != nil {
		logger.Errorf("got error: %v", err)
	}
	logger.Info(value)
}
```
Outputs:
```
DEBU[0000] [GET] https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml: code=200 
INFO[0000] 9.277404108343935                            
```

## Tests
```
go test -coverprofile=coverage.out ./... -v
```

## Docs
```
godoc
```