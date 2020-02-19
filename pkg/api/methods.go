package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"moul.io/http2curl"
)

func do(provider *Provider, url fmt.Stringer, method string, header http.Header, body io.Reader, out interface{}) error {
	request, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return err
	}

	if header != nil {
		request.Header = header
	}

	command, _ := http2curl.GetCurlCommand(request)
	log.Debug(command)

	resp, err := provider.Client.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		// TODO (m.kalinin): add resp.body limit length
		log.Debugln(fmt.Sprintf("Got %v status_code with resp: %s", resp.StatusCode, string(b)))
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return ExtractError(bytes.NewReader(b))
	}

	// TODO (m.kalinin): decode subscriptions after create
	if out == nil {
		return nil
	}

	return json.NewDecoder(bytes.NewReader(b)).Decode(out)
}

func Get(provider *Provider, url fmt.Stringer, out interface{}) error {
	return do(provider, url, http.MethodGet, nil, nil, out)
}

func GetWithHeaders(provider *Provider, url fmt.Stringer, header http.Header, out interface{}) error {
	return do(provider, url, http.MethodGet, header, nil, out)
}

func Post(provider *Provider, url fmt.Stringer, body io.Reader, out interface{}) error {
	return do(provider, url, http.MethodPost, nil, body, out)
}

func PostWithHeaders(provider *Provider, url fmt.Stringer, header http.Header, body io.Reader, out interface{}) error {
	return do(provider, url, http.MethodPost, header, body, out)
}

func PatchWithHeaders(provider *Provider, url fmt.Stringer, header http.Header, body io.Reader, out interface{}) error {
	return do(provider, url, http.MethodPatch, header, body, out)
}

func DeleteWithHeaders(provider *Provider, url fmt.Stringer, header http.Header, body io.Reader) error {
	return do(provider, url, http.MethodDelete, header, body, nil)
}
