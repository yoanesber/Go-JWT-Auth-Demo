package customtype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

type Date struct {
	time.Time
}

// To unmarshal a JSON date string in the format "YYYY-MM-DD" into a Date struct.
// It handles empty strings and null values by returning a zero Date value.
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	d.Time = t
	return nil
}

// MarshalJSON formats the Date struct into a JSON string in the format "YYYY-MM-DD".
// It returns an error if the date is not in the expected format.
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(dateFormat))
}

// Value implements the driver.Valuer interface for the Date type.
// It converts the Date to a time.Time value for database storage.
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// Scan implements the sql.Scanner interface for the Date type.
// It converts a time.Time value from the database into a Date struct.
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Date", value)
	}
}

// String formats the Date struct into a string in the format "YYYY-MM-DD".
// It is useful for displaying the date in a human-readable format.
func (d Date) String() string {
	return d.Time.Format(dateFormat)
}
