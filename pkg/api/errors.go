package api

import (
	"encoding/json"
	"errors"
	"io"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ExtractError(body io.Reader) error {
	response := ErrorResponse{}
	_ = json.NewDecoder(body).Decode(&response)
	return errors.New(response.Error)
}
