package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

type IntArray []int

// Value 实现 driver.Valuer 接口，用于写入数据库
func (a *IntArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner 接口，用于从数据库中读取
func (a *IntArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("scan source is not []byte")
	}
	return json.Unmarshal(bytes, a)
}
