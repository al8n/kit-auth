package httpcodec

import (
	"context"
	"encoding/json"
	"net/http"
)

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorWrapper{Error: err.Error()})
}

type ErrorWrapper struct {
	Error string `json:"error"`
}
