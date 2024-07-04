package utils

import (
	"encoding/json"
	"fmt"
	"sync"
)

func PrintMap(mapp *sync.Map) (string, error) {
	m := map[string]interface{}{}
	mapp.Range(func(key, value interface{}) bool {
		m[fmt.Sprint(key)] = value
		return true
	})

	b, err := json.MarshalIndent(m, "", " ")

	return string(b), err
}
