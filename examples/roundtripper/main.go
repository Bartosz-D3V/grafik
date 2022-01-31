package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type authedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (at *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.token))
	return at.wrapped.RoundTrip(req)
}

func main() {
	const githubUrl = "https://api.github.com/graphql"
	githubApiToken := os.Getenv("GITHUB_TOKEN")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &authedTransport{
			token:   githubApiToken,
			wrapped: http.DefaultTransport,
		},
	}
	client := New(githubUrl, httpClient)

	res, err := client.GetData(context.Background(), nil)
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
	var graphqlRes GetDataResponse
	err = json.Unmarshal(b, &graphqlRes)
	if err != nil {
		panic(err)
	}

	log.Printf("Fetched %d repositories metadata from Github profile\n", len(graphqlRes.Data.Viewer.Repositories.Edges))
	for _, repo := range graphqlRes.Data.Viewer.Repositories.Edges {
		log.Printf("Repository name=%s\nProgramming languages used=%s\nWatchers=%d\n\n", repo.Node.Name, mapLanguages(repo), repo.Node.Watchers.TotalCount)
	}
}

func mapLanguages(repo RepositoryEdge) string {
	langNodes := repo.Node.Languages.Nodes
	languages := make([]string, len(langNodes))
	for i, language := range langNodes {
		languages[i] = language.Name
	}
	return strings.Join(languages, ", ")
}
