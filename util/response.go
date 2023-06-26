package util

import (
	"encoding/json"
	"net/http"
)

func Success(w http.ResponseWriter, code int, data interface{}, message string) {
	response, err := json.Marshal(map[string]interface{}{
		"data":    data,
		"success": true,
		"message": message,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}

func Error(w http.ResponseWriter, code int, data interface{}, message string) {
	response, err := json.Marshal(map[string]interface{}{
		"data":    data,
		"message": message,
		"success": false,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}
