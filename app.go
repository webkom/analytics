package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type App struct {
	Router             *mux.Router
	DB                 *sql.DB
	batchSaveHistogram prometheus.Histogram
}

func (a *App) Initialize(postgresUrl string) {
	var err error
	a.DB, err = sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.InitializePrometheus()
	a.InitializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) InitializeRoutes() {
	batchHandler := &BatchHandler{a.DB, a.batchSaveHistogram}
	a.Router.HandleFunc("/v1/batch", batchHandler.ServeHTTP)

	a.Router.Handle("/metrics", promhttp.Handler())
}

func (a *App) InitializePrometheus() {
	var (
		batchSaveHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "batch_save_nanoseconds_histogram",
			Help: "Track the patch save operation",
		})
	)
	prometheus.MustRegister(batchSaveHistogram)
	a.batchSaveHistogram = batchSaveHistogram
}

func (a *App) Migrate() {
	const tableCreationQuery = `
	CREATE TABLE IF NOT EXISTS identifies
	(
	id			SERIAL NOT NULL PRIMARY KEY,
   	anonymous_id    	TEXT,
   	user_id			TEXT,
   	context			JSONB,
   	traits			JSONB,
   	received_at		TIMESTAMP NOT NULL,
	sent_at			TIMESTAMP NOT NULL,
	timestamp		TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS pages
	(
	id			SERIAL NOT NULL PRIMARY KEY,
   	anonymous_id    	TEXT,
   	user_id			TEXT,
   	context			JSONB,
   	name			TEXT,
   	properties		JSONB,
   	received_at		TIMESTAMP NOT NULL,
	sent_at			TIMESTAMP NOT NULL,
	timestamp		TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS screens
	(
	id			SERIAL NOT NULL PRIMARY KEY,
   	anonymous_id    	TEXT,
   	user_id			TEXT,
   	context			JSONB,
   	name			TEXT,
   	category		TEXT,
   	properties		JSONB,
   	received_at		TIMESTAMP NOT NULL,
	sent_at			TIMESTAMP NOT NULL,
	timestamp		TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tracks
	(
	id			SERIAL NOT NULL PRIMARY KEY,
   	anonymous_id    	TEXT,
   	user_id			TEXT,
   	context			JSONB,
   	properties		JSONB,
   	event 			TEXT,
   	received_at		TIMESTAMP NOT NULL,
	sent_at			TIMESTAMP NOT NULL,
	timestamp		TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS "groups"
	(
	id			SERIAL NOT NULL PRIMARY KEY,
   	anonymous_id    	TEXT,
   	user_id			TEXT,
   	context			JSONB,
   	traits			JSONB,
   	group_id		TEXT,
   	received_at		TIMESTAMP NOT NULL,
	sent_at			TIMESTAMP NOT NULL,
	timestamp		TIMESTAMP NOT NULL
	);
	`
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
