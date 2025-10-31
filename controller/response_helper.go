package controller

import (
	"encoding/json"
	"net/http"

	"github.com/LanangDepok/ebook-store/model"
)

func respondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := model.Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := model.Response{
		Status:  "error",
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
