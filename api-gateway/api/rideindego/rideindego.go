package rideindego

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/arthben/BackendGolang/api-gateway/internal/client"
	dbase "github.com/arthben/BackendGolang/api-gateway/internal/database"
	"github.com/google/uuid"
)

const (
	Timeout          = 10 * time.Second
	LastUpdateLayout = "2006-01-02T15:04:05.999Z"
	BaseURL          = "https://www.rideindego.com/stations/json/"
)

type Service struct {
	db dbase.DBService
}

func NewService(db dbase.DBService) *Service {
	return &Service{db: db}
}

func (r *Service) Search(at time.Time, kioskId string) (*FetchResponse, int, error) {
	tbl, err := r.db.SearchRideIndego(context.Background(), at, kioskId)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	// compose return value
	features := []Features{}
	for _, feature := range tbl.Features {
		fID := feature.FeatureID
		properties := parseProperties(tbl.Properties[fID], tbl.PropertiesBikes[fID])

		features = append(features, Features{
			Type: feature.FeatureType,
			Geometry: Geometry{
				Type:        feature.GeometryType,
				Coordinates: unmarshalCoordinate(feature.GeometryCoordinate),
			},
			Properties: properties,
		})
	}

	resp := FetchResponse{
		Type:        tbl.Master.TypeCollection,
		Features:    features,
		LastUpdated: tbl.Master.LastUpdate.UTC(),
	}

	return &resp, http.StatusOK, nil
}

func (r *Service) FetchAndStore() (int, error) {

	// fetch data
	rawData, httpStatus, err := client.SendHTTPRequest(http.MethodGet, Timeout, nil, BaseURL)
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
	err = r.storeToDB(&jsonData)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error while store RideIndego data")
	}
	return http.StatusOK, nil
}

func (r *Service) storeToDB(fetchResponse *FetchResponse) error {
	fetchID := uuid.NewString()
	featureID := int(fetchResponse.LastUpdated.Unix())

	// save master
	master := &dbase.RideIndegoMaster{
		FetchID:        fetchID,
		TypeCollection: fetchResponse.Type,
		LastUpdate:     fetchResponse.LastUpdated,
	}

	var (
		features   []*dbase.RideIndegoFeatures
		properties []*dbase.RideIndegoProperties
		propBikes  []*dbase.RideIndegoBikes
	)

	for i, val := range fetchResponse.Features {
		featureID += i

		// compose feature table
		features = append(features, &dbase.RideIndegoFeatures{
			FetchID:            fetchID,
			FeatureID:          featureID,
			FeatureType:        val.Type,
			GeometryType:       val.Geometry.Type,
			GeometryCoordinate: fmt.Sprintf("POINT(%f %f)", val.Geometry.Coordinates[0], val.Geometry.Coordinates[1]),
		})

		properties = append(properties, &dbase.RideIndegoProperties{
			FetchID:                fetchID,
			FeatureID:              featureID,
			PropertiesID:           val.Properties.ID,
			Name:                   val.Properties.Name,
			Notes:                  val.Properties.Notes,
			KioskID:                val.Properties.KioskID,
			EventEnd:               val.Properties.EventEnd,
			Latitude:               val.Properties.Latitude,
			OpenTime:               val.Properties.OpenTime,
			TimeZone:               val.Properties.TimeZone,
			CloseTime:              val.Properties.CloseTime,
			IsVirtual:              val.Properties.IsVirtual,
			KioskType:              val.Properties.KioskType,
			Longitude:              val.Properties.Longitude,
			EventStart:             val.Properties.EventStart,
			PublicText:             val.Properties.PublicText,
			TotalDocks:             val.Properties.TotalDocks,
			AddressCity:            val.Properties.AddressCity,
			Coordinates:            fmt.Sprintf("POINT(%f %f)", val.Properties.Coordinates[0], val.Properties.Coordinates[1]),
			KioskStatus:            val.Properties.KioskStatus,
			AddressState:           val.Properties.AddressState,
			IsEventBased:           val.Properties.IsEventBased,
			AddressStreet:          val.Properties.AddressStreet,
			AddressZipCode:         val.Properties.AddressZipCode,
			BikesAvailable:         val.Properties.BikesAvailable,
			DocksAvailable:         val.Properties.DocksAvailable,
			TrikesAvailable:        val.Properties.TrikesAvailable,
			KioskPublicStatus:      val.Properties.KioskPublicStatus,
			SmartBikesAvailable:    val.Properties.SmartBikesAvailable,
			RewardBikesAvailable:   val.Properties.RewardBikesAvailable,
			RewardDocksAvailable:   val.Properties.RewardDocksAvailable,
			ClassicBikesAvailable:  val.Properties.ClassicBikesAvailable,
			KioskConnectionStatus:  val.Properties.KioskConnectionStatus,
			ElectricBikesAvailable: val.Properties.ElectricBikesAvailable,
		})

		for j, bike := range val.Properties.Bikes {
			propBikes = append(propBikes, &dbase.RideIndegoBikes{
				FetchID:      fetchID,
				FeatureID:    featureID,
				PropertiesID: val.Properties.ID,
				Index:        j,
				Battery:      bike.Battery,
				DockNumber:   bike.DockNumber,
				IsElectric:   bike.IsElectric,
				IsAvailable:  bike.IsAvailable,
			})
		}
	}

	paramStoreData := dbase.ParamStoreRideIndego{
		Master:          master,
		Features:        features,
		Properties:      properties,
		PropertiesBikes: propBikes,
	}
	return r.db.StoreRideIndego(context.Background(), paramStoreData)
}

