package main

import (
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

	res, err := client.GetRocketResults(maxResults, nil)
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

	log.Printf("Found %d errors in GraphQL response \n", len(graphqlRes.Errors))
	for _, e := range graphqlRes.Errors {
		log.Printf("Message=%s\nLine=%d\nColumn=%d\nCode=%s", e.Message, e.Locations[0].Line, e.Locations[0].Column, e.Extensions.Code)
	}
}
