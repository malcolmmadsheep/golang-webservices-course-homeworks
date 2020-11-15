package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

const ValidAccessToken = "ValidAccessToken"

type ServerUser struct {
	XMLName   xml.Name `xml:"row" json:"-"`
	ID        int      `xml:"id" json:"Id"`
	FirstName string   `xml:"last_name" json:"-"`
	LastName  string   `xml:"first_name" json:"-"`
	Age       int      `xml:"age" json:"Age"`
	About     string   `xml:"about" json:"About"`
	Gender    string   `xml:"gender" json:"Gender"`
	FullName  string   `json:"Name"`
}

type ServerUsers struct {
	XMLName xml.Name     `xml:"root"`
	Users   []ServerUser `xml:"row"`
}

type FindUserRequestParamsTestcase struct {
	ID                   string
	Request              SearchRequest
	ExpectedErrorMessage string
	IsError              bool
}

func handleError(rw http.ResponseWriter, status int, err error) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(fmt.Sprintf(`{"status":%d,"Error":"%s"}`, status, err)))
}

func SearchServer(rw http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("AccessToken")

	if accessToken != ValidAccessToken {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	orderField := r.URL.Query().Get("order_field")

	if orderField != "" && orderField != "Id" && orderField != "Age" && orderField != "Name" {
		handleError(rw, http.StatusBadRequest, fmt.Errorf("ErrorBadOrderField"))
		return
	}

	if orderField == "" {
		orderField = "Name"
	}

	orderBy := 1
	query := r.URL.Query().Get("query")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if r.URL.Query().Get("order_by") != "" {
		orderBy, err = strconv.Atoi(r.URL.Query().Get("order_by"))

		if err != nil {
			handleError(rw, http.StatusInternalServerError, err)
			return
		}
	}

	root := ServerUsers{}
	fileContent, err := ioutil.ReadFile("./dataset.xml")

	if err != nil {
		handleError(rw, http.StatusInternalServerError, err)
		return
	}

	err = xml.Unmarshal(fileContent, &root)

	if err != nil {
		handleError(rw, http.StatusInternalServerError, err)
		return
	}

	for i := 0; i < len(root.Users); i++ {
		root.Users[i].FullName = root.Users[i].FirstName + " " + root.Users[i].LastName
	}

	foundUsers := make([]ServerUser, 0, limit)

	if query == "" {
		foundUsers = append(foundUsers, root.Users[0:limit]...)
	} else {
		var user ServerUser

		for i := 0; i < len(root.Users) && len(foundUsers) < limit; i++ {
			user = root.Users[i]

			if strings.Contains(user.About, query) || strings.Contains(user.FullName, query) {
				foundUsers = append(foundUsers, user)
			}
		}
	}

	sort.SliceStable(foundUsers, func(i, j int) bool {
		userI := foundUsers[i]
		userJ := foundUsers[j]

		if orderField == "Name" {
			isBigger := userI.FullName >= userJ.FullName

			if orderBy == -1 {
				return !isBigger
			}

			return isBigger
		}

		if orderField == "Age" {
			return userI.Age*orderBy > userJ.Age*orderBy
		}

		return userI.ID*orderBy > userJ.ID*orderBy
	})

	foundUsersStr, err := json.Marshal(foundUsers)

	if err != nil {
		handleError(rw, http.StatusInternalServerError, err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(foundUsersStr)
}

// Functionality testing

func TestFindUserEmptyQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	searchClient := SearchClient{
		ValidAccessToken,
		ts.URL,
	}

	searchRequest := SearchRequest{
		Limit: 50,
	}

	res, _ := searchClient.FindUsers(searchRequest)

	if res == nil {
		t.Error("expected response, got nil")
	}

	if len(res.Users) != 25 {
		t.Errorf("expected 25 users, got %d", len(res.Users))
	}

	if !res.NextPage {
		t.Errorf("expected to have next page")
	}

	ts.Close()
}

func TestFindUserSmallAmountOfUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	searchClient := SearchClient{
		ValidAccessToken,
		ts.URL,
	}

	searchRequest := SearchRequest{
		Query: "Brooks",
		Limit: 10,
	}

	res, err := searchClient.FindUsers(searchRequest)

	if err != nil {
		t.Errorf("expected nil, got error '%s'", err)
	}

	if len(res.Users) != 1 {
		t.Errorf("expected single user in result, got %d", len(res.Users))
	}

	if res.NextPage {
		t.Errorf("expected not to have next page")
	}

	ts.Close()
}

// Function errors testing

func TestFindUserInvalidRequestParams(t *testing.T) {
	cases := []FindUserRequestParamsTestcase{
		{"negative offset", SearchRequest{Offset: -1}, "offset must be > 0", true},
		{"negative limit", SearchRequest{Limit: -1}, "limit must be > 0", true},
	}

	searchClient := SearchClient{}

	for _, testCase := range cases {
		resp, err := searchClient.FindUsers(testCase.Request)

		if resp != nil {
			t.Errorf("[%s] Expected nil, got response", testCase.ID)
		}

		if err == nil || err.Error() != testCase.ExpectedErrorMessage {
			t.Errorf(`[%s] Expected error "%s", received "%s"`, testCase.ID, err.Error(), testCase.ExpectedErrorMessage)
		}
	}
}

// Server error response testing
func TestFindUserUnauthorized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	searchClient := SearchClient{
		AccessToken: "INVALID_TOKEN",
		URL:         ts.URL,
	}

	searchRequest := SearchRequest{}

	resp, err := searchClient.FindUsers(searchRequest)

	if resp != nil {
		t.Errorf("Expected nil, got response")
	}

	if err == nil || err.Error() != "Bad AccessToken" {
		t.Errorf(`Expected error "Bad AccessToken", got "%s"`, err)
	}

	ts.Close()
}

