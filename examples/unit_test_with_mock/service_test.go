package unit_test_example_with_mock

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestService_ReturnAverageCostForPerLaunch(t *testing.T) {
	t.Parallel()
	svc := service{mockSpaceXClient{true}}

	res, err := svc.ReturnAverageCostForPerLaunch()

	assert.NoError(t, err)
	assert.Equal(t, 49000000, res)
}

func TestService_ReturnAverageCostForPerLaunch_Failure(t *testing.T) {
	t.Parallel()
	svc := service{mockSpaceXClient{false}}

	res, err := svc.ReturnAverageCostForPerLaunch()

	assert.Errorf(t, err, "SpaceXClient failed: GraphQL call failed")
	assert.Equal(t, 0, res)
}

type mockSpaceXClient struct {
	returnValidResponse bool
}

func (c mockSpaceXClient) GetRocketResults(int, *http.Header) (*http.Response, error) {
	if c.returnValidResponse {
		res := createGraphQLResponse()
		resJson, _ := json.Marshal(res)
		r := io.NopCloser(bytes.NewReader(resJson))
		return &http.Response{
			Status:     "OK",
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	return nil, errors.New("GraphQL call failed")
}

func createGraphQLResponse() GetRocketResultsResponse {
	return GetRocketResultsResponse{
		Data: GetRocketResultsData{
			RocketsResult: RocketsResult{
				Data: []Rocket{
					{
						CostPerLaunch: 57000000,
					},
					{
						CostPerLaunch: 10000000,
					},
					{
						CostPerLaunch: 80000000,
					},
				},
			},
		},
	}
}
