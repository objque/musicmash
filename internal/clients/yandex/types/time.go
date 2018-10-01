package types

import (
	"fmt"
	"strconv"
	"time"
)

type Time struct {
	Value time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Value.Format(strconv.Quote("2006-01-02"))), nil
}

var layouts = []string{"2006-01-02T00:00:00+03:00", "2006-01-02", "2006-01", "2006"}

func (t *Time) UnmarshalJSON(data []byte) error {
	input := string(data)
	var err error
	for _, layout := range layouts {
		t.Value, err = time.Parse(strconv.Quote(layout), input)
		if err != nil {
			continue
		}
		return nil
	}
	return fmt.Errorf("can't parse datetime from '%s'", input)
}
