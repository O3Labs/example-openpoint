package local

import (
	"encoding/json"
	"errors"
	// "log"
	"sync"
)

type Local struct {
	client map[string][]byte
}

var local *Local
var once sync.Once

func NewLocalClient() *Local {
	once.Do(func() {
		local = &Local{client: map[string][]byte{}}
	})
	return local
}

func (r *Local) Get(key string) interface{} {
	return r.client[key]
}

func (r *Local) GetScan(key string, v interface{}) error {
	result := r.client[key]
	if result == nil {
		return errors.New("nil")
	}
	json.Unmarshal(result, &v)
	return nil
}

func (r *Local) Set(key string, value interface{}) {
	b, _ := json.Marshal(value)
	r.client[key] = b
}
