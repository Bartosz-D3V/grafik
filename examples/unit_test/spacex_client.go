// Generated with grafik. DO NOT EDIT

package unit_test_with_grafik

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Result struct {
	TotalCount int `json:"totalCount"`
}

type Rocket struct {
	CostPerLaunch int    `json:"cost_per_launch"`
	Country       string `json:"country"`
	Name          string `json:"name"`
}

type RocketsResult struct {
	Result Result   `json:"result"`
	Data   []Rocket `json:"data"`
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
}`

type SpaceXClient interface {
	GetRocketResults(ctx context.Context, limit int, header *http.Header) (*http.Response, error)
}

func (c *spaceXClient) GetRocketResults(ctx context.Context, limit int, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["limit"] = limit

	return c.ctrl.Execute(ctx, getRocketResults, params, header)
}

type GetRocketResultsResponse struct {
	Data   GetRocketResultsData `json:"data"`
	Errors []GraphQLError       `json:"errors"`
}

type GetRocketResultsData struct {
	RocketsResult RocketsResult `json:"rocketsResult"`
}

type GraphQLError struct {
	Message    string                 `json:"message"`
	Locations  []GraphQLErrorLocation `json:"locations"`
	Extensions GraphQLErrorExtensions `json:"extensions"`
}

type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type GraphQLErrorExtensions struct {
	Code string `json:"code"`
}

type spaceXClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) SpaceXClient {
	return &spaceXClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
