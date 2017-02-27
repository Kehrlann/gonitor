package main

import (
	"container/ring"
)

const tries = 10
const failureThreshold = 3
const recoveryThreshold = tries

// Analyze takes a channel of response codes, and, every time it receives a new
// response code, analyzes it's history, detects state transitions, e.g. running -> down
// or down -> running, and emits the corresponding messages on the appropriate channel
func Analyze(url string, responseCodes <-chan int, messages chan<- *StateChangeMessage) {

	isDown := false
	lastTenCodes := ring.New(tries)

	for code := range responseCodes {
		lastTenCodes.Value = code
		lastTenCodes = lastTenCodes.Next()

		codes := RingToIntSlice(lastTenCodes) // slightly overkill, but nice for testing and printing

		failure, recovery := computeState(codes)
		if failure && !isDown {
			isDown = true
			messages <- ErrorMessage(url, codes)
		} else if isDown && recovery {
			// TODO : make recovery message
			isDown = false
			messages <- RecoveryMessage(url, codes)
		}
	}
}

// computeState takes HTTP codes, and tells you whether this could trigger an
// alarm, and whether it could trigger or a recovery.
func computeState(codesToAnalyze []int) (isFailure bool, canRecover bool) {
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
