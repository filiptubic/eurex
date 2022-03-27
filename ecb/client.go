package ecb

import (
	"encoding/xml"
	"io/ioutil"
	"time"

	"net/http"
	"net/url"

	"github.com/avast/retry-go"
	log "github.com/sirupsen/logrus"
)

const (
	getRatesPath = "stats/eurofxref/eurofxref-hist-90d.xml"
)

type ECBClientInterface interface {
	GetRates() (*ECBResponseData, error)
}

type ECBClientMock struct {
	GetRatesMock func() (*ECBResponseData, error)
}

func (c *ECBClientMock) GetRates() (*ECBResponseData, error) {
	return c.GetRatesMock()
}

type ECBOptions struct {
	retry int
	wait  time.Duration
}

func NewECBOptions(retry int, wait time.Duration) *ECBOptions {
	return &ECBOptions{
		retry: retry,
		wait:  wait,
	}
}

type ECBClient struct {
	logger  *log.Logger
	scheme  string
	host    string
	options *ECBOptions
}

func NewECBClient(scheme, host string, options *ECBOptions, logger *log.Logger) *ECBClient {
	return &ECBClient{
		logger:  logger,
		scheme:  scheme,
		host:    host,
		options: options,
	}
}

func (c *ECBClient) GetRates() (*ECBResponseData, error) {
	url := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   getRatesPath,
	}

	var resp *http.Response
	var err error

	err = retry.Do(
		func() error {
			resp, err = http.Get(url.String())
			if err != nil {
				return err
			}

			if resp.StatusCode/100 == 5 {
				return ECBClientError{statusCode: resp.StatusCode}
			}

			return nil
		},
		retry.Attempts(uint(c.options.retry+1)),
		retry.Delay(c.options.wait),
		retry.OnRetry(func(n uint, err error) {
			c.logger.Errorf("[retry=%d] retrying on %v", n, err)
		}),
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		c.logger.Errorf("[GET] %v: code=%d", url.String(), resp.StatusCode)
		return nil, ECBClientError{statusCode: resp.StatusCode}
	}
	c.logger.Debugf("[GET] %v: code=%d", url.String(), resp.StatusCode)

	respDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ecbData := ECBResponseData{}
	xml.Unmarshal(respDataBytes, &ecbData)
	return &ecbData, err
}
