package datatypes

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Timestamptz implements custom type for storing and reading time with timezone in the database.
type Timestamptz time.Time

func (t *Timestamptz) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.unmarshal(string(v))
	case string:
		return t.unmarshal(v)
	default:
		return fmt.Errorf("unexpected Timestamptz type: %#v", v)
	}
}

func (t Timestamptz) Value() (driver.Value, error) {
	return driver.Value(time.Time(t).Format(time.RFC3339)), nil
}

func (Timestamptz) GormDataType() string {
	return "text"
}

// MarshalJSON writes a quoted string in the custom format.
func (t Timestamptz) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

// String returns the time in the custom format.
func (t Timestamptz) String() string {
	return time.Time(t).String()
}

func (t *Timestamptz) unmarshal(value string) error {
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}
	*t = Timestamptz(parsedTime)
	return nil
}
