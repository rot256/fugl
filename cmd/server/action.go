package main

import (
	"os/exec"
	"strings"
	"time"
)

/* Runs an optional command unpon failure to submit new canary
 *
 * This can be used in a dead man's switch type scenario
 */

func actionRunner(command string, state *ServerState) {
	// check if feature enabled
	if command == "" {
		return
	}
	parts := strings.Fields(command)
	logInfo("Failure action:", parts)

	// wait for deadline
	for ticker := time.NewTicker(time.Second * 5); ; <-ticker.C {
		// check if canary on server
		state.canaryLock.RLock()
		if state.latestCanary == nil {
			state.canaryLock.RUnlock()
			continue
		}

		// check time
		now := time.Now()
		deadline := state.latestCanary.Deadline.Time()
		state.canaryLock.RUnlock()
		if now.After(deadline) {
			logInfo("Running failure action")
			out, err := exec.Command(parts[0], parts[1:]...).Output()
			if err != nil {
				logWarning("Failed to execute failure action:", err)
			}
			logInfo("Failure action output:\n" + string(out))
			select {} // wait forever
		}

		// sleep
		logDebug("Action Runner going to sleep:", deadline.Sub(now).Seconds())
		time.Sleep(deadline.Sub(now))
		logDebug("Woke from sleep")
	}
}
