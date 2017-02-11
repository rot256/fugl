package fugl

import (
	"errors"
	"fmt"
	"time"
)

type CanaryTime time.Time

type Canary struct {
	Version  int64      `json:"version"`  // Canary struct version
	Author   string     `json:"author"`   // Publishing entity of the canary
	Creation CanaryTime `json:"creation"` // Time of creation
	Expiry   CanaryTime `json:"expiry"`   // Expiry time of canary
	Promises []string   `json:"promises"` // Set of promises (may be empty)
	Nonce    string     `json:"nonce"`    // Random nonce
	Final    bool       `json:"final"`    // Is this canary final?
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
