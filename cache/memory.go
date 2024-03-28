package cache

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// Memory struct contains *memcache.Client
type Memory struct {
	sync.Map
}

// NewMemory create new memcache
func NewMemory() *Memory {
	return &Memory{}
}

// Get return cached value
func (m *Memory) Get(key string) interface{} {
	value, loaded := m.LoadAndDelete(key)
	if !loaded {
		return ""
	}
	values, ok := value.([]string)
	if !ok {
		return ""
	}
	expiredAt, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		return ""
	}
	if time.Now().Unix() > expiredAt {
		return ""
	}
	return values[0]
}

// IsExist check value exists in mcache.
func (m *Memory) IsExist(key string) bool {
	_, loaded := m.LoadAndDelete(key)
	return loaded
}

// Set cached value with key and expire time.
func (m *Memory) Set(key string, val interface{}, timeout time.Duration) (err error) {
	m.Store(key, []string{val.(string), strconv.FormatInt(time.Now().Add(timeout).Unix(), 10)})
	return nil
}

// Delete delete value in mcache.
func (m *Memory) Delete(key string) error {
	_, loaded := m.LoadAndDelete(key)
	if !loaded {
		return errors.New("delete error")
	}
	return nil
}
