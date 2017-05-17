package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kehrlann/gonitor/monitor/alert"
	"github.com/kehrlann/gonitor/config"
	log "github.com/Sirupsen/logrus"
	"github.com/kehrlann/gonitor/monitor"
)

func main() {
	flag.Usage = printUsage
	path := flag.String("config", config.DEFAULT_CONFIG_PATH, "Path to the config file")
	flag.Parse()

	log.SetLevel(log.DebugLevel)
	log.Info("Starting Gonitor ...")
	configuration, err := config.LoadConfig(*path)

	switch err := err.(type) {
	case nil:
		break
	case *config.NoDefaultConfigError:
		fmt.Println(err.HelpMessage)
		os.Exit(1)
	default:
		log.Fatalf("Error loading config : `%v`", err)
		break
	}

	log.Info("Starting monitoring ...")
	messages := make(chan *alert.StateChangeMessage)
	monitor.Monitor(configuration.Resources, messages)
	alert.EmitMessages(messages, configuration)
}

func printUsage() {
	fmt.Println("Gonitor is a website monitoring tool.")
	fmt.Println()
	fmt.Println("Usage: gonitor [-config path_to_config]")
	fmt.Println()
	flag.PrintDefaults()
}
