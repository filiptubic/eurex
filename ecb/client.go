package ecb

import (
	"encoding/xml"
	"io/ioutil"
	"time"

	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

const (
	GetRatesPath = "stats/eurofxref/eurofxref-hist-90d.xml"
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
		Path:   GetRatesPath,
	}

	var resp *http.Response
	var err error
	retry := 0
	for {
		resp, err = http.Get(url.String())
		if err != nil {
			return nil, err
		}

		if resp.StatusCode/100 != 5 || retry == c.options.retry {
			break
		}
		retry++

		c.logger.Errorf("[GET] %v: retrying on 5xx error", url.String())
		time.Sleep(c.options.wait)
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
