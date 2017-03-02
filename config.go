package main

import "time"

// Config represents an Url to be watched, and all the necessary parameters
type Config struct {
	url               string        // the Url to watch
	intervalInSeconds time.Duration // the interval at which to poll the resource. Defaults to 60s.
	timeoutInSeconds  time.Duration // the timeout in seconds. Defaults to 2s.
	numberOfTries     int           // number of attempts at getting the resource. Defaults to 10.
	failureThreshold  int           // the number of failures within the tries that would raise an alarm. Defaults to 5.
}
