// Package client contains the code used internally by grafik to prepare & send HTTP requests.
package client

// GraphQLRequest is a root level struct generated by grafik.
// It corresponds to GraphQL HTTP request as per specification: https://graphql.org/learn/serving-over-http/#post-request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}
