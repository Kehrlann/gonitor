package main

import (
	"net/http"
	"time"
)

var timeout = 1 * time.Second

// Run takes a resource and polls the given HTTP url for errors , and emits failure / recovery messages accordingly
func Run(resource Resource, every time.Duration, messages chan<- *StateChangeMessage) {
	responseCodes := make(chan int)
	go Analyze(resource, responseCodes, messages)
	for range time.Tick(every) {
		responseCodes <- fetch(resource.Url)
	}
}

func fetch(url string) int {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0
	}

	return resp.StatusCode
}
