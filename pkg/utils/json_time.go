package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type JSONTime time.Time

// Format ke JSON (RFC3339 tanpa nano, dengan offset)
func (jt JSONTime) MarshalJSON() ([]byte, error) {
	t := time.Time(jt)
	formatted := t.Format("2006-01-02T15:04:05-07:00")
	// formatted := t.Format("02/01/2006 15:04:05")
	return []byte(`"` + formatted + `"`), nil
}

// Parse dari JSON ke time.Time
func (jt *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02T15:04:05-07:00", s)
	if err != nil {
		return err
	}
	*jt = JSONTime(t)
	return nil
}

// Supaya GORM bisa simpan ke MySQL
func (jt JSONTime) Value() (driver.Value, error) {
	return time.Time(jt), nil
}

// Supaya GORM bisa baca dari MySQL
func (jt *JSONTime) Scan(value interface{}) error {
	if val, ok := value.(time.Time); ok {
		*jt = JSONTime(val)
		return nil
	}
	return fmt.Errorf("cannot scan %v into JSONTime", value)
}
