// Generated with ggrafik. DO NOT EDIT

package main

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type Country struct {
	Name      string     `json:"name"`
	Native    string     `json:"native"`
	Capital   string     `json:"capital"`
	Currency  string     `json:"currency"`
	Languages []Language `json:"languages"`
	Emoji     string     `json:"emoji"`
}

type Language struct {
	Name string `json:"name"`
}

const getPolandInfo = `query getPolandInfo {
    country(code: "PL") {
        name
        native
        emoji
        currency
        languages {
            name
        }
    }
}`

type CountriesClient interface {
	GetPolandInfo(header *http.Header) (*http.Response, error)
}

func (c *countriesClient) GetPolandInfo(header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(getPolandInfo, params, header)
}

type GetPolandInfoResponse struct {
	Data   GetPolandInfoData `json:"data"`
	Errors []Error           `json:"errors"`
}

type GetPolandInfoData struct {
	Country Country `json:"country"`
}

type Error struct {
	Message    string     `json:"message"`
	Locations  []Location `json:"locations"`
	Extensions Extension  `json:"extensions"`
}

type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Extension struct {
	Code string `json:"code"`
}

type countriesClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CountriesClient {
	return &countriesClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
