package client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_Execute_DefaultHeader_Success(t *testing.T) {
	exp := createCountriesResponse()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleRequest(t, exp, w, r)
	}))

	query := `
    query getContinentNameByCode($code: ID!) {
        continent(code: $code) {
            code,
            name
        }
    }
`
	client := New(svr.URL, svr.Client())
	params := make(map[string]interface{})
	params["code"] = "EU"
	res, err := client.Execute(query, params, nil)

	b, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	var countriesRes countriesResponse
	err = json.Unmarshal(b, &countriesRes)
	assert.NoError(t, err)
	assert.EqualValues(t, createCountriesResponse(), countriesRes)
}

func TestClient_Execute_CustomHeader_Success(t *testing.T) {
	expHeader := http.Header{
		"Authorization":  {"Bearer token"},
		"X-Trace-Id":     {"081cdde6-43a6-4c1c-8be1-2137af048273"},
		"Correlation-Id": {"4d82349b-821e-401c-b597-ee8f29cdb70a"},
	}
	exp := createCountriesResponse()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range expHeader {
			assert.Equal(t, v[0], r.Header.Get(k))
		}
		handleRequest(t, exp, w, r)
	}))

	query := `
    query($code: ID!) {
        continent(code: $code) {
            code,
            name
        }
    }
`
	client := New(svr.URL, svr.Client())
	params := make(map[string]interface{})
	params["code"] = "EU"
	res, err := client.Execute(query, params, &expHeader)

	b, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	var countriesRes countriesResponse
	err = json.Unmarshal(b, &countriesRes)
	assert.NoError(t, err)
	assert.EqualValues(t, createCountriesResponse(), countriesRes)
}

func TestClient_Execute_Marshall_Error(t *testing.T) {
	client := New("localhost:8080", http.DefaultClient)
	params := make(map[string]interface{})
	params["breakingParam"] = make(chan int)
	res, err := client.Execute("", params, nil)

	assert.Nil(t, res)
	expErr := GraphQLCallError{
		Message: "Parsing GraphQL request failed",
		Reason:  "json: unsupported type: chan int",
	}
	assert.ErrorIs(t, err, expErr)
}

func TestClient_Execute_NewRequest_Error(t *testing.T) {
	client := New("http://localhost:%%%8080", http.DefaultClient)
	params := make(map[string]interface{})
	res, err := client.Execute("", params, nil)

	assert.Nil(t, res)
	expErr := GraphQLCallError{
		Message: "Preparation of GraphQL call failed",
		Reason:  "parse \"http://localhost:%%%8080\": invalid port \":%%%8080\" after host",
	}
	assert.ErrorIs(t, err, expErr)
}

func TestClient_Execute_HttpDo_Error(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	svr.URL = "smtp://localhost:8000"

	client := New(svr.URL, svr.Client())
	params := make(map[string]interface{})
	res, err := client.Execute("", params, nil)

	assert.Nil(t, res)
	expErr := GraphQLCallError{
		Message: "GraphQL call failed",
		Reason:  `Post "smtp://localhost:8000": unsupported protocol scheme "smtp"`,
	}
	assert.ErrorIs(t, err, expErr)
}

func handleRequest(t *testing.T, exp countriesResponse, w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(exp)
	assert.NoError(t, err)

	_, err = w.Write(b)
	assert.NoError(t, err)

	// Check content-type
	assert.NotNil(t, r.Header.Get("Content-Type"))
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

	// Check payload - variables
	resBytes, err := io.ReadAll(r.Body)
	assert.NoError(t, err)
	var graphqlReq GraphQLRequest
	err = json.Unmarshal(resBytes, &graphqlReq)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(graphqlReq.Variables))
	assert.Equal(t, "EU", graphqlReq.Variables["code"])

	// Check payload - query
	assert.NotNil(t, graphqlReq.Query)
	assert.False(t, strings.Contains(graphqlReq.Query, "\n"))
	assert.False(t, strings.Contains(graphqlReq.Query, "\t"))
}

func createCountriesResponse() countriesResponse {
	return countriesResponse{
		Data: data{
			Continent: continent{
				Code: "EU",
				Name: "Europe",
			},
		},
	}
}

type continent struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type data struct {
	Continent continent `json:"continent"`
}

type countriesResponse struct {
	Data data `json:"data"`
}
