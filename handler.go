package main

import (
	"database/sql"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"net/http"
	"time"
)

type BatchHandler struct {
	DB                 *sql.DB
	batchSaveHistogram prometheus.Histogram
}

func (h *BatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.create
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *BatchHandler) create(w http.ResponseWriter, r *http.Request) {
	var batchEvents BatchEvents
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&batchEvents); err != nil {
		log.Warn(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	start := time.Now()
	err := batchEvents.createBatchEvents(h.DB)
	elapsed := float64(time.Since(start).Nanoseconds())
	h.batchSaveHistogram.Observe(elapsed)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, batchEvents)
}