func TestFindUserTimeoutError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(1050 * time.Millisecond)
	}))

	searchClient := SearchClient{
		ValidAccessToken,
		ts.URL,
	}

	searchQuery := SearchRequest{}

	resp, err := searchClient.FindUsers(searchQuery)

	if resp != nil {
		t.Errorf("Expected nil, got response")
	}

	if err == nil || !strings.HasPrefix(err.Error(), "timeout for") {
		t.Errorf("Expected timeout error, got %s", err)
	}

	ts.Close()
}

func TestFindUserUnknownError(t *testing.T) {
	searchClient := SearchClient{
		ValidAccessToken,
		"invalidProtocolScheme://search.me",
	}

	searchQuery := SearchRequest{}

	resp, err := searchClient.FindUsers(searchQuery)

	if resp != nil {
		t.Errorf("Expected nil, got response")
	}

	if err == nil || !strings.HasPrefix(err.Error(), "unknown error") {
		t.Errorf(`Expected unknown error, got "%s"`, err.Error())
	}
}

func TestFindUserInternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))

	searchClient := SearchClient{
		ValidAccessToken,
		ts.URL,
	}

	searchQuery := SearchRequest{}

	resp, err := searchClient.FindUsers(searchQuery)

	if resp != nil {
		t.Errorf("Expected nil, got response")
	}

	if err == nil || err.Error() != "SearchServer fatal error" {
		t.Errorf(`Expected "SearchServer fatal error", got "%s"`, err)
	}

	ts.Close()
}

func TestFindUserBadRequestInvalidResponseJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"status": 400, "err": "bad req error"`))
	}))

	searchClient := SearchClient{
		ValidAccessToken,
		ts.URL,
	}

	res, err := searchClient.FindUsers(SearchRequest{})

	if res != nil {
		t.Errorf("Expected nil, got response")
	}

	if err == nil || !strings.HasPrefix(err.Error(), "cant unpack error json") {
		t.Errorf(`Expected JSON unpack error, got "%s"`, err)
	}

	ts.Close()
}

func TestFindUserInvalidOrderField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	searchClient := SearchClient{
		ValidAccessToken,
		ts.URL,
	}

	searchRequest := SearchRequest{
		OrderField: "invalid_order_field",
	}

	res, err := searchClient.FindUsers(searchRequest)

	if res != nil {
		t.Errorf("Expected nil, got response")
	}

	if err == nil || err.Error() != "OrderField invalid_order_field invalid" {
		t.Errorf(`Expected ErrorBadOrderField error, got "%s"`, err)
	}

	ts.Close()
}

func TestFindUserUnknownBadRequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"status":400,"Error":"I do not like it"}`))
	}))

	searchClient := SearchClient{
		ValidAccessToken,
		ts.URL,
	}

	searchRequest := SearchRequest{}

	res, err := searchClient.FindUsers(searchRequest)

	if res != nil {
		t.Error("Expected nil, got response")
	}

	if err == nil || err.Error() != "unknown bad request error: I do not like it" {
		t.Errorf(`Expected unknown bad request error, got "%s"`, err)
	}

	ts.Close()
}

func TestFindUserInvalidResponseJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"users":[]}`))
	}))

	searchClient := SearchClient{
		URL: ts.URL,
	}

	searchRequest := SearchRequest{}

	res, err := searchClient.FindUsers(searchRequest)

	if res != nil {
		t.Error("Expected nil, got res")
	}

	if err == nil || !strings.HasPrefix(err.Error(), "cant unpack result json") {
		t.Errorf(`Expected "cant unpack result json" error, got %s`, err)
	}

	ts.Close()
}
