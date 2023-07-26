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

	db := database.NewInMemoryDB()

	r.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var payload Payload
		err := decoder.Decode(&payload)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.Set(payload.ID, payload)
		if err != nil {
			log.Printf("Error saving to in-memory DB: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	r.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		keys, err := db.Keys()
		if err != nil {
			log.Printf("Error fetching from in-memory DB: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keys)
	}).Methods("GET")

	r.HandleFunc("/events/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		payload, err := db.Get(id)
		if err != nil {
			log.Printf("Error fetching from in-memory DB: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}).Methods("GET")

	fmt.Println("Starting server on port: 8080")
	http.ListenAndServe(":8080", r)
}
