// Generated with ggrafik. DO NOT EDIT

package main

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type UsersConstraint string

const (
	UsersPkey UsersConstraint = "users_pkey"
)

type UsersUpdateColumn string

const (
	Id        UsersUpdateColumn = "id"
	Name      UsersUpdateColumn = "name"
	Rocket    UsersUpdateColumn = "rocket"
	Timestamp UsersUpdateColumn = "timestamp"
	Twitter   UsersUpdateColumn = "twitter"
)

type UsersMutationResponse struct {
	AffectedRows int     `json:"affected_rows"`
	Returning    []Users `json:"returning"`
}

type Users struct {
	Id string `json:"id"`
}

type UsersOnConflict struct {
	Constraint    UsersConstraint     `json:"constraint"`
	UpdateColumns []UsersUpdateColumn `json:"update_columns"`
}

const addOrUpdateHardcodedUser = `mutation addOrUpdateHardcodedUser($rocketName: String, $usersOnConflict: users_on_conflict) {
        insert_users(objects: {id: "5b8bcf27-9561-4123-87ff-75088c9da9c7", rocket: $rocketName}, on_conflict: $usersOnConflict) {
            affected_rows
            returning {
                id
            }
        }
}`

type UsersClient interface {
	AddOrUpdateHardcodedUser(rocketName string, usersOnConflict UsersOnConflict, header *http.Header) (*http.Response, error)
}

func (c *usersClient) AddOrUpdateHardcodedUser(rocketName string, usersOnConflict UsersOnConflict, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 2)
	params["rocketName"] = rocketName
	params["usersOnConflict"] = usersOnConflict

	return c.ctrl.Execute(addOrUpdateHardcodedUser, params, header)
}

type AddOrUpdateHardcodedUserResponse struct {
	Data AddOrUpdateHardcodedUserData `json:"data"`
}

type AddOrUpdateHardcodedUserData struct {
	InsertUsers UsersMutationResponse `json:"insert_users"`
}

type usersClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) UsersClient {
	return &usersClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
