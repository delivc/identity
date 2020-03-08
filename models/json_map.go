package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONMap is a JSON object
type JSONMap map[string]interface{}

// Value returns the jsonmap as string
func (j JSONMap) Value() (driver.Value, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return driver.Value(""), err
	}
	return driver.Value(string(data)), nil
}

// Scan scans a json interface / map
func (j JSONMap) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	default:
		return errors.New("Invalid data type for JSONMap")
	}

	if len(source) == 0 {
		source = []byte("{}")
	}
	return json.Unmarshal(source, &j)
}
