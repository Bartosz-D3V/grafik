package evaluator

import (
	"bytes"
	"fmt"
	"github.com/Bartosz-D3V/grafik/test"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"io/ioutil"
	"path"
	"testing"
)

func TestEvaluator_Generate_FlatStructure(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/simple_type/schema.graphql")
	query := loadQuery(t, schema, "test/simple_type/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "FilesClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type File struct {
	Name string %[1]cjson:"name"%[1]c
}

const getFileNameWithId = %[1]cquery GetFileNameWithId($id: ID!) {
    getFile(id: $id) {
        name
    }
}%[1]c

const renameFileWithId = %[1]cmutation RenameFileWithId($id: ID!, $name: String!) {
    renameFile(id: $id, name: $name) {
        name
    }
}%[1]c

type FilesClient interface {
	GetFileNameWithId(ctx context.Context, id string, header *http.Header) (*http.Response, error)
	RenameFileWithId(ctx context.Context, id string, name string, header *http.Header) (*http.Response, error)
}

func (c *filesClient) GetFileNameWithId(ctx context.Context, id string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["id"] = id

	return c.ctrl.Execute(ctx, getFileNameWithId, params, header)
}

func (c *filesClient) RenameFileWithId(ctx context.Context, id string, name string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 2)
	params["id"] = id
	params["name"] = name

	return c.ctrl.Execute(ctx, renameFileWithId, params, header)
}

type GetFileNameWithIdResponse struct {
	Data   GetFileNameWithIdData %[1]cjson:"data"%[1]c
	Errors []GraphQLError        %[1]cjson:"errors"%[1]c
}

type GetFileNameWithIdData struct {
	GetFile File %[1]cjson:"getFile"%[1]c
}

type RenameFileWithIdResponse struct {
	Data   RenameFileWithIdData %[1]cjson:"data"%[1]c
	Errors []GraphQLError       %[1]cjson:"errors"%[1]c
}

type RenameFileWithIdData struct {
	RenameFile File %[1]cjson:"renameFile"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type filesClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) FilesClient {
	return &filesClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_ArrayStructure(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/array/schema.graphql")
	query := loadQuery(t, schema, "test/array/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "FilmsClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Film struct {
	Producers []string %[1]cjson:"producers"%[1]c
}

type FilmsConnection struct {
	Films []Film %[1]cjson:"films"%[1]c
}

const getAllFilmsProducers = %[1]cquery GetAllFilmsProducers {
    allFilms {
        films {
            producers
        }
    }
}%[1]c

type FilmsClient interface {
	GetAllFilmsProducers(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *filmsClient) GetAllFilmsProducers(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getAllFilmsProducers, params, header)
}

type GetAllFilmsProducersResponse struct {
	Data   GetAllFilmsProducersData %[1]cjson:"data"%[1]c
	Errors []GraphQLError           %[1]cjson:"errors"%[1]c
}

type GetAllFilmsProducersData struct {
	AllFilms FilmsConnection %[1]cjson:"allFilms"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type filmsClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) FilmsClient {
	return &filmsClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_2DArrayStructure(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/2d_array/schema.graphql")
	query := loadQuery(t, schema, "test/2d_array/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "MathClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Result struct {
	X int %[1]cjson:"x"%[1]c
	Y int %[1]cjson:"y"%[1]c
	Z int %[1]cjson:"z"%[1]c
}

const getAllResults = %[1]cquery GetAllResults {
    allResults {
        x,
        y,
        z
    }
    allResultsSimplified
}%[1]c

type MathClient interface {
	GetAllResults(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *mathClient) GetAllResults(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getAllResults, params, header)
}

type GetAllResultsResponse struct {
	Data   GetAllResultsData %[1]cjson:"data"%[1]c
	Errors []GraphQLError    %[1]cjson:"errors"%[1]c
}

type GetAllResultsData struct {
	AllResults           [][]Result %[1]cjson:"allResults"%[1]c
	AllResultsSimplified [][]string %[1]cjson:"allResultsSimplified"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type mathClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) MathClient {
	return &mathClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_3DArrayStructure(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/3d_array/schema.graphql")
	query := loadQuery(t, schema, "test/3d_array/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "MathClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Result struct {
	X int %[1]cjson:"x"%[1]c
	Y int %[1]cjson:"y"%[1]c
	Z int %[1]cjson:"z"%[1]c
}

const getAllResults = %[1]cquery GetAllResults {
    allResults {
        x,
        y,
        z
    }
}%[1]c

type MathClient interface {
	GetAllResults(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *mathClient) GetAllResults(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getAllResults, params, header)
}

type GetAllResultsResponse struct {
	Data   GetAllResultsData %[1]cjson:"data"%[1]c
	Errors []GraphQLError    %[1]cjson:"errors"%[1]c
}

type GetAllResultsData struct {
	AllResults [][][]Result %[1]cjson:"allResults"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type mathClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) MathClient {
	return &mathClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_NestedStructure(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/nested_type/schema.graphql")
	query := loadQuery(t, schema, "test/nested_type/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "SpecificHeroClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Character struct {
	HomeWorld Planet  %[1]cjson:"homeWorld"%[1]c
	Species   Species %[1]cjson:"species"%[1]c
}

type Location struct {
	PosX int %[1]cjson:"posX"%[1]c
	PoxY int %[1]cjson:"poxY"%[1]c
}

type Planet struct {
	Location Location %[1]cjson:"location"%[1]c
}

type Species struct {
	Origin Planet %[1]cjson:"origin"%[1]c
}

const getHeroWithId123ABC = %[1]cquery GetHeroWithId123ABC {
    getHero(characterSelector: {idSelector: {id: "123ABC"}}) {
        homeWorld {
            location {
                posX
            }
        }
        species {
            origin {
                location {
                    poxY
                }
            }
        }
    }
}%[1]c

type SpecificHeroClient interface {
	GetHeroWithId123ABC(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *specificHeroClient) GetHeroWithId123ABC(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getHeroWithId123ABC, params, header)
}

type GetHeroWithId123ABCResponse struct {
	Data   GetHeroWithId123ABCData %[1]cjson:"data"%[1]c
	Errors []GraphQLError          %[1]cjson:"errors"%[1]c
}

type GetHeroWithId123ABCData struct {
	GetHero Character %[1]cjson:"getHero"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type specificHeroClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) SpecificHeroClient {
	return &specificHeroClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_Enum(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/enum/schema.graphql")
	query := loadQuery(t, schema, "test/enum/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CompanyClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Department struct {
	Name DepartmentName %[1]cjson:"name"%[1]c
}

type DepartmentName string

const (
	IT      DepartmentName = "IT"
	SALES   DepartmentName = "SALES"
	HR      DepartmentName = "HR"
	SUPPORT DepartmentName = "SUPPORT"
)

const getDepartment = %[1]cquery getDepartment {
    getDepartment {
        name
    }
}%[1]c

type CompanyClient interface {
	GetDepartment(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *companyClient) GetDepartment(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getDepartment, params, header)
}

type GetDepartmentResponse struct {
	Data   GetDepartmentData %[1]cjson:"data"%[1]c
	Errors []GraphQLError    %[1]cjson:"errors"%[1]c
}

type GetDepartmentData struct {
	GetDepartment Department %[1]cjson:"getDepartment"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type companyClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CompanyClient {
	return &companyClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_Input(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/input/schema.graphql")
	query := loadQuery(t, schema, "test/input/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CapsulesClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Capsule struct {
	Id         string %[1]cjson:"id"%[1]c
	Type       string %[1]cjson:"type"%[1]c
}

type Date interface {

}

const getCapsulesByFullSelector = %[1]cquery GetCapsulesByFullSelector($order: String, $mission: String, $originalLaunch: Date, $id: ID, $sort: String) {
    capsules(order: $order, find: {landings: 10, mission: $mission, original_launch: $originalLaunch, id: $id}, sort: $sort) {
        id
        type
    }
}%[1]c

type CapsulesClient interface {
	GetCapsulesByFullSelector(ctx context.Context, order string, mission string, originalLaunch Date, id string, sort string, header *http.Header) (*http.Response, error)
}

func (c *capsulesClient) GetCapsulesByFullSelector(ctx context.Context, order string, mission string, originalLaunch Date, id string, sort string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 5)
	params["order"] = order
	params["mission"] = mission
	params["originalLaunch"] = originalLaunch
	params["id"] = id
	params["sort"] = sort

	return c.ctrl.Execute(ctx, getCapsulesByFullSelector, params, header)
}

type GetCapsulesByFullSelectorResponse struct {
	Data   GetCapsulesByFullSelectorData %[1]cjson:"data"%[1]c
	Errors []GraphQLError                %[1]cjson:"errors"%[1]c
}

type GetCapsulesByFullSelectorData struct {
	Capsules []Capsule %[1]cjson:"capsules"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type capsulesClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CapsulesClient {
	return &capsulesClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_Input_2DArray(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/input_2d_array/schema.graphql")
	query := loadQuery(t, schema, "test/input_2d_array/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CapsulesClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Capsule struct {
	Id         string %[1]cjson:"id"%[1]c
}

type Limit struct {
	Size int %[1]cjson:"size"%[1]c
}

type Position struct {
	X int %[1]cjson:"x"%[1]c
	Y int %[1]cjson:"y"%[1]c
}

const getCapsulesByPositions = %[1]cquery GetCapsulesByPositions($find: [[Position]], $limit: [[Limit]], $selector: [[String]]) {
   capsules(find: $find, limit: $limit, selector: $selector) {
       id
   }
}%[1]c

type CapsulesClient interface {
	GetCapsulesByPositions(ctx context.Context, find [][]Position, limit [][]Limit, selector [][]string, header *http.Header) (*http.Response, error)
}

func (c *capsulesClient) GetCapsulesByPositions(ctx context.Context, find [][]Position, limit [][]Limit, selector [][]string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 3)
	params["find"] = find
	params["limit"] = limit
	params["selector"] = selector

	return c.ctrl.Execute(ctx, getCapsulesByPositions, params, header)
}

type GetCapsulesByPositionsResponse struct {
	Data   GetCapsulesByPositionsData %[1]cjson:"data"%[1]c
	Errors []GraphQLError             %[1]cjson:"errors"%[1]c
}

type GetCapsulesByPositionsData struct {
	Capsules []Capsule %[1]cjson:"capsules"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type capsulesClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CapsulesClient {
	return &capsulesClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Generate_Input_3DArray(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/input_3d_array/schema.graphql")
	query := loadQuery(t, schema, "test/input_3d_array/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CapsulesClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Capsule struct {
	Id         string %[1]cjson:"id"%[1]c
}

type Limit struct {
	Size int %[1]cjson:"size"%[1]c
}

type Position struct {
	X int %[1]cjson:"x"%[1]c
	Y int %[1]cjson:"y"%[1]c
}

const getCapsulesByPositions = %[1]cquery GetCapsulesByPositions($find: [[[Position]]], $limit: [[[Limit]]], $selector: [[[String]]]) {
   capsules(find: $find, limit: $limit, selector: $selector) {
       id
   }
}%[1]c

type CapsulesClient interface {
	GetCapsulesByPositions(ctx context.Context, find [][][]Position, limit [][][]Limit, selector [][][]string, header *http.Header) (*http.Response, error)
}

func (c *capsulesClient) GetCapsulesByPositions(ctx context.Context, find [][][]Position, limit [][][]Limit, selector [][][]string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 3)
	params["find"] = find
	params["limit"] = limit
	params["selector"] = selector

	return c.ctrl.Execute(ctx, getCapsulesByPositions, params, header)
}

type GetCapsulesByPositionsResponse struct {
	Data   GetCapsulesByPositionsData %[1]cjson:"data"%[1]c
	Errors []GraphQLError             %[1]cjson:"errors"%[1]c
}

type GetCapsulesByPositionsData struct {
	Capsules []Capsule %[1]cjson:"capsules"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type capsulesClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CapsulesClient {
	return &capsulesClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_CircularType(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/circular_type/schema.graphql")
	query := loadQuery(t, schema, "test/circular_type/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "MovieClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Actor struct {
	ActedIn []Movie %[1]cjson:"actedIn"%[1]c
}

type Movie struct {
	Title string %[1]cjson:"title"%[1]c
	Actor Actor  %[1]cjson:"actor"%[1]c
}

const getAllMoviesWhereActorsOfTheMovieActedIn = %[1]cquery GetAllMoviesWhereActorsOfTheMovieActedIn($title: String!) {
    movie(title: $title) {
        actor {
            actedIn {
                title
            }
        }
    }
}%[1]c

type MovieClient interface {
	GetAllMoviesWhereActorsOfTheMovieActedIn(ctx context.Context, title string, header *http.Header) (*http.Response, error)
}

func (c *movieClient) GetAllMoviesWhereActorsOfTheMovieActedIn(ctx context.Context, title string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["title"] = title

	return c.ctrl.Execute(ctx, getAllMoviesWhereActorsOfTheMovieActedIn, params, header)
}

type GetAllMoviesWhereActorsOfTheMovieActedInResponse struct {
	Data   GetAllMoviesWhereActorsOfTheMovieActedInData %[1]cjson:"data"%[1]c
	Errors []GraphQLError                               %[1]cjson:"errors"%[1]c
}

type GetAllMoviesWhereActorsOfTheMovieActedInData struct {
	Movie Movie %[1]cjson:"movie"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type movieClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) MovieClient {
	return &movieClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_FragmentType(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/fragment/schema.graphql")
	query := loadQuery(t, schema, "test/fragment/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "RocketClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Rocket struct {
	Active      bool   %[1]cjson:"active"%[1]c
	Country     string %[1]cjson:"country"%[1]c
	Description string %[1]cjson:"description"%[1]c
	Id          string %[1]cjson:"id"%[1]c
	Name        string %[1]cjson:"name"%[1]c
}

const getShortRocketInfo = %[1]cquery GetShortRocketInfo {
    rockets {
        ...RocketShortInfo
    }
}
fragment RocketShortInfo on Rocket {
    id
    name
    description
    ...AdditionalRocketInfo
}
fragment AdditionalRocketInfo on Rocket {
    country
    ... on Rocket {
        ...InformatoryRocketInfo
    }
}
fragment InformatoryRocketInfo on Rocket {
    active
}%[1]c

type RocketClient interface {
	GetShortRocketInfo(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *rocketClient) GetShortRocketInfo(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getShortRocketInfo, params, header)
}

type GetShortRocketInfoResponse struct {
	Data   GetShortRocketInfoData %[1]cjson:"data"%[1]c
	Errors []GraphQLError         %[1]cjson:"errors"%[1]c
}

type GetShortRocketInfoData struct {
	Rockets []Rocket %[1]cjson:"rockets"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type rocketClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) RocketClient {
	return &rocketClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_SelectionSet(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/selection_set/schema.graphql")
	query := loadQuery(t, schema, "test/selection_set/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CountriesClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
    "context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Continent struct {
	Code string %[1]cjson:"code"%[1]c
	Name string %[1]cjson:"name"%[1]c
}

type Country struct {
	Code string %[1]cjson:"code"%[1]c
	Name string %[1]cjson:"name"%[1]c
}

const getCountriesAndContinents = %[1]cquery getCountriesAndContinents {
    continents {
        code
        name
    }
    countries {
        code
        name
    }
}%[1]c

type CountriesClient interface {
	GetCountriesAndContinents(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *countriesClient) GetCountriesAndContinents(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getCountriesAndContinents, params, header)
}

type GetCountriesAndContinentsResponse struct {
	Data   GetCountriesAndContinentsData %[1]cjson:"data"%[1]c
	Errors []GraphQLError                %[1]cjson:"errors"%[1]c
}

type GetCountriesAndContinentsData struct {
	Continents []Continent %[1]cjson:"continents"%[1]c
	Countries  []Country   %[1]cjson:"countries"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type countriesClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CountriesClient {
	return &countriesClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Interface(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/interface/schema.graphql")
	query := loadQuery(t, schema, "test/interface/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CharacterClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type CharacterFragment struct {
	Name            string %[1]cjson:"name"%[1]c
	HomePlanet      string %[1]cjson:"homePlanet"%[1]c
	PrimaryFunction string %[1]cjson:"primaryFunction"%[1]c
}

const getCharacters = %[1]cquery getCharacters {
    characters {
        ... on Human {
            ... on Human {
                homePlanet
            }
        }
        ... on Droid {
            primaryFunction
        }
        ... on Character {
            name
        }
    }
}%[1]c

type CharacterClient interface {
	GetCharacters(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *characterClient) GetCharacters(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getCharacters, params, header)
}

type GetCharactersResponse struct {
	Data   GetCharactersData %[1]cjson:"data"%[1]c
	Errors []GraphQLError    %[1]cjson:"errors"%[1]c
}

type GetCharactersData struct {
	Characters []CharacterFragment %[1]cjson:"characters"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type characterClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CharacterClient {
	return &characterClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Interface_3DArray(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/interface_3d_array/schema.graphql")
	query := loadQuery(t, schema, "test/interface_3d_array/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CharacterClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type CharacterFragment struct {
	Name            string %[1]cjson:"name"%[1]c
	HomePlanet      string %[1]cjson:"homePlanet"%[1]c
	PrimaryFunction string %[1]cjson:"primaryFunction"%[1]c
}

const getCharacters = %[1]cquery getCharacters {
    characters {
        ... on Human {
            homePlanet
        }
        ... on Droid {
            primaryFunction
        }
        ... on Character {
            name
        }
    }
}%[1]c

type CharacterClient interface {
	GetCharacters(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *characterClient) GetCharacters(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getCharacters, params, header)
}

type GetCharactersResponse struct {
	Data   GetCharactersData %[1]cjson:"data"%[1]c
	Errors []GraphQLError    %[1]cjson:"errors"%[1]c
}

type GetCharactersData struct {
	Characters [][][]CharacterFragment %[1]cjson:"characters"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type characterClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CharacterClient {
	return &characterClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Interface_No_Implementation(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/interface_no_impl/schema.graphql")
	query := loadQuery(t, schema, "test/interface_no_impl/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CharacterClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type CharacterFragment struct {
	Id string %[1]cjson:"id"%[1]c
}

type Response struct {
	SuperType CharacterFragment %[1]cjson:"superType"%[1]c
}

const getCharactersId = %[1]cquery getCharactersId {
    characters {
        superType {
            id
        }
    }
}%[1]c

type CharacterClient interface {
	GetCharactersId(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *characterClient) GetCharactersId(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getCharactersId, params, header)
}

type GetCharactersIdResponse struct {
	Data   GetCharactersIdData %[1]cjson:"data"%[1]c
	Errors []GraphQLError      %[1]cjson:"errors"%[1]c
}

type GetCharactersIdData struct {
	Characters Response %[1]cjson:"characters"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type characterClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CharacterClient {
	return &characterClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_InterfaceWithSelectionSet(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/interface_selection_set/schema.graphql")
	query := loadQuery(t, schema, "test/interface_selection_set/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "PlanetClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type CharacterFragment struct {
	Name            string %[1]cjson:"name"%[1]c
	HomePlanet      string %[1]cjson:"homePlanet"%[1]c
	PrimaryFunction string %[1]cjson:"primaryFunction"%[1]c
}

type PlanetFragment struct {
	Name        string %[1]cjson:"name"%[1]c
	Temperature int    %[1]cjson:"temperature"%[1]c
	Age         int    %[1]cjson:"age"%[1]c
}

const getCharacters = %[1]cquery getCharacters {
    characters {
        ... on Human {
            homePlanet
        }
        ... on Droid {
            primaryFunction
        }
        ... on Character {
            name
        }
    }
    planets {
        ... on IcePlanet {
            temperature
        }
        ... on RockyPlanet {
            age
        }
        ... on Planet {
            name
        }
    }
}%[1]c

type PlanetClient interface {
	GetCharacters(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *planetClient) GetCharacters(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getCharacters, params, header)
}

type GetCharactersResponse struct {
	Data   GetCharactersData %[1]cjson:"data"%[1]c
	Errors []GraphQLError    %[1]cjson:"errors"%[1]c
}

type GetCharactersData struct {
	Characters []CharacterFragment %[1]cjson:"characters"%[1]c
	Planets    []PlanetFragment    %[1]cjson:"planets"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type planetClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) PlanetClient {
	return &planetClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Union(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/union/schema.graphql")
	query := loadQuery(t, schema, "test/union/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CharacterClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type CharacterResultUnion struct {
	HomePlanet      string %[1]cjson:"homePlanet"%[1]c
	PrimaryFunction string %[1]cjson:"primaryFunction"%[1]c
}

const getCharacters = %[1]cquery getCharacters {
    characters {
        ... on Human {
            homePlanet
        }
        ... on Droid {
            primaryFunction
        }
    }
}%[1]c

type CharacterClient interface {
	GetCharacters(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *characterClient) GetCharacters(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getCharacters, params, header)
}

type GetCharactersResponse struct {
	Data   GetCharactersData %[1]cjson:"data"%[1]c
	Errors []GraphQLError    %[1]cjson:"errors"%[1]c
}

type GetCharactersData struct {
	Characters []CharacterResultUnion %[1]cjson:"characters"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type characterClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CharacterClient {
	return &characterClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_WithPointers(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/fragment/schema.graphql")
	query := loadQuery(t, schema, "test/fragment/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "RocketClient",
		UsePointers: true,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Rocket struct {
	Active      *bool   %[1]cjson:"active"%[1]c
	Country     *string %[1]cjson:"country"%[1]c
	Description *string %[1]cjson:"description"%[1]c
	Id          *string %[1]cjson:"id"%[1]c
	Name        *string %[1]cjson:"name"%[1]c
}

const getShortRocketInfo = %[1]cquery GetShortRocketInfo {
    rockets {
        ...RocketShortInfo
    }
}
fragment RocketShortInfo on Rocket {
    id
    name
    description
    ...AdditionalRocketInfo
}
fragment AdditionalRocketInfo on Rocket {
    country
    ... on Rocket {
        ...InformatoryRocketInfo
    }
}
fragment InformatoryRocketInfo on Rocket {
    active
}%[1]c

type RocketClient interface {
	GetShortRocketInfo(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *rocketClient) GetShortRocketInfo(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getShortRocketInfo, params, header)
}

type GetShortRocketInfoResponse struct {
	Data   *GetShortRocketInfoData %[1]cjson:"data"%[1]c
	Errors []GraphQLError         %[1]cjson:"errors"%[1]c
}

type GetShortRocketInfoData struct {
	Rockets []Rocket %[1]cjson:"rockets"%[1]c
}

type GraphQLError struct {
	Message    *string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions *GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   *int %[1]cjson:"line"%[1]c
	Column *int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code *string %[1]cjson:"code"%[1]c
}

type rocketClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) RocketClient {
	return &rocketClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_SubQuery(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/subquery/schema.graphql")
	query := loadQuery(t, schema, "test/subquery/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "GitClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type Collaborator struct {
	Id string %[1]cjson:"id"%[1]c
}

type PullRequest struct {
	BranchName    string         %[1]cjson:"branchName"%[1]c
	Collaborators []Collaborator %[1]cjson:"collaborators"%[1]c
}

type Repository struct {
	Name         string         %[1]cjson:"name"%[1]c
	Author       UserFragment   %[1]cjson:"author"%[1]c
	PullRequests []PullRequest  %[1]cjson:"pullRequests"%[1]c
	Users        []UserFragment %[1]cjson:"users"%[1]c
}

type UserFragment struct {
	Name         string %[1]cjson:"name"%[1]c
	TotalCommits int    %[1]cjson:"totalCommits"%[1]c
}

const getRepositoryInformation = %[1]cquery getRepositoryInformation {
    repositories(first: 10) {
        name
        author {
            ... on Author {
                name
            }
        }
        pullRequests(after: "ABC") {
            branchName
            collaborators(first: 5) {
                id
            }
        }
        users(before: "123") {
            totalCommits
        }
    }
}%[1]c

type GitClient interface {
	GetRepositoryInformation(ctx context.Context, header *http.Header) (*http.Response, error)
}

func (c *gitClient) GetRepositoryInformation(ctx context.Context, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(ctx, getRepositoryInformation, params, header)
}

type GetRepositoryInformationResponse struct {
	Data   GetRepositoryInformationData %[1]cjson:"data"%[1]c
	Errors []GraphQLError               %[1]cjson:"errors"%[1]c
}

type GetRepositoryInformationData struct {
	Repositories []Repository %[1]cjson:"repositories"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type gitClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) GitClient {
	return &gitClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func TestEvaluator_Comments(t *testing.T) {
	t.Parallel()
	schema := loadSchema(t, "test/comments/schema.graphql")
	query := loadQuery(t, schema, "test/comments/query.graphql")
	info := AdditionalInfo{
		PackageName: "grafik_client",
		ClientName:  "CommentsClient",
		UsePointers: false,
	}
	e := New(schema, query, info)

	out := getSourceString(t, e)
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with grafik. DO NOT EDIT

package grafik_client

import (
	"context"
	GraphqlClient "github.com/Bartosz-D3V/grafik/client"
	"net/http"
)

type File struct {
	Name string %[1]cjson:"name"%[1]c
}

const getFileNameWithId = %[1]cquery GetFileNameWithId($id: ID!) {
    getFile(id: $id) {
        name 
    }
}%[1]c

type CommentsClient interface {
	GetFileNameWithId(ctx context.Context, id string, header *http.Header) (*http.Response, error)
}

func (c *commentsClient) GetFileNameWithId(ctx context.Context, id string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["id"] = id

	return c.ctrl.Execute(ctx, getFileNameWithId, params, header)
}

type GetFileNameWithIdResponse struct {
	Data   GetFileNameWithIdData %[1]cjson:"data"%[1]c
	Errors []GraphQLError        %[1]cjson:"errors"%[1]c
}

type GetFileNameWithIdData struct {
	GetFile File %[1]cjson:"getFile"%[1]c
}

type GraphQLError struct {
	Message    string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation %[1]cjson:"locations"%[1]c
	Extensions GraphQLErrorExtensions %[1]cjson:"extensions"%[1]c
}

type GraphQLErrorLocation struct {
	Line   int %[1]cjson:"line"%[1]c
	Column int %[1]cjson:"column"%[1]c
}

type GraphQLErrorExtensions struct {
	Code string %[1]cjson:"code"%[1]c
}

type commentsClient struct {
	ctrl GraphqlClient.Client
}

func New(endpoint string, client *http.Client) CommentsClient {
	return &commentsClient{
		ctrl: GraphqlClient.New(endpoint, client),
	}
}
`, '`'))

	assert.Equal(t, expOut, out)
}

func loadSchema(t *testing.T, schemaName string) *ast.Schema {
	schemaLoc := path.Join("../", schemaName)
	file, err := ioutil.ReadFile(schemaLoc)
	assert.NoError(t, err)
	assert.NotNil(t, file)

	return gqlparser.MustLoadSchema(&ast.Source{
		Input: string(file),
		Name:  path.Base(schemaName),
	})
}

func loadQuery(t *testing.T, schema *ast.Schema, queryName string) *ast.QueryDocument {
	queryLoc := path.Join("../", queryName)
	file, err := ioutil.ReadFile(queryLoc)
	assert.NoError(t, err)
	assert.NotNil(t, file)

	return gqlparser.MustLoadQuery(schema, string(file))
}

func getSourceString(t *testing.T, e Evaluator) string {
	fileContent := e.Generate()
	src := &bytes.Buffer{}
	_, err := fileContent.WriteTo(src)
	assert.NoError(t, err)

	return src.String()
}
