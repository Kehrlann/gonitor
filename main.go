package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kehrlann/gonitor/config"
)

func main() {
	flag.Usage = printUsage
	path := flag.String("config", config.DEFAULT_CONFIG_PATH, "Path to the config file")
	flag.Parse()

	log.SetLevel(log.DebugLevel)
	log.Info("Starting Gonitor ...")
	configuration, err := config.LoadConfig(*path)

	if err != nil {
		switch err := err.(type) {
		default:
			log.Fatalf("Error loading config : `%v`", err)
			break
		case *config.NoDefaultConfigError:
			fmt.Println(err.HelpMessage)
			os.Exit(1)
		}
	}

	log.Info("Starting monitoring ...")
	messages := make(chan *StateChangeMessage)
	for _, resource := range configuration.Resources {
		go Run(resource, messages)
	}
	emitMessages(messages, &configuration.Smtp, configuration.GlobalCommand)
}

func emitMessages(messages <-chan *StateChangeMessage, smtp *config.Smtp, globalCommand string) {
	for m := range messages {
		log.Info(m)
		go SendMail(smtp, m)
		go ExecCommand(m, globalCommand)
	}
}

func printUsage() {
	fmt.Println("Gonitor is a website monitoring tool.")
	fmt.Println()
	fmt.Println("Usage: gonitor [-config path_to_config]")
	fmt.Println()
	flag.PrintDefaults()
}
