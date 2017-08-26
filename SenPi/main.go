package main

import (
	"flag"
	"github.com/fu-inet-swp17/team2/SenPi/config"
	"github.com/fu-inet-swp17/team2/SenPi/db"
	"github.com/fu-inet-swp17/team2/SenPi/riot"
	"github.com/op/go-logging"
	"os"
	"runtime"
	"time"
)

var log = logging.MustGetLogger("main")

func readCommandLineArgs() config.Configuration {
	// prepare command line arguments
	configPathPtr := flag.String("config", "", "Path to the config file")
	sampleConfigPtr := flag.Bool("sampleconfig", false, "If a conf.json file with default values should be created in the current directory")
	initDBPtr := flag.Bool("initdb", false, "If the database should be initialized using the configuration provided in the conf.json file")
	clearDBPtr := flag.Bool("cleardb", false, "If the database should be cleared")
	reinitDBPtr := flag.Bool("reinitdb", false, "If the database should be cleared and initialized afterwards. Basically combines cleardb and initdb.")

	// parse command line arguments
	flag.Parse()

	// write sample config if set and exit
	if *sampleConfigPtr {
		filePath, err := config.WriteSampleConfig(*configPathPtr)
		if err != nil {
			log.Critical(err.Error())
			runtime.Goexit()
		}
		log.Notice("Config file created at \"" + filePath + "\"")
		log.Notice("Server stopping")
		runtime.Goexit()
	}

	// read config
	configuration, err := config.ReadConfig(*configPathPtr)
	if err != nil {
		log.Critical(err.Error())
		runtime.Goexit()
	}

	// initialize db package
	db.Init(configuration.SQL)

	if *reinitDBPtr {
		*clearDBPtr = true
		*initDBPtr = true
	}

	if *clearDBPtr {
		log.Notice("Clearing database")
		err = db.ClearDatabase()
		if err != nil {
			log.Critical(err.Error())
			runtime.Goexit()
		}
		log.Notice("Database cleared")
	}

	if *initDBPtr {
		log.Notice("Initialize database")
		err = db.InitDatabase()
		if err != nil {
			log.Critical(err.Error())
			runtime.Goexit()
		}
		log.Notice("Database initialized")
	}

	return configuration
}

func main() {
	defer os.Exit(0)

	// setup logging
	logging.SetFormatter(logging.MustStringFormatter(`%{color}%{level} > %{shortfunc}: %{color:reset} %{message}`))

	// start initializing
	log.Notice("Server starting")

	// read the configuration
	configuration := readCommandLineArgs()

	// start listening for devices to register
	completion := make(chan struct{})
	go riot.StartListeningForAnnouncements(configuration, completion)

	// start polling for information
	cancelPolling := make(chan struct{})
	ticker := time.NewTicker(time.Duration(configuration.PollingInterval) * time.Second)
	go riot.StartScheduledPolling(ticker, configuration, cancelPolling)

	// initializing finished
	log.Notice("Server started")

	// wait for the listening goroutine to exit
	<-completion

	log.Notice("Server stopping")

	// stop the scheduled polling goroutine
	cancelPolling <- struct{}{}
}
