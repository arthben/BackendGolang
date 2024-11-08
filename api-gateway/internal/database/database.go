package database

import (
	"context"
	"fmt"
	"time"

	"github.com/arthben/BackendGolang/api-gateway/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type DBService interface {
	Close() error
	SaveRideIndego(ctx context.Context, tbl *TableIndegoWeather) (err error)
	SaveOpenWeather(ctx context.Context, tbl *TableIndegoWeather) (err error)
	FindRideIndego(ctx context.Context, at time.Time, kioskId string) (*TableIndegoWeather, error)
	FindOpenWeather(ctx context.Context, at time.Time) (*TableIndegoWeather, error)
}

type dbase struct {
	db *sqlx.DB
}

func NewPool(cfg *config.EnvParams) (DBService, error) {
	// create connection and maintain pool internaly
	fmt.Printf("cfg.DB.DSN: %v\n", cfg.DB.DSN)

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

func (d *dbase) FindOpenWeather(ctx context.Context, at time.Time) (*TableIndegoWeather, error) {
	sql := `SELECT last_update, raw_data FROM openweather WHERE last_update >= $1`
	resp := TableIndegoWeather{}
	err := d.db.GetContext(ctx, &resp, sql, at)
	return &resp, err
}

func (d *dbase) FindRideIndego(ctx context.Context, at time.Time, kioskId string) (*TableIndegoWeather, error) {
	var sql string

	if len(kioskId) > 0 {
		sql = `SELECT last_update, 
			   jsonb_path_query(raw_data, '$.features[*] ? (@.properties.kioskId == ` + kioskId + `)') AS raw_data `
	} else {
		sql = "SELECT last_update, raw_data"
	}

	sql += " FROM rideindego WHERE last_update >= $1 ORDER BY last_update ASC LIMIT 1"

	resp := TableIndegoWeather{}
	err := d.db.GetContext(ctx, &resp, sql, at)
	return &resp, err
}

func (d *dbase) SaveRideIndego(ctx context.Context, tbl *TableIndegoWeather) (err error) {
	sql := "INSERT INTO rideindego(last_update, raw_data) VALUES(:last_update, :raw_data)"
	return d.insert(ctx, tbl, sql)
}

func (d *dbase) SaveOpenWeather(ctx context.Context, tbl *TableIndegoWeather) (err error) {
	sql := "INSERT INTO openweather(last_update, raw_data) VALUES(:last_update, :raw_data)"
	return d.insert(ctx, tbl, sql)
}

func (d *dbase) insert(ctx context.Context, tbl *TableIndegoWeather, sql string) (err error) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	// any error will be rollback
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	prep, err := tx.PrepareNamedContext(ctx, sql)
	if err != nil {
		return
	}

	_, err = prep.ExecContext(ctx, tbl)
	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}
