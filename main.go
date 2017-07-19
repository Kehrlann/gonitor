package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kehrlann/gonitor/monitor/alert"
	"github.com/kehrlann/gonitor/config"
	log "github.com/sirupsen/logrus"
	"github.com/kehrlann/gonitor/monitor"
	"github.com/kehrlann/gonitor/server/web"
)

func main() {
	flag.Usage = printUsage
	path := flag.String("config", config.DEFAULT_CONFIG_PATH, "Path to the config file")
	example := flag.Bool("example", false, "Print an example config")
	flag.Parse()

	if *example {
		config.PrintExampleConfig()
		return
	}

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
	messages := monitor.Monitor(configuration.Resources)

	// This is where we do the glue
	// Options :
	// 1. web.serve returns a channel of ws  connections, the emitters eats it
	// 2. EmitMessages is non blocking and return a websocket connection handler, which is passed to serve ; and serve
	//		becomes blocking
	// 3. We create a broker that we pass to both, which mediates their interactions
	// 4. We pass the messages channel to the web.serve so it can use it to send the messages. Breaks the
	//		logical boundaries (web.serve creates some sort of emitters), but is clean regarding the glue code
	// "REGSITER CLIENT"
	incoming_connections := web.Serve()
	alert.EmitMessages(messages, incoming_connections, configuration)
}

func printUsage() {
	fmt.Println("Gonitor is a website monitoring tool.")
	fmt.Println()
	fmt.Println("Usage: gonitor [-config path_to_config]")
	fmt.Println()
	flag.PrintDefaults()
}
