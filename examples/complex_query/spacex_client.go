// Generated with grafik. DO NOT EDIT

package main

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Dragon struct {
	Name               string                   `json:"name"`
	Type               string                   `json:"type"`
	Wikipedia          string                   `json:"wikipedia"`
	PressurizedCapsule DragonPressurizedCapsule `json:"pressurized_capsule"`
}

type DragonPressurizedCapsule struct {
	PayloadVolume Volume `json:"payload_volume"`
}

type Info struct {
	Ceo string `json:"ceo"`
}

type Launchpad struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
}

type Location struct {
	Name string `json:"name"`
}

type Mission struct {
	Manufacturers []string `json:"manufacturers"`
}

type Roadster struct {
	Name      string `json:"name"`
	Wikipedia string `json:"wikipedia"`
}

type Volume struct {
	CubicMeters int `json:"cubic_meters"`
}

const getBatchInfo = `query getBatchInfo($limit: Int) {
    missions(limit: $limit) {
        manufacturers
    }
    launchpads(limit: $limit) {
        name
        location {
            name
        }
    }
    roadster {
        name
        wikipedia
    }
    company {
        ceo
    }
    dragons(limit: $limit) {
        wikipedia
        name
        type
        pressurized_capsule {
            payload_volume {
                cubic_meters
            }
        }
    }
}`

type SpaceXClient interface {
	GetBatchInfo(ctx context.Context, limit int, header http.Header) (*http.Response, error)
}

func (c *spaceXClient) GetBatchInfo(ctx context.Context, limit int, header http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["limit"] = limit

	return c.ctrl.Execute(ctx, getBatchInfo, params, header)
}

type GetBatchInfoResponse struct {
	Data   GetBatchInfoData `json:"data"`
	Errors []GraphQLError   `json:"errors"`
}

type GetBatchInfoData struct {
	Missions   []Mission   `json:"missions"`
	Launchpads []Launchpad `json:"launchpads"`
	Roadster   Roadster    `json:"roadster"`
	Company    Info        `json:"company"`
	Dragons    []Dragon    `json:"dragons"`
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
