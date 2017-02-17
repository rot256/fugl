package main

import (
	"github.com/rot256/fugl"
	"golang.org/x/crypto/openpgp"
	"net/http"
	"sync"
	"time"
)

type ServerState struct {
	storeDir       string          // directory for storing new canaries
	latestCanary   *fugl.Canary    // cached latest canary (parsed proof)
	latestProof    string          // newest proof
	canaryKey      *openpgp.Entity // parsed public key
	canaryKeyArmor string          // ascii armored pgp key
	canaryLock     sync.RWMutex
}

func SendRequestError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
	return
}

/* Serves the public key */

type GetKeyHandler struct {
	state *ServerState
}

func (h *GetKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.state.canaryKeyArmor == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(h.state.canaryKeyArmor))
}

/* Serves the latest published canary */

type LatestHandler struct {
	state *ServerState
}

func (h *LatestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.state.canaryLock.RLock()
	defer h.state.canaryLock.RUnlock()
	if h.state.latestProof == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(h.state.latestProof))
}

/* Add a new canary */

type SubmitHandler struct {
	state *ServerState
}

func (h *SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse and verify signature
	proof := r.PostFormValue(fugl.SERVER_SUBMIT_FIELD_NAME)
	if proof == "" {
		SendRequestError(w, "No proof form attribute specified")
		return
	}
	logDebug("New proof submission:\n", proof)

	canary, _, err := fugl.OpenProof(h.state.canaryKey, proof)
	if err != nil {
		SendRequestError(w, err.Error())
		return
	}
	if canary == nil {
		SendRequestError(w, "Unable to load canary from proof")
		return
	}

	// check version field
	if canary.Version != fugl.CanaryVersion {
		logWarning("Invalid canary version field")
		SendRequestError(w, "Unsupported canary version")
		return
	}

	// verify expires in the future
	if time.Now().After(canary.Expiry.Time()) {
		SendRequestError(w, "Canary must have a deadline in the future")
		return
	}

	// take write lock
	h.state.canaryLock.Lock()
	defer h.state.canaryLock.Unlock()

	// verify canary fields
	err = fugl.CheckCanary(canary, h.state.latestCanary, time.Now())
	if err != nil {
		SendRequestError(w, err.Error())
		return
	}

	// save to disk
	err = fugl.SaveToDirectory(proof, h.state.storeDir, canary.Expiry.Time())
	if err != nil {
		logError("Failed to save valid proof to store:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.state.latestProof = proof
	h.state.latestCanary = canary
	logInfo("Succesfully added a new canary")
	w.WriteHeader(http.StatusNoContent)
}
