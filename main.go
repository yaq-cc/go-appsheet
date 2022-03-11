package appsheet

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func (c *AppSheetClient) Execute(ctx context.Context, r *AppSheetRequest) (*http.Response, error) {
	url := strings.Replace(URL, "{appId}", c.ApplicationId, -1)
	url = strings.Replace(url, "{tableName}", r.Table, -1)
	var body bytes.Buffer
	json.NewEncoder(&body).Encode(r)
	req, err := http.NewRequestWithContext(ctx, "POST", url, &body)
	if err != nil {
		log.Fatal(err)
	}
	return c.Client.Do(req)
}

type isRow interface {
	isRow()
}

type AppSheetRequest struct {
	Table      string     `json:"-"`
	Action     string     `json:"Action"`
	Properties *Properties `json:"Properties"`
	Rows       []isRow    `json:"Rows"`
}

func NewAppSheetRequest(table, action string) *AppSheetRequest {
	return &AppSheetRequest{
		Table: table,
		Action: action,
		Properties: &Properties{
			Locale: "en-US",
		},
		Rows: []isRow{},
	}
}

func (r *AppSheetRequest) AddRows(rows ...isRow) *AppSheetRequest {
	r.Rows = append(r.Rows, rows...)
	return r
}

type Properties struct {
	Locale       string            `json:"Locale"`
	Location     string            `json:"Location"`
	Timezone     string            `json:"Timezone"`
	UserSettings map[string]string `json:"UserSettings"`
}
