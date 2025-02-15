package core

import (
	"encoding/json"

	"github.com/Streamlet/gosql"
)

func ConfigSet(c *gosql.Connection, name string, value interface{}) *Error {
	var stringVal string
	if _, ok := value.(string); !ok {
		jsonVal, err := json.Marshal(value)
		if err != nil {
			return NewError(Error_JsonEncodeError, err.Error())
		}
		stringVal = string(jsonVal)
	} else {
		stringVal = value.(string)
	}
	_, err := c.Update("INSERT INTO config (name, value) VALUES (?, ?) ON DUPLICATE KEY UPDATE value = ?",
		name, stringVal, stringVal)
	if err != nil {
		return NewError(Error_DbError, err.Error())
	}
	return nil
}

func ConfigGet[T any](c *gosql.Connection, name string) (*T, *Error) {
	type ConfigItem struct {
		Value *string `db:"value"`
	}

	items, err := gosql.Select[ConfigItem](c, "SELECT value FROM config WHERE name = ? LIMIT 1", name)
	if err != nil {
		return nil, NewError(Error_DbError, err.Error())
	}
	if len(items) != 1 || items[0].Value == nil {
		return nil, nil
	}
	var t T
	if err := json.Unmarshal([]byte(*items[0].Value), &t); err != nil {
		return nil, NewError(Error_JsonDecodeError, err.Error())
	}

	return &t, nil
}
