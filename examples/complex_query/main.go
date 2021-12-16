package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
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

	res, err := client.GetBatchInfo(context.Background(), maxResults, nil)
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
	var graphqlRes GetBatchInfoResponse
	err = json.Unmarshal(b, &graphqlRes)
	if err != nil {
		panic(err)
	}

	printResults(graphqlRes)
}

func printResults(graphqlRes GetBatchInfoResponse) {
	log.Printf("CEO of the SpaceX: %s\n", graphqlRes.Data.Company.Ceo)
	log.Printf("Missions fetched: %d\n", len(graphqlRes.Data.Missions))
	for i, mission := range graphqlRes.Data.Missions {
		log.Printf("Mission %d: Manufacturers: %s\n", i, strings.Join(mission.Manufacturers, ", "))
	}
	log.Printf("Launchpads fetched: %d\n", len(graphqlRes.Data.Launchpads))
	for i, launchpad := range graphqlRes.Data.Launchpads {
		log.Printf("Launchpad %d: Name: %s Location: %s\n", i, launchpad.Name, launchpad.Location.Name)
	}
	log.Printf("Dragons fetched: %d\n", len(graphqlRes.Data.Dragons))
	for i, dragon := range graphqlRes.Data.Dragons {
		log.Printf("Dragon %d: Name: %s Type: %s Wikipedia: %s Payload (m³): %d\n", i, dragon.Name, dragon.Type, dragon.Wikipedia, dragon.PressurizedCapsule.PayloadVolume.CubicMeters)
	}
	log.Printf("Roadster fetched: %s Wikipedia: %s\n", graphqlRes.Data.Roadster.Name, graphqlRes.Data.Roadster.Wikipedia)
}
