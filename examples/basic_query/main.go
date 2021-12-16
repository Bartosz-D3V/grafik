package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	const (
		spacexUrl  = "https://api.spacex.land/graphql"
		maxResults = 2
	)
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := New(spacexUrl, httpClient)

	res, err := client.GetRocketResults(context.Background(), maxResults, nil)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var graphqlRes GetRocketResultsResponse
	err = json.Unmarshal(b, &graphqlRes)
	if err != nil {
		panic(err)
	}

	log.Printf("Fetched %d/%d results \n", maxResults, graphqlRes.Data.RocketsResult.Result.TotalCount)
	for i, r := range graphqlRes.Data.RocketsResult.Data {
		log.Printf("#%d: Rocket %s was launched from %s and costed $%d", i, r.Name, r.Country, r.CostPerLaunch)
	}
}
