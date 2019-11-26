package httputils

import (
	"fmt"
	"net/http"
	"time"
)

const dateLayout = "2006-01-02T15:04:05"

func GetQueryTime(r *http.Request, key string) (*time.Time, error) {
	since := r.URL.Query().Get(key)
	if len(since) == 0 {
		return nil, fmt.Errorf("%v argument didn't provided", key)
	}

	date, err := time.Parse(dateLayout, since)
	if err != nil {
		return nil, fmt.Errorf("%v argument %s does not match %s layout", key, since, dateLayout)
	}
	return &date, nil
}
