package main

import (
	log "github.com/Sirupsen/logrus"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = printUsage
	path := flag.String("config", DEFAULT_CONFIG_PATH, "Path to the config file")
	flag.Parse()

	log.SetLevel(log.DebugLevel)
	log.Info("Starting Gonitor ...")
	config, err := LoadConfig(*path)

	if err != nil {
		switch err := err.(type) {
		default:
			log.Fatalf("Error loading config : `%v`", err)
			break
		case *NoDefaultConfigError:
			fmt.Println(err.HelpMessage)
			os.Exit(1)
		}
	}

	log.Info("Starting monitoring ...")
	messages := make(chan *StateChangeMessage)
	for _, resource := range config.Resources {
		go Run(resource, messages)
	}
	emitMessages(messages, &config.Smtp)
}

func emitMessages(messages <-chan *StateChangeMessage, smtp *Smtp) {
	for m := range messages {
		log.Info(m)
		go SendMail(smtp, m)
	}
}

func printUsage() {
	fmt.Println("Gonitor is a website monitoring tool.")
	fmt.Println()
	fmt.Println("Usage: gonitor [-config path_to_config]")
	fmt.Println()
	flag.PrintDefaults()
}
