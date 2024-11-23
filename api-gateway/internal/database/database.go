package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/arthben/BackendGolang/api-gateway/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

type DBService interface {
	Close() error
	StoreRideIndego(context.Context, ParamStoreRideIndego) error
	SearchRideIndego(ctx context.Context, at time.Time, kioskID string) (SearchResRideIndego, error)
	StoreOpenWeather(ctx context.Context, master *OpenWeatherMaster, details []*OpenWewatherDetail) error
	SearchOpenWeather(ctx context.Context, at time.Time) (SearchResOpenWeather, error)
}

type dbase struct {
	db *sqlx.DB
}

func NewPool(cfg *config.EnvParams) (DBService, error) {
	// create connection and maintain pool internaly
	db, err := sqlx.Connect("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxPool)
	db.SetMaxIdleConns(cfg.DB.MinPool)

	return &dbase{db: db}, nil
}

func (d *dbase) Close() error {
	return d.db.Close()
}

func (d *dbase) SearchOpenWeather(
	ctx context.Context,
	at time.Time,
) (searchResult SearchResOpenWeather, err error) {
	findme := readOpenWeather{db: d.db, ctx: ctx}
	searchResult.Master, err = findme.readMaster(at)
	if err != nil {
		handleError("readMaster", err)
	}

	searchResult.Detail, err = findme.readDetail(searchResult.Master.FetchID)
	if err != nil {
		handleError("readDetail", err)
	}

	return
}

func (d *dbase) SearchRideIndego(
	ctx context.Context,
	at time.Time,
	kioskID string,
) (searchResult SearchResRideIndego, err error) {

	var (
		kiosk     int = -1
		featureID int = -1
	)

	findme := readRideIndego{db: d.db, ctx: ctx}

	if len(kioskID) > 0 {
		if kiosk, err = strconv.Atoi(kioskID); err != nil {
			err = errors.New("invalid kioskID")
			return
		}

		masterExtended := RideIndegoMasterExtends{}
		if err = findme.readMaster(at, kiosk, &masterExtended); err != nil {
			handleError("readMaster", err)
			return
		}

		featureID = masterExtended.FeatureID
		searchResult.Master = &RideIndegoMaster{
			FetchID:        masterExtended.FetchID,
			TypeCollection: masterExtended.TypeCollection,
			LastUpdate:     masterExtended.LastUpdate,
		}

	} else {
		master := RideIndegoMaster{}
		if err = findme.readMaster(at, kiosk, &master); err != nil {
			return
		}
		searchResult.Master = &master
	}

	// find features data
	searchResult.Features, err = findme.readFeatures(searchResult.Master.FetchID, featureID)
	if err != nil {
		handleError("readFeatures", err)
	}

	// find properties data
	searchResult.Properties, err = findme.readProperties(searchResult.Master.FetchID, featureID)
	if err != nil {
		handleError("readProperties", err)
	}

	// find properti bikes
	searchResult.PropertiesBikes, err = findme.readPropBikes(searchResult.Master.FetchID, featureID)
	if err != nil {
		handleError("readPropertiesBikes", err)
	}

	return
}

func (d *dbase) StoreRideIndego(ctx context.Context, pInput ParamStoreRideIndego) (err error) {

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	storeData := storeRideIndeGo{tx: tx, ctx: ctx}
	if err = storeData.insertMaster(pInput.Master); err != nil {
		handleError("insertMaster", err)
		return
	}
	if err = storeData.insertFeatures(pInput.Features); err != nil {
		handleError("insertFeatures", err)
		return
	}
	if err = storeData.insertProperties(pInput.Properties); err != nil {
		handleError("insertProperties", err)
		return
	}
	if err = storeData.insertPropertiesBike(pInput.PropertiesBikes); err != nil {
		handleError("insertPropertiesBike", err)
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}

	return nil
}

func handleError(msg string, err error) {
	log.Error().Err(err).Msg(msg)
	fmt.Printf("%v - err: %v\n", msg, err)
}

func (d *dbase) StoreOpenWeather(ctx context.Context, master *OpenWeatherMaster, detail []*OpenWewatherDetail) (err error) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	storeData := storeWeather{tx: tx, ctx: ctx}
	if err = storeData.insertMaster(master); err != nil {
		handleError("inserMaster", err)
		return
	}
	if err = storeData.insertDetail(detail); err != nil {
		handleError("insertDetail", err)
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return nil
}
