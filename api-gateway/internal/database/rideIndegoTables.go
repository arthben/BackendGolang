package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

// Parameter for store data
type ParamStoreRideIndego struct {
	Master          *RideIndegoMaster
	Features        []*RideIndegoFeatures
	Properties      []*RideIndegoProperties
	PropertiesBikes []*RideIndegoBikes
}

type SearchResRideIndego struct {
	Master          *RideIndegoMaster
	Features        []*RideIndegoFeatures
	Properties      map[int]*RideIndegoProperties
	PropertiesBikes map[int][]*RideIndegoBikes
}

// Structure table rideindego_properties_bikes
type RideIndegoBikes struct {
	FetchID      string `db:"fetch_id"`
	FeatureID    int    `db:"feat_id"`
	PropertiesID int    `db:"id"`
	Index        int    `db:"idx"`
	Battery      int    `db:"battery"`
	DockNumber   int    `db:"dock_number"`
	IsElectric   bool   `db:"is_electric"`
	IsAvailable  bool   `db:"is_available"`
}

// Structure table rideindego_properties
type RideIndegoProperties struct {
	FetchID                string  `db:"fetch_id"`
	FeatureID              int     `db:"feat_id"`
	PropertiesID           int     `db:"id"`
	Name                   string  `db:"name"`
	Notes                  string  `db:"notes"`
	KioskID                int     `db:"kiosk_id"`
	EventEnd               string  `db:"event_end"`
	Latitude               float64 `db:"latitude"`
	OpenTime               string  `db:"open_time"`
	TimeZone               string  `db:"time_zone"`
	CloseTime              string  `db:"close_time"`
	IsVirtual              bool    `db:"is_virtual"`
	KioskType              int     `db:"kiosk_type"`
	Longitude              float64 `db:"longitude"`
	EventStart             string  `db:"event_start"`
	PublicText             string  `db:"public_text"`
	TotalDocks             int     `db:"total_docks"`
	AddressCity            string  `db:"address_city"`
	Coordinates            string  `db:"coordinates"`
	KioskStatus            string  `db:"kiosk_status"`
	AddressState           string  `db:"address_state"`
	IsEventBased           bool    `db:"is_event_based"`
	AddressStreet          string  `db:"address_street"`
	AddressZipCode         string  `db:"address_zip_code"`
	BikesAvailable         int     `db:"bikes_available"`
	DocksAvailable         int     `db:"docks_available"`
	TrikesAvailable        int     `db:"trikes_available"`
	KioskPublicStatus      string  `db:"kiosk_public_status"`
	SmartBikesAvailable    int     `db:"smart_bikes_available"`
	RewardBikesAvailable   int     `db:"reward_bikes_available"`
	RewardDocksAvailable   int     `db:"reward_docks_available"`
	ClassicBikesAvailable  int     `db:"classic_bikes_available"`
	KioskConnectionStatus  string  `db:"kiosk_connection_status"`
	ElectricBikesAvailable int     `db:"electric_bikes_available"`
}

// Structure table rideindego_features
type RideIndegoFeatures struct {
	FetchID            string `db:"fetch_id"`
	FeatureID          int    `db:"feat_id"`
	FeatureType        string `db:"ftype"`
	GeometryType       string `db:"geo_type"`
	GeometryCoordinate string `db:"geo_coordinate"`
}

// Structure table rideindego_master
type RideIndegoMaster struct {
	FetchID        string    `db:"fetch_id"`
	TypeCollection string    `db:"type_collection"`
	LastUpdate     time.Time `db:"last_update"`
}

// this struct only for store temporary
type RideIndegoMasterExtends struct {
	FetchID        string    `db:"fetch_id"`
	TypeCollection string    `db:"type_collection"`
	LastUpdate     time.Time `db:"last_update"`
	FeatureID      int       `db:"feat_id"`
}

type readRideIndego struct {
	db  *sqlx.DB
	ctx context.Context
}

func withKioskID(kioskID int) bool {
	return kioskID != -1
}

func withFeatureID(featureID int) bool {
	return featureID != -1
}

func (r *readRideIndego) readMaster(at time.Time, kioskID int, v interface{}) error {
	var (
		sql  string
		args []interface{}
	)

	// time parameter is mandatory
	args = append(args, at)

	if withKioskID(kioskID) {
		// find based on kioskID
		sql = `SELECT m.fetch_id, m.type_collection, m.last_update, p.feat_id
			   FROM rideindego_master m
			   LEFT OUTER JOIN rideindego_properties p on p.fetch_id=m.fetch_id 
			   WHERE m.last_update >= $1 and p.kiosk_id = $2
			   ORDER by m.last_update ASC
			   LIMIT 1`
		args = append(args, kioskID)

	} else {
		// find only based on time
		sql = `SELECT fetch_id, type_collection, last_update 
			   FROM rideindego_master
			   WHERE last_update >= $1
			   ORDER by last_update ASC
			   LIMIT 1`
	}

	err := r.db.GetContext(r.ctx, v, sql, args...)
	return err
}

func (r *readRideIndego) readFeatures(fetchID string, featureID int) ([]*RideIndegoFeatures, error) {
	var args []interface{}

	// fetchID is mandatory
	args = append(args, fetchID)

	sql := `SELECT fetch_id, feat_id, ftype, geo_type, ST_AsText(geo_coordinate) AS geo_coordinate
			FROM rideindego_features
			WHERE fetch_id=$1`

	if withFeatureID(featureID) {
		sql += " AND feat_id=$2"
		args = append(args, featureID)
	}

	var features []*RideIndegoFeatures
	err := r.db.SelectContext(r.ctx, &features, sql, args...)
	return features, err
}

