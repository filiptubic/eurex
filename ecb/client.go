package ecb

import (
	"encoding/xml"
	"io/ioutil"

	"net/http"

	log "github.com/sirupsen/logrus"
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

type ECBClient struct {
	logger *log.Logger
	Url    string
}

func NewECBClient(url string, logger *log.Logger) *ECBClient {
	return &ECBClient{
		logger: logger,
		Url:    url,
	}
}

func (c *ECBClient) GetRates() (*ECBResponseData, error) {
	resp, err := http.Get(c.Url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		c.logger.Errorf("[GET] %s: code=%d", c.Url, resp.StatusCode)
		return nil, ECBClientError{statusCode: resp.StatusCode}
	}
	c.logger.Debugf("[GET] %s: code=%d", c.Url, resp.StatusCode)

	respDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ecbData := ECBResponseData{}
	xml.Unmarshal(respDataBytes, &ecbData)
	return &ecbData, err
}
