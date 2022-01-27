// Generated with grafik. DO NOT EDIT

package main

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Users struct {
	Id string `json:"id"`
}

type UsersConstraint string

const (
	UsersPkey UsersConstraint = "users_pkey"
)

type UsersMutationResponse struct {
	AffectedRows int     `json:"affected_rows"`
	Returning    []Users `json:"returning"`
}

type UsersOnConflict struct {
	Constraint    UsersConstraint     `json:"constraint"`
	UpdateColumns []UsersUpdateColumn `json:"update_columns"`
}

type UsersUpdateColumn string

const (
	Id        UsersUpdateColumn = "id"
	Name      UsersUpdateColumn = "name"
	Rocket    UsersUpdateColumn = "rocket"
	Timestamp UsersUpdateColumn = "timestamp"
	Twitter   UsersUpdateColumn = "twitter"
)

const addOrUpdateHardcodedUser = `mutation addOrUpdateHardcodedUser($rocketName: String, $usersOnConflict: users_on_conflict) {
        insert_users(objects: {id: "5b8bcf27-9561-4123-87ff-75088c9da9c7", rocket: $rocketName}, on_conflict: $usersOnConflict) {
            affected_rows
            returning {
                id
            }
        }
}`

type SpaceXClient interface {
	AddOrUpdateHardcodedUser(ctx context.Context, rocketName string, usersOnConflict UsersOnConflict, header *http.Header) (*http.Response, error)
}

func (c *spaceXClient) AddOrUpdateHardcodedUser(ctx context.Context, rocketName string, usersOnConflict UsersOnConflict, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 2)
	params["rocketName"] = rocketName
	params["usersOnConflict"] = usersOnConflict

	return c.ctrl.Execute(ctx, addOrUpdateHardcodedUser, params, header)
}

type AddOrUpdateHardcodedUserResponse struct {
	Data   AddOrUpdateHardcodedUserData `json:"data"`
	Errors []GraphQLError               `json:"errors"`
}

type AddOrUpdateHardcodedUserData struct {
	InsertUsers UsersMutationResponse `json:"insert_users"`
}

type GraphQLError struct {
	Message    string                 `json:"message"`
	Locations  []GraphQLErrorLocation `json:"locations"`
	Extensions GraphQLErrorExtensions `json:"extensions"`
}

type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type GraphQLErrorExtensions struct {
	Code string `json:"code"`
}

type spaceXClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) SpaceXClient {
	return &spaceXClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
