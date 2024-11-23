package openweather

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/arthben/BackendGolang/api-gateway/internal/client"
	"github.com/arthben/BackendGolang/api-gateway/internal/config"
	dbase "github.com/arthben/BackendGolang/api-gateway/internal/database"
	"github.com/google/uuid"
)

const (
	Timeout = 10 * time.Second
	BaseURL = "https://api.openweathermap.org/data/2.5/weather"
)

type Service struct {
	db  dbase.DBService
	cfg *config.EnvParams
}

func NewService(db dbase.DBService, cfg *config.EnvParams) *Service {
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
	var jsonData FetchResponse
	err = client.GetJSON(rawData, &jsonData)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error reading response body")
	}

	// save to database
	err = o.storeToDB(&jsonData)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error while store Openweather data")
	}

	return http.StatusOK, nil
}

func (o *Service) storeToDB(fetchResponse *FetchResponse) error {
	fetchID := uuid.NewString()

	// save master
	master := &dbase.OpenWeatherMaster{
		FetchID:       fetchID,
		Base:          fetchResponse.Base,
		Clouds:        fetchResponse.Clouds.All,
		COD:           fetchResponse.Cod,
		Coord:         fmt.Sprintf("POINT(%f %f)", fetchResponse.Coord.Lat, fetchResponse.Coord.Lon),
		DT:            fetchResponse.Dt,
		ID:            fetchResponse.ID,
		MainFeelsLike: fetchResponse.Main.FeelsLike,
		MainGrndLevel: fetchResponse.Main.GrndLevel,
		MainHumidity:  fetchResponse.Main.Humidity,
		MainPressure:  fetchResponse.Main.Pressure,
		MainSeaLevel:  fetchResponse.Main.SeaLevel,
		MainTemp:      fetchResponse.Main.Temp,
		MainTempMax:   fetchResponse.Main.TempMax,
		MainTempMin:   fetchResponse.Main.TempMin,
		Name:          fetchResponse.Name,
		RainOneHour:   fetchResponse.Rain.OneH,
		SysCountry:    fetchResponse.Sys.Country,
		SysID:         fetchResponse.Sys.ID,
		SysSunrise:    fetchResponse.Sys.Sunrise,
		SysSunset:     fetchResponse.Sys.Sunset,
		SysType:       fetchResponse.Sys.Type,
		Timezone:      fetchResponse.Timezone,
		Visibility:    fetchResponse.Visibility,
		WindDeg:       fetchResponse.Wind.Deg,
		WindSpeed:     fetchResponse.Wind.Speed,
	}

	var details []*dbase.OpenWewatherDetail

	for i, detail := range fetchResponse.Weather {
		details = append(details, &dbase.OpenWewatherDetail{
			FetchID:     fetchID,
			Index:       i,
			Description: detail.Description,
			Icon:        detail.Icon,
			ID:          detail.ID,
			Main:        detail.Main,
		})
	}

	err := o.db.StoreOpenWeather(context.Background(), master, details)
	return err
}

func (o *Service) Search(at time.Time) (*FetchResponse, int, error) {
	tbl, err := o.db.SearchOpenWeather(context.Background(), at)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	// compose return value
	weathers := make([]Weather, 0)
	for _, weather := range tbl.Detail {
		weathers = append(weathers, Weather{
			Description: weather.Description,
			Icon:        weather.Icon,
			ID:          weather.ID,
			Main:        weather.Main,
		})
	}

	resp := FetchResponse{
		Base: tbl.Master.Base,
		Clouds: Clouds{
			All: tbl.Master.Clouds,
		},
		Cod:   tbl.Master.COD,
		Coord: unmarshalCoordinate(tbl.Master.Coord),
		Dt:    tbl.Master.DT,
		ID:    tbl.Master.ID,
		Main: Main{
			FeelsLike: tbl.Master.MainFeelsLike,
			GrndLevel: tbl.Master.MainGrndLevel,
			Humidity:  tbl.Master.MainHumidity,
			Pressure:  tbl.Master.MainPressure,
			SeaLevel:  tbl.Master.MainSeaLevel,
			Temp:      tbl.Master.MainTemp,
			TempMax:   tbl.Master.MainTempMax,
			TempMin:   tbl.Master.MainTempMin,
		},
		Name: tbl.Master.Name,
		Rain: Rain{
			OneH: tbl.Master.RainOneHour,
		},
		Sys: Sys{
			Country: tbl.Master.SysCountry,
			ID:      tbl.Master.SysID,
			Sunrise: tbl.Master.SysSunrise,
			Sunset:  tbl.Master.SysSunset,
			Type:    tbl.Master.SysType,
		},
		Timezone:   tbl.Master.Timezone,
		Visibility: tbl.Master.Visibility,
		Weather:    weathers,
		Wind: Wind{
			Deg:   tbl.Master.WindDeg,
			Speed: tbl.Master.WindSpeed,
		},
	}
	return &resp, http.StatusOK, nil
}

func unmarshalCoordinate(coordinate string) Coord {
	var lon, lat float64

	_, err := fmt.Sscanf(coordinate, "POINT(%f %f)", &lon, &lat)
	if err != nil {
		return Coord{}
	}

	return Coord{
		Lat: lat,
		Lon: lon,
	}
}
