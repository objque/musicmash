package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/musicmash/musicmash/internal/version"
)

var userAgent = fmt.Sprintf("musicmash-server/%v", version.Commit)

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

	resp, err := n.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("tried to send request with notifications, but got err: %w", err)
	}

	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notifications server return status code %v, but expect 200", resp.StatusCode)
	}

	return nil
}
