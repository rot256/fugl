package fugl

import (
	"errors"
	"fmt"
	"reflect"
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
	News     []string   `json:"news"`     // list of BBC news articles
}

func (c Canary) Equal(other Canary) bool {
	return (c.Version == other.Version) &&
		(c.Author == other.Author) &&
		c.Creation.Time().Equal(other.Creation.Time()) &&
		c.Expiry.Time().Equal(other.Expiry.Time()) &&
		reflect.DeepEqual(c.Promises, other.Promises) &&
		(c.Nonce == other.Nonce) &&
		(c.Final == other.Final)
}

/* specifies the time format used in the canaries
 */

func (t CanaryTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.String())
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
	return time.Time(t).Round(1 * time.Second)
}

func (t CanaryTime) String() string {
	return t.Time().Format(CanaryTimeFormat)
}
