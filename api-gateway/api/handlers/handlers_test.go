package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/arthben/BackendGolang/api-gateway/internal/config"
	"github.com/arthben/BackendGolang/api-gateway/internal/database"
	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	At string `json:"at"`
}

var (
	cfg    *config.EnvParams
	dbPool database.DBService
)

func init() {
	os.Chdir("../..")
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(0)
	}
	cfg = conf

	repo, err := database.NewPool(cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(0)
	}

	dbPool = repo
}

func setupRouter() *gin.Engine {
	h := NewHandler(dbPool, cfg)
	h.BuildHandler()
	return h.router
}

func TestStationsWithKioskId(t *testing.T) {
	ro := setupRouter()

	scenarios := []struct {
		name           string
		header         map[string]string
		paramAt        string
		paramKioskId   string
		expectedStatus int
	}{
		{
			name: "Success",
			header: map[string]string{
				"Authorization": "Bearer secret_token_static",
			},
			paramAt:        "2024-11-08T01:00:00Z",
			paramKioskId:   "3005",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Fail - KioskId not valid",
			header: map[string]string{
				"Authorization": "Bearer secret_token_static",
			},
			paramAt:        "2024-11-08T01:00:00Z",
			paramKioskId:   "00000",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Fail - Without Token",
			header:         map[string]string{},
			paramAt:        "2024-11-08T01:00:00Z",
			paramKioskId:   "3005",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail - Invalid Token",
			header: map[string]string{
				"Authorization": "Bearer should_be_secret_token_static",
			},
			paramAt:        "2024-11-08T01:00:00Z",
			paramKioskId:   "3005",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail - No Data Found",
			header: map[string]string{
				"Authorization": "Bearer secret_token_static",
			},
			paramAt:        "2024-12-08T01:00:00Z",
			paramKioskId:   "3005",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, ts := range scenarios {
		t.Run(ts.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/stations/"+ts.paramKioskId+"?at="+ts.paramAt, nil)
			if err != nil {
				t.Errorf("Expected status %d but error occured - %v", ts.expectedStatus, err.Error())
				return
			}

			for key, v := range ts.header {
				req.Header.Add(key, v)
			}

			w := httptest.NewRecorder()
			ro.ServeHTTP(w, req)

			if w.Code != ts.expectedStatus {
				t.Errorf("Expected status %d but got response %d", ts.expectedStatus, w.Code)
				return
			}

			if ts.expectedStatus == http.StatusOK {
				var res ResponseBody
				err = json.Unmarshal(w.Body.Bytes(), &res)
				if err != nil {
					t.Errorf("Expected status %d but error occured - %v", ts.expectedStatus, err.Error())
					return
				}

				expectedResponseBody := ResponseBody{At: ts.paramAt}
				if !reflect.DeepEqual(res, expectedResponseBody) {
					t.Errorf("Expected status %d but response body mismatch - %v", ts.expectedStatus, res.At)
					return
				}
			}
		})
	}
}

func TestStations(t *testing.T) {
	ro := setupRouter()

	scenarios := []struct {
		name           string
		header         map[string]string
		paramAt        string
		expectedStatus int
	}{
		{
			name: "Success",
			header: map[string]string{
				"Authorization": "Bearer secret_token_static",
			},
			paramAt:        "2024-11-08T01:00:00Z",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Fail - Without Token",
			header:         map[string]string{},
			paramAt:        "2024-11-08T01:00:00Z",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail - Invalid Token",
			header: map[string]string{
				"Authorization": "Bearer should_be_secret_token_static",
			},
			paramAt:        "2024-11-08T01:00:00Z",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail - No Data Found",
			header: map[string]string{
				"Authorization": "Bearer secret_token_static",
			},
			paramAt:        "2024-12-08T01:00:00Z",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, ts := range scenarios {
		t.Run(ts.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/stations?at="+ts.paramAt, nil)
			if err != nil {
				t.Errorf("Expected status %d but error occured - %v", ts.expectedStatus, err.Error())
				return
			}

			for key, v := range ts.header {
				req.Header.Add(key, v)
			}

			w := httptest.NewRecorder()
			ro.ServeHTTP(w, req)

			if w.Code != ts.expectedStatus {
				t.Errorf("Expected status %d but got response %d", ts.expectedStatus, w.Code)
				return
			}

			if ts.expectedStatus == http.StatusOK {
				var res ResponseBody
				err = json.Unmarshal(w.Body.Bytes(), &res)
				if err != nil {
					t.Errorf("Expected status %d but error occured - %v", ts.expectedStatus, err.Error())
					return
				}

				expectedResponseBody := ResponseBody{At: ts.paramAt}
				if !reflect.DeepEqual(res, expectedResponseBody) {
					t.Errorf("Expected status %d but response body mismatch - %v", ts.expectedStatus, res.At)
					return
				}
			}
		})
	}
}

func TestFetchAndStore(t *testing.T) {
	ro := setupRouter()

	scenarios := []struct {
		name           string
		header         map[string]string
		expectedStatus int
	}{
		{
			name: "Success",
			header: map[string]string{
				"Authorization": "Bearer secret_token_static",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Fail - Without Token",
			header:         map[string]string{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail - Invalid Token",
			header: map[string]string{
				"Authorization": "Bearer should_be_secret_token_static",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, ts := range scenarios {
		t.Run(ts.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/v1/indego-data-fetch-and-store-it-db", nil)
			if err != nil {
				t.Errorf("Expected status %d but error occured - %v", ts.expectedStatus, err.Error())
				return
			}

			for key, v := range ts.header {
				req.Header.Add(key, v)
			}

			w := httptest.NewRecorder()
			ro.ServeHTTP(w, req)

			if w.Code != ts.expectedStatus {
				t.Errorf("Expected status %d but got response %d", ts.expectedStatus, w.Code)
				return
			}

			if ts.expectedStatus == http.StatusOK {
				var res map[string]string
				err = json.Unmarshal(w.Body.Bytes(), &res)
				if err != nil {
					t.Errorf("Expected status %d but error occured - %v", ts.expectedStatus, err.Error())
					return
				}

				if res["status"] != "Fetch and store success" {
					t.Errorf("Expected status %d but response body mismatch - %v", ts.expectedStatus, res["status"])
					return
				}

			}
		})
	}
}
