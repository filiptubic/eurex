package ecb

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type ECBClient struct {
	url string
}

func (c *ECBClient) GetRates() (*ECBResponseData, error) {
	resp, err := http.Get(c.url)
	if resp.StatusCode/100 != 2 {
		return nil, ECBClientError{statusCode: resp.StatusCode}
	}
	if err != nil {
		return nil, err
	}
	respDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ecbData := ECBResponseData{}
	xml.Unmarshal(respDataBytes, &ecbData)
	return &ecbData, err
}
