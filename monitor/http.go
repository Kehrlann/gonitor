package monitor

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

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
