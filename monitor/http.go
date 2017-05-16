package monitor

import (
	"net/http"
	"time"

	"github.com/kehrlann/gonitor/emit"
	"github.com/kehrlann/gonitor/config"
	log "github.com/Sirupsen/logrus"
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

func fetch(client *http.Client, url string) int {
	log.Debugf("Getting %v", url)
	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		log.Warnf("Error fetching %v : `%v`", url, err)
		return 0
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	log.Debugf("%v, status code : %v, response time : %s", url, statusCode, time.Since(start))
	return statusCode
}
