package itunes

import "time"

const layout1 = "Jan 2, 2006"
const layout2 = "2 Jan 2006"

func parseTime(value string) (*time.Time, error) {
	t, err := time.Parse(layout1, value)
	if err != nil {
		t, err = time.Parse(layout2, value)
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}
