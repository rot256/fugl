package fugl

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestCanarySerializeCycle(t *testing.T) {
	canary := Canary{
		Version: 1,
		Author: "John Doe",
		Creation: CanaryTime(time.Now()),
		Expiry: CanaryTime(time.Now()),
		Promises: []string{"example"},
		Nonce: "nonce",
		Final: false,
	}

	bs, err := json.Marshal(canary)
	if err != nil {
		t.Fatalf("error encoding Canary to json, err=%v", err)
	}

	var out Canary
	err = json.Unmarshal(bs, &out)
	if err != nil {
		t.Fatalf("error decoding Canary from json, err=%v", err)
	}

	if !canary.Equal(out) {
		fmt.Printf("canary=%v\n", canary)
		fmt.Printf("out   =%v\n", out)

		t.Fatal("serialization cycle mismatch")
	}
}
