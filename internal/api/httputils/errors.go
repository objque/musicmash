package httputils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, err error) {
	WriteErrorWithCode(w, http.StatusBadRequest, err)
}

func WriteErrorWithCode(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	b, _ := json.Marshal(&ErrorResponse{Error: err.Error()})
	_, _ = w.Write(b)
}
