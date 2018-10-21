package itunes

import "go.uber.org/ratelimit"

const rateLimit = 30 // requests per second

type Provider struct {
	Token           string
	URL             string
	Region          string
	requestsLimiter ratelimit.Limiter
}

func (p *Provider) WaitRateLimit() {
	p.requestsLimiter.Take()
}

func NewProvider(url, token string) *Provider {
	rl := ratelimit.New(rateLimit)
	return &Provider{URL: url, Token: token, Region: "us", requestsLimiter: rl}
}
