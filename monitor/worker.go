package monitor

import (
	"time"

	"github.com/kehrlann/gonitor/emit"
	"github.com/kehrlann/gonitor/config"
)

// Run takes a resource and polls the given HTTP url for errors , and emits failure / recovery messages accordingly
func Run(resource config.Resource, messages chan<- *emit.StateChangeMessage) {
	responseCodes := make(chan int)

	go analyze(resource, responseCodes, messages)

	for range time.Tick(time.Duration(resource.IntervalInSeconds) * time.Second) {
		responseCodes <- Fetch(resource)
	}
}
//
//func run(resources []config.Resource, messages chan<- *emit.StateChangeMessage) {
//
//}
