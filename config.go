package main

import "time"

// Resource represents an Url to be watched, and all the necessary parameters
type Config struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Url               string        `json:"url"`               // the Url to watch
	IntervalInSeconds time.Duration `json:"intervalInSeconds"` // the interval at which to poll the resource. Defaults to 60s.
	TimeoutInSeconds  time.Duration `json:"timeoutInseconds"`  // the timeout in seconds. Defaults to 2s.
	NumberOfTries     int           `json:"numberOfTries"`     // number of attempts at getting the resource. Defaults to 10.
	FailureThreshold  int           `json:"failureThreshold"`  // the number of failures within the tries that would raise an alarm. Defaults to 5.
}

// LoadConfig loads all resources from a JSON file
func LoadConfig(path string) Config {
	return Config{}
}
