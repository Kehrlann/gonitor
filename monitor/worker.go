package monitor

import (
	"net/http"
	"time"

	"github.com/kehrlann/gonitor/emit"
	"github.com/kehrlann/gonitor/config"
)

// Run takes a resource and polls the given HTTP url for errors , and emits failure / recovery messages accordingly
func Run(resource config.Resource, messages chan<- *emit.StateChangeMessage) {
	responseCodes := make(chan int)

	client := &http.Client{
		Timeout: time.Duration(resource.TimeoutInSeconds) * time.Second,
	}

	go analyze(resource, responseCodes, messages)

	for range time.Tick(time.Duration(resource.IntervalInSeconds) * time.Second) {
		responseCodes <- fetch(client, resource.Url)
	}
}