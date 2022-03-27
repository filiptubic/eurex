package ecb

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/avast/retry-go"
	log "github.com/sirupsen/logrus"
)

func TestECBClient_GetRates(t *testing.T) {
	tt := []struct {
		name    string
		handler http.HandlerFunc
		verify  func(data *ECBResponseData, err error)
	}{
		{
			name: "valid response",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data := ECBResponseData{
					Data: []DataXML{
						{
							Date: DateXML("2002-2-2"),
							Rates: []RateXML{
								{Currency: "USD", Rate: 1.234},
							},
						},
					},
				}
				dataMarshaled, _ := xml.Marshal(data)
				_, _ = w.Write(dataMarshaled)
			}),
			verify: func(data *ECBResponseData, err error) {
				if err != nil {
					t.Fatal(err)
				}
				if len(data.Data) != 1 {
					t.Fatalf("expecting one data, got: %d", len(data.Data))
				}
				if data.Data[0].Date != DateXML("2002-2-2") {
					t.Errorf("expecting '2002-2-2' got: %s", data.Data[0].Date)
				}
			},
		},
		{
			name: "4xx error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			verify: func(data *ECBResponseData, err error) {
				if _, ok := err.(ECBClientError); !ok {
					t.Fatalf("expecting ECBClientError from server, got: %v", err)
				}
			},
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			verify: func(data *ECBResponseData, err error) {
				if _, ok := err.(retry.Error); !ok {
					t.Fatalf("expecting retry.Error from server, got: %v", err)
				}
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(test.handler)
			defer ts.Close()

			url, _ := url.Parse(ts.URL)

			client := NewECBClient(url.Scheme, url.Host, NewECBOptions(0, time.Second), log.New())
			test.verify(client.GetRates())
		})
	}
}

func TestECBClient_GetRates_retry(t *testing.T) {
	retries := []int{0, 3, 5}
	for _, retry := range retries {
		called := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()
		url, _ := url.Parse(ts.URL)
		client := NewECBClient(url.Scheme, url.Host, NewECBOptions(retry, time.Second*0), log.New())
		client.GetRates()

		if called-1 != retry {
			t.Errorf("expecting %d retries, got %d", retry, called-1)
		}
	}
}
