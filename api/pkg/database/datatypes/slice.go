package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Scan is implementing sql.Scanner interface.
func Scan(dest interface{}, value interface{}) error {
	switch t := value.(type) {
	case string:
		return json.Unmarshal([]byte(t), dest)
	case []byte:
		return json.Unmarshal(t, dest)
	default:
		// Undefined types
		return errors.New(fmt.Sprint("failed to unmarshal JSON value:", value))
	}
}

// Value is implementing driver.Valuer interface.
func Value(src interface{}) (driver.Value, error) {
	value, err := json.Marshal(src)
	return string(json.RawMessage(value)), err
}

// Slice is a generic variation for presenting any driver.Value supported values.
type Slice[T any] []T

func (j *Slice[any]) Scan(value interface{}) error {
	return Scan(j, value)
}

func (j Slice[any]) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return Value(j)
}
