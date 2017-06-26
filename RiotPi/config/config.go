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
		log.Critical("reading config file: ", err)
		panic(err)
	}
	defer file.Close()

	// parsing json
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Critical("parsing config file: ", err)
		panic(err)
	}

	return configuration
}

func WriteSampleConfig() {
	// check if file aleady exists
	if _, err := os.Stat("./conf.json"); err == nil {
		log.Error("file exists")
		return
	}

	// create sample config
	configuration := &Configuration{
		PollingInterval:     10,
		ListeningPort:       1234,
		RiotPort:            5683,
		LowPanInterfaceName: "lowpan0",
		SQL: SQLSettings{
			Address:                 "sqladdr",
			Port:                    1234,
			User:                    "sqluser",
			Password:                "sqlpw",
			DatabaseName:            "sqldb",
			DeviceTableName:         "sqldevicetable",
			DeviceResourceTableName: "sqldeviceresourcetable",
			DataTableName:           "sqldatatable",
		},
	}

	// marshal to json
	jsonConfiguration, err := json.Marshal(configuration)
	if err != nil {
		log.Error("marshaling json: %s", err)
		return
	}

	// creating file
	file, err := os.Create("./conf.json")
	if err != nil {
		log.Error("creating file: %s", err)
		return
	}
	defer file.Close()

	// write to file
	_, err = file.Write(jsonConfiguration)
	if err != nil {
		log.Error("writing to file: %s", err)
		return
	}

	log.Notice("config file created")
}
