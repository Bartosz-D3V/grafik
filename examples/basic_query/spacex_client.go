// Generated with ggrafik. DO NOT EDIT

package main

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type RocketsResult struct {
	Result Result   `json:"result"`
	Data   []Rocket `json:"data"`
}

type Result struct {
	TotalCount int `json:"totalCount"`
}

type Rocket struct {
	CostPerLaunch int    `json:"cost_per_launch"`
	Country       string `json:"country"`
	Name          string `json:"name"`
}

const getRocketResults = `query getRocketResults($limit: Int){
    rocketsResult(limit: $limit) {
        result {
            totalCount
        }
        data {
            name
            country
            cost_per_launch
        }
    }
}
`

type RocketsClient interface {
	GetRocketResults(limit int, header *http.Header) (*http.Response, error)
}

func (c *rocketsClient) GetRocketResults(limit int, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["limit"] = limit

	return c.ctrl.Execute(getRocketResults, params, header)
}

type GetRocketResultsResponse struct {
	Data GetRocketResultsData `json:"data"`
}

type GetRocketResultsData struct {
	RocketsResult RocketsResult `json:"rocketsResult"`
}

type rocketsClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) RocketsClient {
	return &rocketsClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
