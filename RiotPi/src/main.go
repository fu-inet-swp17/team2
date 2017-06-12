/*
	author: Niclas Kristek
	github.com/nkristek
*/
package main

import (
	"./config"                 // configuration file handling
	"./db"                     // database handling
	"./riot"                   // communication with riot
	"flag"                     // command line argument parsing
	"github.com/op/go-logging" // logging
	"os"                       // file handling and command line arguments
	"runtime"                  // exit from another function
	//"errors"      // combine errors from multiple goroutines
	//"sync"      	// waitgroup
)

// vars

var log = logging.MustGetLogger("main")

// functions

func readCommandLineArgs() config.Configuration {
	configPathPtr := flag.String("config", "conf.json", "Path to the config file")
	sampleConfigPtr := flag.Bool("sampleconfig", false, "If a conf.json file with default values should be created in the current directory")
	initDBPtr := flag.Bool("initdb", false, "If the database should be initialized using the configuration provided in the conf.json file")
	flag.Parse()

	if *sampleConfigPtr {
		config.WriteSampleConfig()
		log.Notice("server stopping")
		runtime.Goexit()
	}

	configuration := config.ReadConfig(*configPathPtr)

	db.InitWithConfiguration(configuration)

	if *initDBPtr {
		db.InitDatabase()
		// TODO: remove the next line after testing
		db.RegisterDevice("[fe80::e870:6d5b:5292:432a%lowpan0]:5683")
	}

	return configuration
}

// main

func main() {
	defer os.Exit(0)

	// setup logging
	logging.SetFormatter(logging.MustStringFormatter(`%{color}%{level} > %{shortfunc}: %{color:reset} %{message}`))

	log.Notice("server starting")

	// read the configuration
	configuration := readCommandLineArgs()

	completion := make(chan struct{})
	//riot.StartListeningForDevices(configuration, completion)

	// start polling for information
	cancelPolling := make(chan struct{})
	riot.StartScheduledPolling(configuration.PollingInterval, cancelPolling)

	log.Notice("server started")

	// execute stuff

	<-completion

	cancelPolling <- struct{}{}

	log.Notice("server stopping")
}
