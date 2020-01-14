package httputils

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteInternalError(w http.ResponseWriter) {
	WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
}

func WriteError(w http.ResponseWriter, err error) {
	WriteErrorWithCode(w, http.StatusBadRequest, err)
}

func WriteErrorWithCode(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	b, _ := json.Marshal(&ErrorResponse{Error: err.Error()})
	_, _ = w.Write(b)
}
