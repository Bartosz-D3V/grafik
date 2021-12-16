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
	const spacexUrl = "https://api.spacex.land/graphql"
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := New(spacexUrl, httpClient)
	headers := &http.Header{
		"Date": []string{time.Now().String()},
	}
	onConflict := UsersOnConflict{
		Constraint: UsersPkey,
		UpdateColumns: []UsersUpdateColumn{
			Id,
			Name,
			Rocket,
			Timestamp,
			Twitter,
		},
	}
	res, err := client.AddOrUpdateHardcodedUser(context.Background(), "Falcon 1", onConflict, headers)
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
	var graphqlRes AddOrUpdateHardcodedUserResponse
	err = json.Unmarshal(b, &graphqlRes)
	if err != nil {
		panic(err)
	}

	log.Printf("Updated user with id=%s. Number of affected users=%d", graphqlRes.Data.InsertUsers.Returning[0].Id, graphqlRes.Data.InsertUsers.AffectedRows)
}
