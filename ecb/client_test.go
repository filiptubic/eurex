package ecb

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
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
					Data: []Data{
						{
							Date: Date("2002-2-2"),
							Rates: []Rate{
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
				if data.Data[0].Date != Date("2002-2-2") {
					t.Errorf("expecting '2002-2-2' got: %s", data.Data[0].Date)
				}
			},
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			verify: func(data *ECBResponseData, err error) {
				if err == nil {
					t.Fatalf("expecting error from server")
				}
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(test.handler)
			defer ts.Close()

			client := New(&ECBClient{Url: ts.URL})
			test.verify(client.client.GetRates())
		})
	}
}
