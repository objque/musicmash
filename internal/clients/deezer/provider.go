package deezer

import "go.uber.org/ratelimit"

const rateLimit = 9 // requests per second

type Provider struct {
	URL             string
	requestsLimiter ratelimit.Limiter
}

func (p *Provider) WaitRateLimit() {
	p.requestsLimiter.Take()
}

func NewProvider(url string) *Provider {
	rl := ratelimit.New(rateLimit)
	return &Provider{URL: url, requestsLimiter: rl}
}
