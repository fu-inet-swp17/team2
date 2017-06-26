/*
	author: Niclas Kristek
	github.com/nkristek
*/
package main

import (
	"flag"                                         // command line argument parsing
	"github.com/fu-inet-swp17/team2/RiotPi/config" // configuration file handling
	"github.com/fu-inet-swp17/team2/RiotPi/data"   // data struct
	"github.com/fu-inet-swp17/team2/RiotPi/db"     // database handling
	"github.com/fu-inet-swp17/team2/RiotPi/riot"   // communication with riot
	"github.com/op/go-logging"                     // logging
	"os"                                           // file handling and command line arguments
	"runtime"                                      // exit from another function
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
	reinitDBPtr := flag.Bool("reinitdb", false, "If the database should be cleared and initialized afterwards. Basically combines cleardb and initdb.")

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

	if *reinitDBPtr {
		*clearDBPtr = true
		*initDBPtr = true
	}

	if *clearDBPtr {
		db.ClearDatabase()
	}

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
	/*
		db.RegisterDevice("fe80::e870:6d5b:5292:432a")
		db.RegisterDeviceResource("/se-app/sensors", data.Device{Id: 1, Address: "fe80::e870:6d5b:5292:432a", LastPing: ""})
	*/
	/*
		db.RegisterDevice("fe80::7b76:6260:b7e3:dbce")
		db.RegisterDeviceResource("/se-app/sensors", data.Device{Id: 1, Address: "fe80::7b76:6260:b7e3:dbce", LastPing: ""})
	*/

	db.RegisterDevice("fe80::9ab0:6d51:52a4:432a")
	db.RegisterDeviceResource("/se-app/temp", data.Device{Id: 1, Address: "fe80::9ab0:6d51:52a4:432a", LastPing: ""})
	//db.RegisterDeviceResource("/se-app/humid", data.Device{Id: 1, Address: "fe80::9ab0:6d51:52a4:432a", LastPing: ""})
	//db.RegisterDeviceResource("/se-app/mag", data.Device{Id: 1, Address: "fe80::9ab0:6d51:52a4:432a", LastPing: ""})

	// initializing finished
	log.Notice("server started")

	<-completion

	cancelPolling <- struct{}{}

	log.Notice("server stopping")
}