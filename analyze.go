package main

import (
	"container/ring"
	"fmt"
)

const tries = 10
const failureThreshold = 3
const recoveryThreshold = tries

// Analyze takes a channel of response codes, and, every time it receives a new
// response code, analyzes it's history, detects state transitions, e.g. running -> down
// or down -> running, and emits the corresponding messages on the appropriate channel
func Analyze(url string, responseCodes <-chan int, messages chan<- *Message) {

	isDown := false
	lastTenCodes := ring.New(tries)

	for code := range responseCodes {
		lastTenCodes.Value = code
		lastTenCodes = lastTenCodes.Next()
		failure, recovery, codes := computeState(lastTenCodes)
		fmt.Println(url, failure, codes)
		if failure && !isDown {
			isDown = true
			messages <- NewMessage(url, false, codes)
		} else if isDown && recovery {
			// TODO : make recovery message
			isDown = false
			messages <- NewMessage(url, true, codes)
		}
	}
}

func computeState(codesToAnalyze *ring.Ring) (
	isFailure bool,
	canRecover bool,
	codes []int) {
	failures, successes, codes := 0, 0, make([]int, 0, tries)
	codesToAnalyze.Do(func(c interface{}) {
		if i, ok := c.(int); ok {
			codes = append(codes, i)
			if i < 200 || i >= 300 {
				failures++
			} else {
				successes++
			}
		}
	})
	return failures >= failureThreshold, successes >= recoveryThreshold, codes
}
