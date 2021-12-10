package evaluator

import (
	"fmt"
	"github.com/Bartosz-D3V/ggrafik/test"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"io/ioutil"
	"path"
	"testing"
)

func TestEvaluator_Generate_FlatStructure(t *testing.T) {
	schema := loadSchema(t, "test/simple_type/simple_type.graphql")
	query := loadQuery(t, schema, "test/simple_type/simple_type_query.graphql")
	info := AdditionalInfo{
		PackageName: "ggrafik_client",
		ClientName:  "FilesClient",
		UsePointers: false,
	}
	e := New("../", schema, query, info)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

package ggrafik_client

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type File struct {
	Name string %[1]cjson:"name"%[1]c
}

const getFileNameWithId = %[1]cquery GetFileNameWithId($id: ID!) {
    getFile(id: $id) {
        name
    }
}

%[1]c

const renameFileWithId = %[1]cmutation RenameFileWithId($id: ID!, $name: String!) {
    renameFile(id: $id, name: $name) {
        name
    }
}%[1]c

type FilesClient interface {
	GetFileNameWithId(id string, header *http.Header) (*http.Response, error)
	RenameFileWithId(id string, name string, header *http.Header) (*http.Response, error)
}

func (c *filesClient) GetFileNameWithId(id string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["id"] = id

	return c.ctrl.Execute(getFileNameWithId, params, header)
}

func (c *filesClient) RenameFileWithId(id string, name string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 2)
	params["id"] = id
	params["name"] = name

	return c.ctrl.Execute(renameFileWithId, params, header)
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
	schema := loadSchema(t, "test/array/array.graphql")
	query := loadQuery(t, schema, "test/array/array_query.graphql")
	info := AdditionalInfo{
		PackageName: "ggrafik_client",
		ClientName:  "FilmsClient",
		UsePointers: false,
	}
	e := New("../", schema, query, info)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

package ggrafik_client

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type FilmsConnection struct {
	Films []Film %[1]cjson:"films"%[1]c
}

type Film struct {
	Producers []string %[1]cjson:"producers"%[1]c
}

const getAllFilmsProducers = %[1]cquery GetAllFilmsProducers {
    allFilms {
        films {
            producers
        }
    }
}%[1]c

type FilmsClient interface {
	GetAllFilmsProducers(header *http.Header) (*http.Response, error)
}

func (c *filmsClient) GetAllFilmsProducers(header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(getAllFilmsProducers, params, header)
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

func TestEvaluator_Generate_NestedStructure(t *testing.T) {
	schema := loadSchema(t, "test/nested_type/nested_type.graphql")
	query := loadQuery(t, schema, "test/nested_type/nested_type_query.graphql")
	info := AdditionalInfo{
		PackageName: "ggrafik_client",
		ClientName:  "SpecificHeroClient",
		UsePointers: false,
	}
	e := New("../", schema, query, info)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

package ggrafik_client

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type Character struct {
	HomeWorld Planet  %[1]cjson:"homeWorld"%[1]c
	Species   Species %[1]cjson:"species"%[1]c
}

type Planet struct {
	Location Location %[1]cjson:"location"%[1]c
}

type Location struct {
	PosX int %[1]cjson:"posX"%[1]c
	PoxY int %[1]cjson:"poxY"%[1]c
}

type Species struct {
	Origin Planet %[1]cjson:"origin"%[1]c
}

type CharacterSelector struct {
	IdSelector IdSelector %[1]cjson:"idSelector"%[1]c
}

type IdSelector struct {
	Id string %[1]cjson:"id"%[1]c
}

const getHeroWithId123ABC = %[1]cquery GetHeroWithId123ABC {
    getHero(characterSelector: {idSelector: {id: "123ABC"}}) {
        homeWorld {
            location {
                posX
            }
        }
    }
}%[1]c

type SpecificHeroClient interface {
	GetHeroWithId123ABC(header *http.Header) (*http.Response, error)
}

func (c *specificHeroClient) GetHeroWithId123ABC(header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(getHeroWithId123ABC, params, header)
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
	schema := loadSchema(t, "test/enum/enum.graphql")
	query := loadQuery(t, schema, "test/enum/enum_query.graphql")
	info := AdditionalInfo{
		PackageName: "ggrafik_client",
		ClientName:  "CompanyClient",
		UsePointers: false,
	}
	e := New("../", schema, query, info)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

package ggrafik_client

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type DepartmentName string

const (
	IT      DepartmentName = "IT"
	SALES   DepartmentName = "SALES"
	HR      DepartmentName = "HR"
	SUPPORT DepartmentName = "SUPPORT"
)

type Department struct {
	Name DepartmentName %[1]cjson:"name"%[1]c
}

const getDepartment = %[1]cquery getDepartment {
    getDepartment {
        name
    }
}%[1]c

type CompanyClient interface {
	GetDepartment(header *http.Header) (*http.Response, error)
}

func (c *companyClient) GetDepartment(header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(getDepartment, params, header)
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
	schema := loadSchema(t, "test/input/input.graphql")
	query := loadQuery(t, schema, "test/input/input_query.graphql")
	info := AdditionalInfo{
		PackageName: "ggrafik_client",
		ClientName:  "CapsulesClient",
		UsePointers: false,
	}
	e := New("../", schema, query, info)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

package ggrafik_client

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type CapsulesFind struct {
	Id             string %[1]cjson:"id"%[1]c
	Landings       int    %[1]cjson:"landings"%[1]c
	Mission        string %[1]cjson:"mission"%[1]c
	OriginalLaunch Date   %[1]cjson:"original_launch"%[1]c
	ReuseCount     int    %[1]cjson:"reuse_count"%[1]c
	Status         string %[1]cjson:"status"%[1]c
	Type           string %[1]cjson:"type"%[1]c
}

type Capsule struct {
	Id         string %[1]cjson:"id"%[1]c
	Landings   int    %[1]cjson:"landings"%[1]c
	ReuseCount int    %[1]cjson:"reuse_count"%[1]c
	Status     string %[1]cjson:"status"%[1]c
	Type       string %[1]cjson:"type"%[1]c
}

const getCapsulesByFullSelector = %[1]cquery GetCapsulesByFullSelector($order: String, $mission: String, $originalLaunch: Date, $id: ID, $sort: String) {
    capsules(order: $order, find: {landings: 10, mission: $mission, original_launch: $originalLaunch, id: $id}, sort: $sort) {
        id
        type
    }
}%[1]c

type CapsulesClient interface {
	GetCapsulesByFullSelector(order string, mission string, originalLaunch Date, id string, sort string, header *http.Header) (*http.Response, error)
}

func (c *capsulesClient) GetCapsulesByFullSelector(order string, mission string, originalLaunch Date, id string, sort string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 5)
	params["order"] = order
	params["mission"] = mission
	params["originalLaunch"] = originalLaunch
	params["id"] = id
	params["sort"] = sort

	return c.ctrl.Execute(getCapsulesByFullSelector, params, header)
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

func TestEvaluator_CircularType(t *testing.T) {
	schema := loadSchema(t, "test/circular_type/circular_type.graphql")
	query := loadQuery(t, schema, "test/circular_type/circular_type_query.graphql")
	info := AdditionalInfo{
		PackageName: "ggrafik_client",
		ClientName:  "MovieClient",
		UsePointers: false,
	}
	e := New("../", schema, query, info)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

package ggrafik_client

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type Movie struct {
	Title string %[1]cjson:"title"%[1]c
	Actor Actor  %[1]cjson:"actor"%[1]c
}

type Actor struct {
	Name    string  %[1]cjson:"name"%[1]c
	ActedIn []Movie %[1]cjson:"actedIn"%[1]c
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
	GetAllMoviesWhereActorsOfTheMovieActedIn(title string, header *http.Header) (*http.Response, error)
}

func (c *movieClient) GetAllMoviesWhereActorsOfTheMovieActedIn(title string, header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 1)
	params["title"] = title

	return c.ctrl.Execute(getAllMoviesWhereActorsOfTheMovieActedIn, params, header)
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
	schema := loadSchema(t, "test/fragment/fragment.graphql")
	query := loadQuery(t, schema, "test/fragment/fragment_query.graphql")
	info := AdditionalInfo{
		PackageName: "ggrafik_client",
		ClientName:  "RocketClient",
		UsePointers: false,
	}
	e := New("../", schema, query, info)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

package ggrafik_client

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type Rocket struct {
	Active         bool   %[1]cjson:"active"%[1]c
	Boosters       int    %[1]cjson:"boosters"%[1]c
	Company        string %[1]cjson:"company"%[1]c
	CostPerLaunch  int    %[1]cjson:"costPerLaunch"%[1]c
	Country        string %[1]cjson:"country"%[1]c
	Description    string %[1]cjson:"description"%[1]c
	Id             string %[1]cjson:"id"%[1]c
	Name           string %[1]cjson:"name"%[1]c
	Stages         int    %[1]cjson:"stages"%[1]c
	SuccessRatePct int    %[1]cjson:"successRatePct"%[1]c
	Type           string %[1]cjson:"type"%[1]c
	Wikipedia      string %[1]cjson:"wikipedia"%[1]c
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
}%[1]c

type RocketClient interface {
	GetShortRocketInfo(header *http.Header) (*http.Response, error)
}

func (c *rocketClient) GetShortRocketInfo(header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(getShortRocketInfo, params, header)
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

func TestEvaluator_WithPointers(t *testing.T) {
	schema := loadSchema(t, "test/fragment/fragment.graphql")
	query := loadQuery(t, schema, "test/fragment/fragment_query.graphql")
	info := AdditionalInfo{
		PackageName: "ggrafik_client",
		ClientName:  "RocketClient",
		UsePointers: true,
	}
	e := New("../", schema, query, info)

	out := string(e.Generate())
	expOut := test.PrepExpCode(t, fmt.Sprintf(`
// Generated with ggrafik. DO NOT EDIT

package ggrafik_client

import (
	GraphqlClient "github.com/Bartosz-D3V/ggrafik/client"
	"net/http"
)

type Rocket struct {
	Active         *bool   %[1]cjson:"active"%[1]c
	Boosters       *int    %[1]cjson:"boosters"%[1]c
	Company        *string %[1]cjson:"company"%[1]c
	CostPerLaunch  *int    %[1]cjson:"costPerLaunch"%[1]c
	Country        *string %[1]cjson:"country"%[1]c
	Description    *string %[1]cjson:"description"%[1]c
	Id             *string %[1]cjson:"id"%[1]c
	Name           *string %[1]cjson:"name"%[1]c
	Stages         *int    %[1]cjson:"stages"%[1]c
	SuccessRatePct *int    %[1]cjson:"successRatePct"%[1]c
	Type           *string %[1]cjson:"type"%[1]c
	Wikipedia      *string %[1]cjson:"wikipedia"%[1]c
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
}%[1]c

type RocketClient interface {
	GetShortRocketInfo(header *http.Header) (*http.Response, error)
}

func (c *rocketClient) GetShortRocketInfo(header *http.Header) (*http.Response, error) {
	params := make(map[string]interface{}, 0)

	return c.ctrl.Execute(getShortRocketInfo, params, header)
}

type GetShortRocketInfoResponse struct {
	Data   *GetShortRocketInfoData %[1]cjson:"data"%[1]c
	Errors []GraphQLError          %[1]cjson:"errors"%[1]c
}

type GetShortRocketInfoData struct {
	Rockets []Rocket %[1]cjson:"rockets"%[1]c
}

type GraphQLError struct {
	Message    *string                 %[1]cjson:"message"%[1]c
	Locations  []GraphQLErrorLocation  %[1]cjson:"locations"%[1]c
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
