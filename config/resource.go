package config

// Resource represents an Url to be watched, and all the necessary config/timing parameters
type Resource struct {
	Url               string // the Url to watch
	IntervalInSeconds int    // the interval at which to poll the resource. Defaults to 60s.
	TimeoutInSeconds  int    // the timeout in seconds. Defaults to 2s.
	NumberOfTries     int    // number of attempts at getting the resource. Defaults to 10.
	FailureThreshold  int    // the number of failures within the tries that would raise an alarm. Defaults to 5.
	Command           string // command to run whenever there is a failure or recovery message
}
