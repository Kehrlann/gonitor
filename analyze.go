package main

import (
	"container/ring"
)

// Analyze takes a channel of response Codes, and, every time it receives a new
// response code, analyzes it's history, detects state transitions, e.g. running -> down
// or down -> running, and emits the corresponding messages on the appropriate channel
func Analyze(resource Resource, responseCodes <-chan int, messages chan<- *StateChangeMessage) {

	isDown := false
	lastTenCodes := ring.New(resource.NumberOfTries)
	firstAlertThreshold := resource.NumberOfTries - 1
	for code := range responseCodes {
		lastTenCodes.Value = code
		lastTenCodes = lastTenCodes.Next()

		codes := RingToIntSlice(lastTenCodes) // slightly overkill, but nice for testing and printing

		failure, recovery := computeState(codes, resource.FailureThreshold, resource.NumberOfTries)
		if firstAlertThreshold > 0 {
			// PASS : we do not send a message until we have enough data
			firstAlertThreshold--
		} else if failure && !isDown {
			isDown = true
			messages <- ErrorMessage(resource, codes)
		} else if isDown && recovery {
			isDown = false
			messages <- RecoveryMessage(resource, codes)
		}
	}
}

// computeState takes HTTP Codes, and tells you whether this could trigger an
// alarm, and whether it could trigger or a recovery.
func computeState(codesToAnalyze []int, failureThreshold int, recoveryThreshold int) (isFailure bool, canRecover bool) {
	failures, successes := 0, 0
	for _, code := range codesToAnalyze {
		if code < 200 || code >= 300 {
			failures++
		} else {
			successes++
		}
	}
	return failures >= failureThreshold, successes >= recoveryThreshold
}
