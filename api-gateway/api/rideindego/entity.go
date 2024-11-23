package rideindego

import (
	"time"
)

type Bikes struct {
	Battery     int  `json:"battery"`
	DockNumber  int  `json:"dockNumber"`
	IsElectric  bool `json:"isElectric"`
	IsAvailable bool `json:"isAvailable"`
}

type Properties struct {
	ID                     int       `json:"id"`
	Name                   string    `json:"name"`
	Bikes                  []Bikes   `json:"bikes"`
	Notes                  string    `json:"notes"`
	KioskID                int       `json:"kioskId"`
	EventEnd               string    `json:"eventEnd"`
	Latitude               float64   `json:"latitude"`
	OpenTime               string    `json:"openTime"`
	TimeZone               string    `json:"timeZone"`
	CloseTime              string    `json:"closeTime"`
	IsVirtual              bool      `json:"isVirtual"`
	KioskType              int       `json:"kioskType"`
	Longitude              float64   `json:"longitude"`
	EventStart             string    `json:"eventStart"`
	PublicText             string    `json:"publicText"`
	TotalDocks             int       `json:"totalDocks"`
	AddressCity            string    `json:"addressCity"`
	Coordinates            []float64 `json:"coordinates"`
	KioskStatus            string    `json:"kioskStatus"`
	AddressState           string    `json:"addressState"`
	IsEventBased           bool      `json:"isEventBased"`
	AddressStreet          string    `json:"addressStreet"`
	AddressZipCode         string    `json:"addressZipCode"`
	BikesAvailable         int       `json:"bikesAvailable"`
	DocksAvailable         int       `json:"docksAvailable"`
	TrikesAvailable        int       `json:"trikesAvailable"`
	KioskPublicStatus      string    `json:"kioskPublicStatus"`
	SmartBikesAvailable    int       `json:"smartBikesAvailable"`
	RewardBikesAvailable   int       `json:"rewardBikesAvailable"`
	RewardDocksAvailable   int       `json:"rewardDocksAvailable"`
	ClassicBikesAvailable  int       `json:"classicBikesAvailable"`
	KioskConnectionStatus  string    `json:"kioskConnectionStatus"`
	ElectricBikesAvailable int       `json:"electricBikesAvailable"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Features struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type FetchResponse struct {
	Type        string     `json:"type"`
	Features    []Features `json:"features"`
	LastUpdated time.Time  `json:"last_updated"`
}
