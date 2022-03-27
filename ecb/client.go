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

// ECBClientInterface defines ECB client API.
type ECBClientInterface interface {
	GetRates() (*ECBResponseData, error)
}

// ECBClientMock type used for mocking http layer in tests.
type ECBClientMock struct {
	GetRatesMock func() (*ECBResponseData, error)
}

// GetRates calls GetRatesMock.
func (c *ECBClientMock) GetRates() (*ECBResponseData, error) {
	return c.GetRatesMock()
}

// ECBOptions allow client configuration, eg: specify retry count and their delay.
type ECBOptions struct {
	retry int
	wait  time.Duration
}

// NewECBOptions creates ECBOptions object.
func NewECBOptions(retry int, wait time.Duration) *ECBOptions {
	return &ECBOptions{
		retry: retry,
		wait:  wait,
	}
}

// ECBClient implements ECBClientInterface, therefore implements how rates are fetched via ECB. Additionally it can be configured using ECBOptions.
type ECBClient struct {
	logger  *log.Logger
	scheme  string
	host    string
	options *ECBOptions
}

// NewECBClient creates new ECBClient.
func NewECBClient(scheme, host string, options *ECBOptions, logger *log.Logger) *ECBClient {
	return &ECBClient{
		logger:  logger,
		scheme:  scheme,
		host:    host,
		options: options,
	}
}

// GetRates make http request to ECB to fetch rates data in form of XML. In case of non 2xx status code it fails with ECBClientError.
// If code is 5xx then it will try to retry requests using policy specified in ECBOptions.
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
