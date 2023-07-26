package database

import (
	"fmt"
	"sync"
)

type InMemoryDB struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		mu:   sync.RWMutex{},
		data: make(map[string]interface{}),
	}
}

func (db *InMemoryDB) Set(key string, value interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[key] = value
	return nil
}

func (db *InMemoryDB) Get(key string) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	value, ok := db.data[key]
	if !ok {
		return nil, fmt.Errorf("no value found for key: %v", key)
	}
	return value, nil
}

func (db *InMemoryDB) Keys() ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	keys := make([]string, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}
	return keys, nil
}
