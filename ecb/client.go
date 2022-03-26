package ecb

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
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
	Url string
}

func (c *ECBClient) GetRates() (*ECBResponseData, error) {
	resp, err := http.Get(c.Url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, ECBClientError{statusCode: resp.StatusCode}
	}

	respDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ecbData := ECBResponseData{}
	xml.Unmarshal(respDataBytes, &ecbData)
	return &ecbData, err
}
