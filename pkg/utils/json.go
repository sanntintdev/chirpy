package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithErr(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}

	if code > 499 {
		log.Printf("Server error: %v", err)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errorResponse{Error: msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request, dest *T) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	
	if err := decoder.Decode(dest); err != nil {
		RespondWithErr(w, http.StatusBadRequest, "Invalid JSON", err)
		return err
	}

	return nil
}
