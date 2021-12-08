package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	const countriesUrl = "https://countries.trevorblades.com/"
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := New(countriesUrl, httpClient)

	res, err := client.GetPolandInfo(nil)
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
	var graphqlRes GetPolandInfoResponse
	err = json.Unmarshal(b, &graphqlRes)
	if err != nil {
		panic(err)
	}

	country := graphqlRes.Data.Country
	log.Printf("%s: %s means %s in %s and we use %s as a currency", country.Emoji, country.Native, country.Name, country.Languages[0].Name, country.Currency)
}
