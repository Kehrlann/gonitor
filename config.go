package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

var DEFAULT_CONFIG_PATH string = "./gonitor.config.json"

// Config is the application config
type Config struct {
	Smtp          Smtp
	GlobalCommand string
	Resources     []Resource
}

type Smtp struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
	To          []string
}

// Resource represents an Url to be watched, and all the necessary config/timing parameters
type Resource struct {
	Url               string // the Url to watch
	IntervalInSeconds int    // the interval at which to poll the resource. Defaults to 60s.
	TimeoutInSeconds  int    // the timeout in seconds. Defaults to 2s.
	NumberOfTries     int    // number of attempts at getting the resource. Defaults to 10.
	FailureThreshold  int    // the number of failures within the tries that would raise an alarm. Defaults to 5.
	Command           string // command to run whenever there is a failure or recovery message
}

// NoDefaultConfigError represents an error thrown when the user hasn't specified a config, and that config wasn't found
type NoDefaultConfigError struct {
	HelpMessage string
}

// Error is the error message for NoDefautConfigError
func (n NoDefaultConfigError) Error() string {
	return fmt.Sprintf("No default config found at %v", DEFAULT_CONFIG_PATH)
}

func NewDefaultConfigError() *NoDefaultConfigError {
	return &NoDefaultConfigError{HelpMessage: fmt.Sprintf(`
It seems you didn't specify a config file. Gonitor attempted to load a default config file, located at  :
	%v

No such file was found. Please create one, or specify an existing config with the -config flag. If you wish to create 
one, here is a starter template. You can either replace the SMTP config with your own, or remove it entirely if you 
don't want e-mail notifications.

{
  "smtp"    :
  {
    "host"        : "smtp.example.com",
    "port"        : 25,
    "username"    : "user",
    "password"    : "password123",
    "fromaddress" : "address@example.com",
    "fromname"    : "Mr Example",
    "to"          : ["recipient@example.com", "admin@example.com"]
  },
  "resources" :
  [
	{
	  "url"               : "http://www.example.com",
	  "intervalInSeconds" : 60,
	  "timeoutInSeconds"  : 2,
	  "numberOfTries"     : 10,
	  "failureThreshold"  : 3
	},
	{
	  "url"               : "http://www.example.test",
	  "intervalInSeconds" : 120,
	  "timeoutInSeconds"  : 10,
	  "numberOfTries"     : 10,
	  "failureThreshold"  : 10
	}
  ]
}`, DEFAULT_CONFIG_PATH)}
}

// LoadConfig loads a config from a JSON file
func LoadConfig(path string) (*Config, error) {
	log.Infof("Loading config from `%v` ...", path)
	file, err := ioutil.ReadFile(path)

	if err != nil {
		if path == DEFAULT_CONFIG_PATH {
			return nil, NewDefaultConfigError()
		}
		return nil, err
	}

	ret := &Config{}
	if err := json.Unmarshal(file, ret); err != nil {
		return nil, err
	}
	log.Info("Config loaded !")
	ret.LogConfig()
	return ret, nil
}

// IsValid tells you whether you can trust this Smtp config to send an e-mail
func (smtp *Smtp) IsValid() bool {
	return smtp.Host != "" && smtp.Port != 0 && len(smtp.To) > 0 && smtp.FromAddress != ""
}

// FormatFromHeader creates the From header used in SMTP messages
func (smtp *Smtp) FormatFromHeader() string {
	return fmt.Sprintf("%v <%v>", smtp.FromName, smtp.FromAddress)
}

// LogConfig logs the config at the info level
func (config *Config) LogConfig() {
	smtp_validity := "valid"
	if !config.Smtp.IsValid() {
		smtp_validity = "INVALID"
	}

	log.Info()
	log.Info("Config is :")
	log.Infof(".. SMTP (%v) :", smtp_validity)
	log.Infof(".... 	Host         :    %v", config.Smtp.Host)
	log.Infof(".... 	Port         :    %v", config.Smtp.Port)
	log.Infof(".... 	Username     :    %v", config.Smtp.Username)
	log.Infof(".... 	Password     :    *** redacted ***")
	log.Infof(".... 	FromAddress  :    %v", config.Smtp.FromAddress)
	log.Infof(".... 	FromName     :    %v", config.Smtp.FromName)
	log.Infof(".... 	To           :    %v", config.Smtp.To)

	for i, resource := range config.Resources {
		log.Infof(".. Resource #%v :", i+1)
		log.Infof(".... 	Url               	:    %v", resource.Url)
		log.Infof(".... 	IntervalInSeconds 	:    %v", resource.IntervalInSeconds)
		log.Infof(".... 	TimeoutInSeconds  	:    %v", resource.TimeoutInSeconds)
		log.Infof(".... 	NumberOfTries     	:    %v", resource.NumberOfTries)
		log.Infof(".... 	FailureThreshold  	:    %v", resource.FailureThreshold)
	}
	log.Info()
}
