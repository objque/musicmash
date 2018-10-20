package deezer

type Provider struct {
	URL string
}

func NewProvider(url string) *Provider {
	return &Provider{URL: url}
}
