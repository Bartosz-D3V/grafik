package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type service struct {
	client SpaceXClient
}

func (s service) ReturnAverageCostForPerLaunch() (int, error) {
	res, err := s.client.GetRocketResults(context.Background(), 50, nil)
	if err != nil {
		return 0, fmt.Errorf("SpaceXClient failed: %s", err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var graphqlRes GetRocketResultsResponse
	err = json.Unmarshal(b, &graphqlRes)
	if err != nil {
		return 0, err
	}

	numOfMissions := len(graphqlRes.Data.RocketsResult.Data)
	totalCost := 0
	for _, rocket := range graphqlRes.Data.RocketsResult.Data {
		totalCost += rocket.CostPerLaunch
	}
	return totalCost / numOfMissions, nil
}
