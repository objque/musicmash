package types

import (
	"bytes"
	"strconv"
	"time"
)

type Time struct {
	Value time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Value.Format(strconv.Quote("2006-01-02"))), nil
}

var emptyDate = []byte{34, 48, 48, 48, 48, 45, 48, 48, 45, 48, 48, 34} // 0000-00-00

func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	t.Value, err = time.ParseInLocation(strconv.Quote("2006-01-02"), string(data), time.FixedZone("UTC", 0))
	if err != nil {
		// some deezer releases have 0000-00-00 release date
		// e.g
		// 'Paddy and the Rats' (id 470404)
		// 'Michael Grimm' (id 567109)
		if bytes.Equal(data, emptyDate) {
			t.Value = time.Date(0001, time.January, 01, 00, 00, 00, 00, time.UTC)
			return nil
		}

		// some oldest releases have year and month
		// for example check albums for user: 14548963, 14589060
		t.Value, err = time.ParseInLocation(strconv.Quote("2006-01"), string(data), time.FixedZone("UTC", 0))
		if err == nil {
			return nil
		}

		// some oldest releases have only year
		// for example check albums for user: 331647880, 14064764, 101672291
		t.Value, err = time.ParseInLocation(strconv.Quote("2006"), string(data), time.FixedZone("UTC", 0))
	}
	return err
}
