package main

import (
	"net/http"
	"time"
)

var timeout = 1 * time.Second

// Run does blah
func Run(url string, every time.Duration, messages chan<- *Message) {
	responseCodes := make(chan int)
	go Analyze(url, responseCodes, messages)
	for _ = range time.Tick(every) {
		responseCodes <- fetch(url)
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
