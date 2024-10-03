package main

import (
	"encoding/json"
	"net/http"
)

func (h handler) plug(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := BuildAggregateResponse{Success: true, Data: AggrLenRecords{LenRecords: 1}}
	json.NewEncoder(w).Encode(response)
}
