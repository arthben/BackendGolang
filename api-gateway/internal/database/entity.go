package database

import (
	"encoding/json"
	"time"
)

type TableIndegoWeather struct {
	LastUpdate time.Time       `db:"last_update"`
	RawData    json.RawMessage `db:"raw_data"`
}
