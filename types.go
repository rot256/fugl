package fugl

import (
	"errors"
	"fmt"
	"time"
)

type CanaryTime time.Time

type Canary struct {
	Version  int64      `json:"version"`  // Struct version
	Message  string     `json:"message"`  // Optional notification
	Previous string     `json:"previous"` // Hash of previous canary
	Deadline CanaryTime `json:"deadline"` // New deadline
	Nonce    string     `json:"nonce"`    // Random nonce
}

type CanaryStatus struct {
	Version int64  `json:"version"` // Current struct version
	Enabled bool   `json:"enabled"` // Canaries enabled?
	Key     string `json:"key"`     // Current PGP key
}

/* specifies the time format used in the canaries
 */

func (t CanaryTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(CanaryTimeFormat))
	return []byte(stamp), nil
}

func (t *CanaryTime) UnmarshalJSON(val []byte) error {
	str := string(val)
	if len(str) < 2 {
		return errors.New("Time field too short")
	}
	if str[0] != '"' || str[len(str)-1] != '"' {
		return errors.New("Time must be json string type")
	}
	date, err := time.Parse(CanaryTimeFormat, str[1:len(str)-1])
	if err != nil {
		return err
	}
	*t = CanaryTime(date)
	return nil
}

func (t CanaryTime) Time() time.Time {
	return time.Time(t)
}
