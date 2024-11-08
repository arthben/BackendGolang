package rideindego

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/arthben/BackendGolang/api-gateway/internal/client"
	"github.com/arthben/BackendGolang/api-gateway/internal/database"
)

const (
	Timeout          = 10 * time.Second
	LastUpdateLayout = "2006-01-02T15:04:05.999Z"
	BaseURL          = "https://www.rideindego.com/stations/json/"
)

type Service struct {
	db database.DBService
}

func NewService(db database.DBService) *Service {
	return &Service{db: db}
}

func (r *Service) FetchAndStore() (int, error) {

	// fetch data
	rawData, httpStatus, err := client.SendHTTPRequest(http.MethodGet, Timeout, nil, BaseURL)
	if err != nil {
		return httpStatus, err
	}

	// get json data
	var jsonData map[string]interface{}
	err = client.GetJSON(rawData, &jsonData)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error reading response body")
	}

	// save to database
	lastUpdate, err := time.Parse(LastUpdateLayout, jsonData["last_updated"].(string))
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error parsing data")
	}

	tbl := &database.TableIndegoWeather{
		LastUpdate: lastUpdate,
		RawData:    rawData,
	}
	err = r.db.SaveRideIndego(context.Background(), tbl)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error while store RideIndego data")
	}

	return http.StatusOK, nil
}

func (r *Service) Find(at time.Time, kioskId string) (string, json.RawMessage, int, error) {
	tbl, err := r.db.FindRideIndego(context.Background(), at, kioskId)
	if err != nil || tbl == nil {
		return "", nil, http.StatusNotFound, errors.New("No Data Found")
	}

	return tbl.LastUpdate.UTC().Format(LastUpdateLayout), tbl.RawData, http.StatusOK, nil

}
