package database

import (
	"fmt"
	"sync"
)

type Payload struct {
	Event     Event   `json:"event"`
	CreatedAt string  `json:"created_at"`
	ID        string  `json:"id"`
	Trigger   Trigger `json:"trigger"`
	Table     Table   `json:"table"`
}

type Event struct {
	SessionVariables map[string]string `json:"session_variables"`
	Op               string            `json:"op"`
	Data             Data              `json:"data"`
}

type Data struct {
	Old interface{} `json:"old"`
	New interface{} `json:"new"`
}

type Trigger struct {
	Name string `json:"name"`
}

type Table struct {
	Schema string `json:"schema"`
	Name   string `json:"name"`
}

type InMemoryDB struct {
	mu   sync.RWMutex
	data map[string]Payload
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string]Payload),
	}
}

func (db *InMemoryDB) Set(key string, value Payload) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[key] = value
	return nil
}

func (db *InMemoryDB) Get(key string) (Payload, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	value, ok := db.data[key]
	if !ok {
		return Payload{}, fmt.Errorf("no value found for key: %v", key)
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