func (r *readRideIndego) readProperties(fetchID string, featureID int) (map[int]*RideIndegoProperties, error) {
	var args []interface{}

	// fetchID is mandatory
	args = append(args, fetchID)

	sql := `SELECT fetch_id, feat_id, id, "name", notes, kiosk_id, event_end, 
			latitude, open_time, time_zone, close_time, is_virtual, kiosk_type, 
			longitude, event_start, public_text, total_docks, address_city, 
			ST_AsText(coordinates) AS coordinates, 
			kiosk_status, address_state, is_event_based, address_street, 
			address_zip_code, bikes_available, docks_available, trikes_available, 
			kiosk_public_status, smart_bikes_available, reward_bikes_available, 
			reward_docks_available, classic_bikes_available, kiosk_connection_status, 
			electric_bikes_available
			FROM rideindego_properties
			WHERE fetch_id=$1`

	if withFeatureID(featureID) {
		sql += " AND feat_id=$2"
		args = append(args, featureID)
	}

	rows, err := r.db.QueryxContext(r.ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	properties := make(map[int]*RideIndegoProperties, 0)
	for rows.Next() {
		prop := RideIndegoProperties{}
		err := rows.StructScan(&prop)
		if err != nil {
			return nil, err
		}

		properties[prop.FeatureID] = &prop
	}

	return properties, nil
}

func (r *readRideIndego) readPropBikes(fetchID string, featureID int) (map[int][]*RideIndegoBikes, error) {
	var args []interface{}

	// fetchID is mandatory
	args = append(args, fetchID)

	sql := `SELECT fetch_id, feat_id, id, idx, battery, dock_number, is_electric, is_available
			FROM rideindego_properties_bikes
			WHERE fetch_id=$1`

	if withFeatureID(featureID) {
		sql += " AND feat_id=$2"
		args = append(args, featureID)
	}

	rows, err := r.db.QueryxContext(r.ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bikes := make(map[int][]*RideIndegoBikes, 0)
	for rows.Next() {
		bike := RideIndegoBikes{}
		err := rows.StructScan(&bike)
		if err != nil {
			return nil, err
		}

		key := bike.FeatureID
		bikes[key] = append(bikes[key], &bike)
	}

	return bikes, nil
}

type storeRideIndeGo struct {
	tx  *sqlx.Tx
	ctx context.Context
}

func (s *storeRideIndeGo) insertMaster(master *RideIndegoMaster) error {
	sql := `INSERT INTO rideindego_master 
			(fetch_id, type_collection, last_update)
			VALUES
			(:fetch_id, :type_collection, :last_update)`

	_, err := s.tx.NamedExecContext(s.ctx, sql, master)
	return err
}

func (s *storeRideIndeGo) insertFeatures(features []*RideIndegoFeatures) error {
	sql := `INSERT INTO rideindego_features
			(fetch_id, feat_id, ftype, geo_type, geo_coordinate)
			VALUES
			(:fetch_id, :feat_id, :ftype, :geo_type, ST_GeomFromText(:geo_coordinate, 4326))`
	_, err := s.tx.NamedExecContext(s.ctx, sql, features)
	return err
}

func (s *storeRideIndeGo) insertProperties(properties []*RideIndegoProperties) error {
	sql := `INSERT INTO rideindego_properties
			(fetch_id, feat_id, id, "name", notes, kiosk_id, event_end, latitude, 
			 open_time, time_zone, close_time, is_virtual, kiosk_type, longitude, 
			 event_start, public_text, total_docks, address_city, coordinates, 
			 kiosk_status, address_state, is_event_based, address_street, 
			 address_zip_code, bikes_available, docks_available, trikes_available, 
			 kiosk_public_status, smart_bikes_available, reward_bikes_available, 
			 reward_docks_available, classic_bikes_available, kiosk_connection_status, 
			 electric_bikes_available)
			VALUES
			(:fetch_id, :feat_id, :id, :name, :notes, :kiosk_id, :event_end, :latitude, 
			 :open_time, :time_zone, :close_time, :is_virtual, :kiosk_type, :longitude, 
			 :event_start, :public_text, :total_docks, :address_city, 
			 ST_GeomFromText(:coordinates, 4326), 
			 :kiosk_status, :address_state, :is_event_based, :address_street, 
			 :address_zip_code, :bikes_available, :docks_available, :trikes_available, 
			 :kiosk_public_status, :smart_bikes_available, :reward_bikes_available, 
			 :reward_docks_available, :classic_bikes_available, :kiosk_connection_status, 
			 :electric_bikes_available)`
	_, err := s.tx.NamedExecContext(s.ctx, sql, properties)
	return err
}

func (s *storeRideIndeGo) insertPropertiesBike(propBikes []*RideIndegoBikes) error {
	sql := `INSERT INTO rideindego_properties_bikes
			(fetch_id, feat_id, id, battery, dock_number, is_electric, is_available)
			VALUES
			(:fetch_id, :feat_id, :id, :battery, :dock_number, :is_electric, :is_available)`
	_, err := s.tx.NamedExecContext(s.ctx, sql, propBikes)
	return err
}
