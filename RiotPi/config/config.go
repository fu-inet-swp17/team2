package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

var Default Configuration = Configuration{
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

// returns GOPATH/src/github.com/fu-inet-swp17/team2/RiotPi/
func GetDefaultPath() (path string, err error) {
	// get the go path
	gopath := os.Getenv("GOPATH")
	if len(gopath) <= 0 {
		err = errors.New("GOPATH not set")
		return
	}

	// remove trailing slash if present
	gopath = strings.TrimSuffix(gopath, "/")

	// build the default path
	path = gopath + "/src/github.com/fu-inet-swp17/team2/RiotPi/"

	return
}

func ReadConfig(path string) (configuration Configuration, err error) {
	// get default path
	if len(path) <= 0 {
		path, err = GetDefaultPath()
		if err != nil {
			err = errors.New("Get default path: " + err.Error())
			return
		}
		path = path + "conf.json"
	}

	// opening file
	file, err := os.Open(path)
	if err != nil {
		err = errors.New("Reading config file: " + err.Error())
		return
	}
	defer file.Close()

	// parsing json
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		err = errors.New("Parsing config file: " + err.Error())
		return
	}

	return
}

func WriteSampleConfig(path string) (filePath string, err error) {
	// get default path
	if len(path) <= 0 {
		path, err = GetDefaultPath()
		if err != nil {
			err = errors.New("Get default path: " + err.Error())
			return
		}
	}

	// check if file aleady exists
	if _, err = os.Stat(path); err == nil {
		err = errors.New("File exists")
		return
	}

	// marshal to json
	jsonConfiguration, err := json.MarshalIndent(Default, "", "    ")
	if err != nil {
		err = errors.New("Marshaling json: " + err.Error())
		return
	}

	// creating file
	file, err := os.Create(path)
	if err != nil {
		err = errors.New("Creating file: " + err.Error())
		return
	}
	defer file.Close()

	// write to file
	_, err = file.Write(jsonConfiguration)
	if err != nil {
		err = errors.New("Writing to file: " + err.Error())
		return
	}

	filePath = path

	return
}
