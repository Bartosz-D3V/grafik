package client

import "fmt"

type GraphQLCallError struct {
	Message string
	Reason  string
}

func (e GraphQLCallError) Error() string {
	return fmt.Sprintf("GraphQL call failed. Message=%s Reason=%s", e.Message, e.Reason)
}
