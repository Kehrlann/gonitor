package main

import (
	"net/http"
	"time"
)

var timeout = 1 * time.Second

// Run takes a resource and polls the given HTTP url for errors , and emits failure / recovery messages accordingly
func Run(resource Resource, messages chan<- *StateChangeMessage) {
	responseCodes := make(chan int)
	go Analyze(resource, responseCodes, messages)
	for range time.Tick(time.Duration(resource.IntervalInSeconds) * time.Second) {
		responseCodes <- fetch(resource.Url, resource.TimeoutInSeconds)
	}
}

func fetch(url string, timeOutInSeconds int) int {
	client := &http.Client{
		Timeout: time.Duration(timeOutInSeconds) * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0
	}

	return resp.StatusCode
}
