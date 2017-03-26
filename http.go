package main

import (
	"net/http"
	"time"
)

// Run takes a resource and polls the given HTTP url for errors , and emits failure / recovery messages accordingly
func Run(resource Resource, messages chan<- *StateChangeMessage) {
	responseCodes := make(chan int)
	client := &http.Client{
		Timeout: time.Duration(resource.TimeoutInSeconds) * time.Second,
	}

	go Analyze(resource, responseCodes, messages)

	for range time.Tick(time.Duration(resource.IntervalInSeconds) * time.Second) {
		responseCodes <- fetch(client, resource.Url)
	}
}

func fetch(client *http.Client, url string) int {

	resp, err := client.Get(url)
	if err != nil {
		return 0
	}

	defer resp.Body.Close()
	return resp.StatusCode
}
