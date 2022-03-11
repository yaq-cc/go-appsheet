package appsheet

import (
	"net/http"
)

var URL string = "https://api.appsheet.com/api/v2/apps/{appId}/tables/{tableName}/Action"

type AppSheetClient struct {
	ApplicationId string
	Client        *http.Client
}

type AppSheetTransport struct {
	DefaultTransport     http.RoundTripper
	ApplicationAccessKey string
}

func (t *AppSheetTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("ApplicationAccessKey", t.ApplicationAccessKey)
	return t.DefaultTransport.RoundTrip(r)
}

func NewAppSheetClient(id, accessKey string) *AppSheetClient {
	return &AppSheetClient{
		ApplicationId: id,
		Client: &http.Client{
			Transport: &AppSheetTransport{
				ApplicationAccessKey: accessKey,
				DefaultTransport:     http.DefaultTransport,
			},
		},
	}
}

type isRow interface {
	isRow()
}

type AppSheetRequest struct {
	Action     string     `json:"Action"`
	Properties Properties `json:"Properties"`
	Rows       []isRow    `json:"Rows"`
}

type Properties struct {
	Locale       string            `json:"Locale"`
	Location     string            `json:"Location"`
	Timezone     string            `json:"Timezone"`
	UserSettings map[string]string `json:"UserSettings"`
}
