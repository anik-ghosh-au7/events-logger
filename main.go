package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/anik-ghosh-au7/events-logger/database"
	"github.com/gorilla/mux"
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

func main() {
	r := mux.NewRouter()

	client := database.NewClient()

	r.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var payload Payload
		err := decoder.Decode(&payload)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = client.Set(payload.ID, payload, 0).Err()
		if err != nil {
			log.Printf("Error saving to Redis: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	r.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		var keys []string

		iter := client.Scan(0, "*", 0).Iterator()
		for iter.Next() {
			keys = append(keys, iter.Val())
		}

		if err := iter.Err(); err != nil {
			log.Printf("Error fetching from Redis: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keys)
	}).Methods("GET")

	r.HandleFunc("/events/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		val, err := client.Get(id).Result()
		if err != nil {
			log.Printf("Error fetching from Redis: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
	}).Methods("GET")

	fmt.Println("Starting server on port: 8080")
	http.ListenAndServe(":8080", r)
}
