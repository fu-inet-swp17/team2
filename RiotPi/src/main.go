/*
	author: Niclas Kristek
	github.com/nkristek
*/
package main

import (
	"./config"                 // configuration file handling
	"./data"                   // data struct
	"./db"                     // database handling
	"./riot"                   // communication with riot
	"flag"                     // command line argument parsing
	"github.com/op/go-logging" // logging
	"os"                       // file handling and command line arguments
	"runtime"                  // exit from another function
)

// vars

var log = logging.MustGetLogger("main")

// functions

func readCommandLineArgs() config.Configuration {
	// prepare command line arguments
	configPathPtr := flag.String("config", "conf.json", "Path to the config file")
	sampleConfigPtr := flag.Bool("sampleconfig", false, "If a conf.json file with default values should be created in the current directory")
	initDBPtr := flag.Bool("initdb", false, "If the database should be initialized using the configuration provided in the conf.json file")
	clearDBPtr := flag.Bool("cleardb", false, "If the database should be cleared")

	// parse command line arguments
	flag.Parse()

	// write sample config if set and exit
	if *sampleConfigPtr {
		config.WriteSampleConfig()
		log.Notice("server stopping")
		runtime.Goexit()
	}

	// read config
	configuration := config.ReadConfig(*configPathPtr)

	// initialize db package
	db.InitWithConfiguration(configuration)

	if *clearDBPtr {
		db.ClearDatabase()
	}

	// init database if set
	if *initDBPtr {
		db.InitDatabase()
	}

	return configuration
}

// main

func main() {
	defer os.Exit(0)

	// setup logging
	logging.SetFormatter(logging.MustStringFormatter(`%{color}%{level} > %{shortfunc}: %{color:reset} %{message}`))

	// start initializing
	log.Notice("server starting")

	// read the configuration
	configuration := readCommandLineArgs()

	completion := make(chan struct{})
	riot.StartListeningForDevices(configuration, completion)

	// start polling for information
	cancelPolling := make(chan struct{})
	riot.StartScheduledPolling(configuration, cancelPolling)

	// TODO: remove lines after debug
	log.Debugf("registering devices for debug")
	db.RegisterDevice("fe80::e870:6d5b:5292:432a")
	db.RegisterDeviceResource("/se-app/sensors", data.Device{Id: 1, Address: "fe80::e870:6d5b:5292:432a", LastPing: ""})
	db.RegisterDevice("fe80::7b76:6260:b7e3:dbce")
	db.RegisterDeviceResource("/se-app/sensors", data.Device{Id: 2, Address: "fe80::7b76:6260:b7e3:dbce", LastPing: ""})

	// initializing finished
	log.Notice("server started")

	<-completion

	cancelPolling <- struct{}{}

	log.Notice("server stopping")
}
