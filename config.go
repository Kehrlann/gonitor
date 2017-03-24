package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config is the application config
type Config struct {
	Resources []Resource
}

// Resource represents an Url to be watched, and all the necessary config/timing parameters
type Resource struct {
	Url               string `json:"url"`               // the Url to watch
	IntervalInSeconds int    `json:"intervalInSeconds"` // the interval at which to poll the resource. Defaults to 60s.
	TimeoutInSeconds  int    `json:"timeoutInSeconds"`  // the timeout in seconds. Defaults to 2s.
	NumberOfTries     int    `json:"numberOfTries"`     // number of attempts at getting the resource. Defaults to 10.
	FailureThreshold  int    `json:"failureThreshold"`  // the number of failures within the tries that would raise an alarm. Defaults to 5.
}

// LoadConfig loads a config from a JSON file
func LoadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	ret := &Config{}
	if err := json.Unmarshal(file, ret); err != nil {
		return nil, err
	}

	return ret, nil
}
