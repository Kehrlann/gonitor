package config

import (
	"encoding/json"
	"io/ioutil"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

var DEFAULT_CONFIG_PATH string = "./gonitor.config.json"

// Configuration is the application config
type Configuration struct {
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


// IsValid tells you whether you can trust this Smtp config to send an e-mail
func (smtp *Smtp) IsValid() bool {
	return smtp.Host != "" && smtp.Port != 0 && len(smtp.To) > 0 && smtp.FromAddress != ""
}

// FormatFromHeader creates the From header used in SMTP messages
func (smtp *Smtp) FormatFromHeader() string {
	return fmt.Sprintf("%v <%v>", smtp.FromName, smtp.FromAddress)
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


// LoadConfig loads a config from a JSON file
func LoadConfig(path string) (*Configuration, error) {
	log.Infof("Loading config from `%v` ...", path)
	file, err := ioutil.ReadFile(path)

	if err != nil {
		if path == DEFAULT_CONFIG_PATH {
			return nil, NewDefaultConfigError()
		}
		return nil, err
	}

	ret := &Configuration{}
	if err := json.Unmarshal(file, ret); err != nil {
		return nil, err
	}
	log.Info("Configuration loaded !")
	ret.LogConfig()
	return ret, nil
}

// LogConfig logs the config at the info level
func (config *Configuration) LogConfig() {
	smtp_validity := "valid"
	if !config.Smtp.IsValid() {
		smtp_validity = "INVALID"
	}

	log.Info()
	log.Info("Configuration is :")
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
