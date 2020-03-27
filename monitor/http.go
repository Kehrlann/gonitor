package monitor

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/Kehrlann/gonitor/config"
)

type fetcher interface {
	fetch(url string) int
}

type httpFetcher struct {
	client *http.Client
}

// Fetch fetches the configured resource and gives back the HTTP response code
func Fetch(resource config.Resource) int{
	client := &http.Client{
		Timeout: time.Duration(resource.TimeoutInSeconds) * time.Second,
	}
	fetcher := &httpFetcher{client}
	return fetcher.fetch(resource.Url)
}

func (fetcher *httpFetcher) fetch(url string) int {
	log.Debugf("Getting %v", url)
	start := time.Now()
	resp, err := fetcher.client.Get(url)
	if err != nil {
		log.Warnf("Error fetching %v : `%v`", url, err)
		return 0
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	log.Debugf("%v, status code : %v, response time : %s", url, statusCode, time.Since(start))
	return statusCode
}

