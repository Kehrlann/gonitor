package monitor

import (
	"time"

	"github.com/kehrlann/gonitor/emit"
	"github.com/kehrlann/gonitor/config"
)

// Run takes a resource and polls the given HTTP url for errors , and emits failure / recovery messages accordingly
func Run(resources []config.Resource, messages chan<- *emit.StateChangeMessage) {
	for _, resource := range resources {
		go run(resource, messages)
	}
}

func run(resource config.Resource, messages chan<- *emit.StateChangeMessage) {
	responseCodes := make(chan int)

	go analyze(resource, responseCodes, messages)

	// TODO : do we want to test that ?? if so, we need to mock the duration interval (val'd ?) and the fetch method
	for range time.Tick(time.Duration(resource.IntervalInSeconds) * time.Second) {
		responseCodes <- Fetch(resource)
	}
}
