package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// 2018-06-18 07:22:59 +0200 CEST
const layout = "2006-01-02 15:04:05 -0700 MST"

type Record struct {
	UUID      string
	Num       float64 //Should we only have 2 digits???
	Timestamp time.Time
}

func (r *Record) Parse(data []string) (err error) {
	if len(data) != 3 {
		return fmt.Errorf("input data should have 3 fields and not %d", len(data))
	}

	if !isValidUUID(data[0]) {
		return fmt.Errorf("value %s is not a valid uuid", data[0])
	}
	r.UUID = data[0]

	if r.Num, err = strconv.ParseFloat(data[1], 64); err != nil {
		err = fmt.Errorf("value %s is not a float: %s", data[1], err)
		return
	}

	if r.Timestamp, err = time.Parse(layout, data[2]); err != nil {
		err = fmt.Errorf("value %s is in the expected time format: %s", data[2], err)
		return
	}

	return
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
