package itunes

import "time"

const layout  =  "Jan 2, 2006"

func parseTime(value string) (*time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
