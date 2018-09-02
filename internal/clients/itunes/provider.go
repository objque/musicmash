package itunes

type Provider struct {
	Token  string
	URL    string
	Region string
}

func NewProvider(url, token string) *Provider {
	return &Provider{URL: url, Token: token, Region: "us"}
}
