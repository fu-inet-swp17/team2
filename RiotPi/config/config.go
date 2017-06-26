/*
	author: Niclas Kristek
	github.com/nkristek
*/
package config

import (
	"encoding/json"            // json formatting
	"github.com/op/go-logging" // logging
	"os"                       // file handling
)

// vars

var log = logging.MustGetLogger("config")

// functions

func ReadConfig(path string) Configuration {
	// opening file
	file, err := os.Open(path)
	if err != nil {
		log.Critical("Reading config file: ", err)
		panic(err)
	}
	defer file.Close()

	// parsing json
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Critical("Parsing config file: ", err)
		panic(err)
	}

	return configuration
}

func WriteSampleConfig() {
	// check if file aleady exists
	if _, err := os.Stat("./conf.json"); err == nil {
		log.Error("File exists")
		return
	}

	// create sample config
	configuration := &Configuration{
		PollingInterval:     360,
		ListeningPort:       2017,
		RiotPort:            5683,
		LowPanInterfaceName: "lowpan0",
		SQL: SQLSettings{
			Address:         "127.0.0.1",
			Port:            3306,
			User:            "SQLUSER",
			Password:        "SQLPW",
			DatabaseName:    "riotdata",
			DeviceTableName: "devices",
			DataTableName:   "sensordata",
		},
	}

	// marshal to json
	jsonConfiguration, err := json.Marshal(configuration)
	if err != nil {
		log.Error("Marshaling json: %s", err)
		return
	}

	// creating file
	file, err := os.Create("./conf.json")
	if err != nil {
		log.Error("Creating file: %s", err)
		return
	}
	defer file.Close()

	// write to file
	_, err = file.Write(jsonConfiguration)
	if err != nil {
		log.Error("Writing to file: %s", err)
		return
	}

	log.Notice("Config file created")
}
