package api

import (
	"fmt"
	"net/http"
)

type Provider struct {
	URL    string
	Client *http.Client
}

func NewProvider(url string, apiVersion int) *Provider {
	return &Provider{
		URL:    fmt.Sprintf("%v/v%d", url, apiVersion),
		Client: &http.Client{},
	}
}