func parseProperties(prop *dbase.RideIndegoProperties, bikes []*dbase.RideIndegoBikes) Properties {

	// keep track each properties have valid Bikes array like fetch result
	// table Bikes have foreign reference from Properties table
	// key name is featureID. this key used on hashmap to avoid
	// query to database for each featureID
	arrBikes := make([]Bikes, 0)
	for _, bike := range bikes {
		arrBikes = append(arrBikes, Bikes{
			Battery:     bike.Battery,
			DockNumber:  bike.DockNumber,
			IsElectric:  bike.IsElectric,
			IsAvailable: bike.IsAvailable,
		})
	}

	return Properties{
		ID:                     prop.PropertiesID,
		Name:                   prop.Name,
		Bikes:                  arrBikes,
		Notes:                  prop.Notes,
		KioskID:                prop.KioskID,
		EventEnd:               prop.EventEnd,
		Latitude:               prop.Latitude,
		OpenTime:               prop.OpenTime,
		TimeZone:               prop.TimeZone,
		CloseTime:              prop.CloseTime,
		IsVirtual:              prop.IsVirtual,
		KioskType:              prop.KioskType,
		Longitude:              prop.Longitude,
		EventStart:             prop.EventStart,
		PublicText:             prop.PublicText,
		TotalDocks:             prop.TotalDocks,
		AddressCity:            prop.AddressCity,
		Coordinates:            unmarshalCoordinate(prop.Coordinates),
		KioskStatus:            prop.KioskStatus,
		AddressState:           prop.AddressState,
		IsEventBased:           prop.IsEventBased,
		AddressStreet:          prop.AddressStreet,
		AddressZipCode:         prop.AddressZipCode,
		BikesAvailable:         prop.BikesAvailable,
		DocksAvailable:         prop.DocksAvailable,
		TrikesAvailable:        prop.TrikesAvailable,
		KioskPublicStatus:      prop.KioskPublicStatus,
		SmartBikesAvailable:    prop.SmartBikesAvailable,
		RewardBikesAvailable:   prop.RewardBikesAvailable,
		RewardDocksAvailable:   prop.RewardDocksAvailable,
		ClassicBikesAvailable:  prop.ClassicBikesAvailable,
		KioskConnectionStatus:  prop.KioskConnectionStatus,
		ElectricBikesAvailable: prop.ElectricBikesAvailable,
	}
}

func unmarshalCoordinate(coordinate string) []float64 {
	var lon, lat float64

	_, err := fmt.Sscanf(coordinate, "POINT(%f %f)", &lon, &lat)
	if err != nil {
		return []float64{}
	}

	return []float64{lon, lat}
}
