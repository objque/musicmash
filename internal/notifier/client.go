package notifier

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func newHTTPClient() *http.Client {
	retryClient := retryablehttp.NewClient()

	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 120 * time.Second
	retryClient.RetryMax = 5

	return retryClient.StandardClient()
}
