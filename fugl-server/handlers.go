package main

import (
	"encoding/json"
	"net/http"
	"time"
)

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

/* Returns canary status on this node
 */
func HandleStatus(w http.ResponseWriter, r *http.Request) {
	// Generate info struct
	var status CanaryStatus
	status.Version = CanaryVersion
	status.Key = CanaryKeyArmor

	// Serialize
	resp, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		logFatal("JSON error:", err)
		return
	}
	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusFound)
	w.Write(resp)
}

/* Serves the latest published canary
 */
func HandleLatest(w http.ResponseWriter, r *http.Request) {
	if LatestProof == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(LatestProof))
}

func SendError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
	return
}

/* Add a new canary to the database
 */
func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	// parse and verify signature
	proof := r.PostFormValue("proof")
	logInfo("New proof submission:", proof)
	canary, err := VerifyProof(CanaryKey, proof)
	if err != nil {
		SendError(w, err.Error())
	}

	// check version field
	if canary.Version != CanaryVersion {
		logWarning("Invalid canary version field")
		SendError(w, "Unsupported canary version")
		return
	}

	// verify previous canary hash
	if LatestProof != "" {
		hash := Sha256StringToHex(proof)
		if hash != canary.Previous {
			SendError(w, "Canary must reference preceeding canary hash")
			return
		}
	}

	// verify deadline in future (avoid bricking)
	if time.Now().After(time.Time(canary.Deadline)) {
		SendError(w, "Canary must have a deadline in the future")
		return
	}

	// save to disk
	err = SaveToStore(proof, time.Now())
	if err != nil {
		logFatal("Failed to save valid proof to store:", err)
	}
	LatestProof = proof
	LatestCanary = canary
	logInfo("Succesfully added a new canary")
	w.WriteHeader(http.StatusNoContent)
}
