package api

import (
	"errors"
	"net/http"
)

func GetUser(r *http.Request) (string, error) {
	userName := r.Header.Get("x-user-name")
	if userName == "" {
		return "", errors.New("x-user-name header didn't provided")
	}
	return userName, nil
}
