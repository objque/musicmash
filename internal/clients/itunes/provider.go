package itunes

import (
	"net/http"
	"time"

	"go.uber.org/ratelimit"
)

const rateLimit = 30 // requests per second

type Provider struct {
	Token           string
	URL             string
	Region          string
	HTTPClient      *http.Client
	requestsLimiter ratelimit.Limiter
}

func (p *Provider) WaitRateLimit() {
	p.requestsLimiter.Take()
}

func NewProvider(url, token string, timeout time.Duration) *Provider {
	rl := ratelimit.New(rateLimit)
	client := http.Client{
		Timeout: timeout,
	}
	provider := &Provider{
		URL:             url,
		Token:           token,
		Region:          "us",
		requestsLimiter: rl,
		HTTPClient:      &client,
	}
	return provider
}
