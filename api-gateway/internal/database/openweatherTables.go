package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type SearchResOpenWeather struct {
	Master *OpenWeatherMaster
	Detail []*OpenWewatherDetail
}

type OpenWewatherDetail struct {
	FetchID     string `db:"fetch_id"`
	Index       int    `db:"idx"`
	Description string `db:"description"`
	Icon        string `db:"icon"`
	ID          int    `db:"id"`
	Main        string `db:"main"`
}

type OpenWeatherMaster struct {
	FetchID       string  `db:"fetch_id"`
	Base          string  `db:"base"`
	Clouds        int     `db:"clouds"`
	COD           int     `db:"cod"`
	Coord         string  `db:"coord"`
	DT            int     `db:"dt"`
	ID            int     `db:"id"`
	MainFeelsLike float64 `db:"main_feels_like"`
	MainGrndLevel int     `db:"main_grnd_level"`
	MainHumidity  int     `db:"main_humidity"`
	MainPressure  int     `db:"main_pressure"`
	MainSeaLevel  int     `db:"main_sea_level"`
	MainTemp      float64 `db:"main_temp"`
	MainTempMax   float64 `db:"main_temp_max"`
	MainTempMin   float64 `db:"main_temp_min"`
	Name          string  `db:"name"`
	RainOneHour   float64 `db:"rain_one_hour"`
	SysCountry    string  `db:"sys_country"`
	SysID         int     `db:"sys_id"`
	SysSunrise    int     `db:"sys_sunrise"`
	SysSunset     int     `db:"sys_sunset"`
	SysType       int     `db:"sys_type"`
	Timezone      int     `db:"timezone"`
	Visibility    int     `db:"visibility"`
	WindDeg       int     `db:"wind_deg"`
	WindSpeed     float64 `db:"wind_speed"`
}

type readOpenWeather struct {
	db  *sqlx.DB
	ctx context.Context
}

func (r *readOpenWeather) readMaster(at time.Time) (*OpenWeatherMaster, error) {
	sql := `SELECT fetch_id, base, clouds, cod, ST_AsText(coord) AS coord, dt, id, 
			main_feels_like, main_grnd_level, main_humidity, main_pressure, 
			main_sea_level, main_temp, main_temp_max, main_temp_min, "name", 
			rain_one_hour, sys_country, sys_id, sys_sunrise, sys_sunset, 
			sys_type, timezone, visibility, wind_deg, wind_speed
			FROM openweather_master
			WHERE dt >= $1
			LIMIT 1`

	var master OpenWeatherMaster
	err := r.db.GetContext(r.ctx, &master, sql, at.Unix())
	return &master, err
}

func (r *readOpenWeather) readDetail(fetchID string) ([]*OpenWewatherDetail, error) {
	sql := `SELECT fetch_id, idx, description, icon, id, main
			FROM openweather_weather
			WHERE fetch_id = $1`

	var detail []*OpenWewatherDetail
	err := r.db.SelectContext(r.ctx, &detail, sql, fetchID)
	return detail, err
}

type storeWeather struct {
	tx  *sqlx.Tx
	ctx context.Context
}

func (s *storeWeather) insertMaster(master *OpenWeatherMaster) error {
	sql := `INSERT INTO openweather_master
			(fetch_id, base, clouds, cod, coord, dt, id, main_feels_like, 
			 main_grnd_level, main_humidity, main_pressure, main_sea_level, 
			 main_temp, main_temp_max, main_temp_min, "name", rain_one_hour, 
			 sys_country, sys_id, sys_sunrise, sys_sunset, sys_type, timezone, 
			 visibility, wind_deg, wind_speed)
			VALUES
			(:fetch_id, :base,:clouds, :cod, ST_GeomFromText(:coord, 4326), 
			 :dt, :id, :main_feels_like, 
			 :main_grnd_level, :main_humidity, :main_pressure, :main_sea_level, 
			 :main_temp, :main_temp_max, :main_temp_min, :name, :rain_one_hour, 
			 :sys_country, :sys_id, :sys_sunrise, :sys_sunset, :sys_type, :timezone, 
			 :visibility, :wind_deg, :wind_speed)`
	_, err := s.tx.NamedExecContext(s.ctx, sql, master)
	return err
}

func (s *storeWeather) insertDetail(detail []*OpenWewatherDetail) error {
	sql := `INSERT INTO openweather_weather
			(fetch_id, description, icon, id, main)
			VALUES
			(:fetch_id, :description, :icon, :id, :main)`
	_, err := s.tx.NamedExecContext(s.ctx, sql, detail)
	return err
}
