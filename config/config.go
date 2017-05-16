package config

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

var DEFAULT_CONFIG_PATH string = "./gonitor.config.json"

// Configuration is the application config
type Configuration struct {
	Smtp          Smtp
	GlobalCommand string
	Resources     []Resource
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
	ret.logConfig()
	return ret, nil
}

// logConfig logs the config at the info level
func (config *Configuration) logConfig() {
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
