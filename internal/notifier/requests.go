package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/musicmash/musicmash/internal/version"
)

var userAgent = fmt.Sprintf("musicmash-server/%v", version.Commit)

func generateRequestID() string {
	return uuid.New().String()
}

func (n *Notifier) sendReleases(releases []*Notification) error {
	b, err := json.Marshal(&releases)
	if err != nil {
		return fmt.Errorf("can't marshal releases: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, n.uri.String(), bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("can't make notification request: %w", err)
	}

	request.Header.Set("User-Agent", userAgent)
	// request-id may be useful on the external notification service
	// for e.g it may check few parallel requests on duplicated req-id
	request.Header.Set("X-Request-Id", generateRequestID())

	resp, err := n.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("tried to send request with notifications, but got err: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		b, _ = ioutil.ReadAll(resp.Body)
		log.Warnf("error body from the notifications server: %v", string(b))
		return fmt.Errorf("notifications server return status code %v, but expect 200", resp.StatusCode)
	}

	return nil
}
