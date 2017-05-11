package main

import (
	"container/ring"

	"github.com/kehrlann/gonitor/config"
)

// Analyze takes a channel of response Codes, and, every time it receives a new
// response code, analyzes it's history, detects state transitions, e.g. running -> down
// or down -> running, and emits the corresponding messages on the appropriate channel
func Analyze(resource config.Resource, responseCodes <-chan int, messages chan<- *StateChangeMessage) {

	isDown := false
	lastHttpReturnCodes := ring.New(resource.NumberOfTries)
	numberOfTriesBeforeFirstAlert := resource.NumberOfTries - 1
	for code := range responseCodes {
		lastHttpReturnCodes.Value = code
		lastHttpReturnCodes = lastHttpReturnCodes.Next()

		codes := RingToIntSlice(lastHttpReturnCodes) // slightly overkill, but nice for testing and printing

		failure, recovery := computeState(codes, resource.FailureThreshold, resource.NumberOfTries)
		if numberOfTriesBeforeFirstAlert > 0 {
			// PASS : we do not send a message until we have enough data
			numberOfTriesBeforeFirstAlert--
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
func computeState(codesToAnalyze []int, failureThreshold int, recoveryThreshold int) (isFailure bool, isRecovery bool) {
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
