package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringArray []string

// Value 实现 driver.Valuer 接口
func (a *StringArray) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Failed to unmarshal JSON value")
	}
	return json.Unmarshal(bytes, &a)
}
