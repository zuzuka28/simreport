package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"document/internal/model"
)

func responseErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	type doc struct {
		Message string `json:"message"`
	}

	body := doc{Message: err.Error()}
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, model.ErrInvalid):
		status = http.StatusBadRequest

	case errors.Is(err, model.ErrNotFound):
		status = http.StatusNotFound

	case errors.Is(err, context.Canceled):
		status = 499
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(body) //nolint:errchkjson
}
