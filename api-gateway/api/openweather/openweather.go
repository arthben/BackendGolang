package openweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/arthben/BackendGolang/api-gateway/internal/client"
	"github.com/arthben/BackendGolang/api-gateway/internal/config"
	"github.com/arthben/BackendGolang/api-gateway/internal/database"
)

const (
	Timeout = 10 * time.Second
	BaseURL = "https://api.openweathermap.org/data/2.5/weather"
)

type Service struct {
	db  database.DBService
	cfg *config.EnvParams
}

func NewService(db database.DBService, cfg *config.EnvParams) *Service {
	return &Service{db: db, cfg: cfg}
}

func (o *Service) FetchAndStore() (int, error) {
	params := url.Values{}
	params.Add("q", "Philadelphia")
	params.Add("appid", o.cfg.OpenWeather.APIKey)

	url := fmt.Sprintf("%s?%s", BaseURL, params.Encode())
	rawData, httpStatus, err := client.SendHTTPRequest(http.MethodGet, Timeout, nil, url)
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
	ts := int64(jsonData["dt"].(float64))
	lastUpdate := time.Unix(ts, 0).UTC()
	tbl := &database.TableIndegoWeather{
		LastUpdate: lastUpdate,
		RawData:    rawData,
	}
	err = o.db.SaveOpenWeather(context.Background(), tbl)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error while store OpenWeather data")
	}

	return http.StatusOK, nil
}

func (o *Service) Find(at time.Time, kioskId string) (json.RawMessage, int, error) {
	tbl, err := o.db.FindOpenWeather(context.Background(), at)
	if err != nil || tbl == nil {
		return nil, http.StatusNotFound, errors.New("No Data Found")
	}

	return tbl.RawData, http.StatusOK, nil
}
