package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGraphQLCallError_Error(t *testing.T) {
	t.Parallel()
	err := GraphQLCallError{
		Message: "explaining message from library",
		Reason:  "originating error",
	}
	expErrMsg := "GraphQL call failed. Message=explaining message from library Reason=originating error"

	assert.Equal(t, expErrMsg, err.Error())
}
