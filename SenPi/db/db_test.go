package db

import (
	"errors"
	"github.com/fu-inet-swp17/team2/SenPi/config"
	"github.com/nkristek/go-senml"
	"strconv"
	"testing"
)

func TestInitAndClearDatabase(t *testing.T) {
	// get test configuration
	configuration, err := getTestConfig()
	if err != nil {
		t.Fatalf("Get test config: " + err.Error())
		return
	}

	// init the db package
	Init(configuration.SQL)

	// init db
	err = InitDatabase()
	if err != nil {
		t.Fatalf("Initialise database: " + err.Error())
		return
	}

	// clear db
	err = ClearDatabase()
	if err != nil {
		t.Fatalf("Clear database: " + err.Error())
		return
	}
}

func TestDBDevice(t *testing.T) {
	// get test configuration
	configuration, err := getTestConfig()
	if err != nil {
		t.Fatalf("Get test config: " + err.Error())
		return
	}

	// init the db package
	Init(configuration.SQL)

	// init db
	err = InitDatabase()
	if err != nil {
		t.Fatalf("Initialise database: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// register a sample device
	sampleDeviceAddress := "SAMPLEADDRESS"
	err = RegisterDevice(sampleDeviceAddress)
	if err != nil {
		t.Fatalf("Register sample device: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// get all devices after register
	registeredDevices, err := GetRegisteredDevices()
	if err != nil {
		t.Fatalf("Get all registered devices: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// validate all devices after register
	if len(registeredDevices) != 1 {
		t.Fatalf("After register, there are " + strconv.Itoa(len(registeredDevices)) + " devices in db, expected is 1.")
		_ = ClearDatabase()
		return
	}
	registeredDevice := registeredDevices[0]

	// validate the registered device
	if registeredDevice.Address != sampleDeviceAddress {
		t.Fatalf("Address of registered device is not equal with the initial address.")
		_ = ClearDatabase()
		return
	}

	// update address of device
	registeredDevice.Address = "NEWSAMPLEADDRESS"
	err = registeredDevice.Update()
	if err != nil {
		t.Fatalf("Update registered device: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// get all devices after update
	updatedDevices, err := GetRegisteredDevices()
	if err != nil {
		t.Fatalf("Get all updated devices: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// validate all devices after update
	if len(updatedDevices) != 1 {
		t.Fatalf("After update, there are " + strconv.Itoa(len(updatedDevices)) + " devices in db, expected is 1.")
		_ = ClearDatabase()
		return
	}
	updatedDevice := updatedDevices[0]

	// validate the device
	if updatedDevice.Address != registeredDevice.Address {
		t.Fatalf("Address of updated device is not equal with the previous device.")
		_ = ClearDatabase()
		return
	}

	// unregister device
	err = updatedDevice.Unregister()
	if err != nil {
		t.Fatalf("Unregister updated device: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// get all devices after delete
	deletedDevices, err := GetRegisteredDevices()
	if err != nil {
		t.Fatalf("Get all deleted devices: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// validate all devices after delete
	if len(deletedDevices) != 0 {
		t.Fatalf("After delete, there are " + strconv.Itoa(len(deletedDevices)) + " devices in db, expected is 0.")
		_ = ClearDatabase()
		return
	}

	// clear db
	err = ClearDatabase()
	if err != nil {
		t.Fatalf("Clear database: " + err.Error())
		_ = ClearDatabase()
		return
	}
}

func TestDBDataPoint(t *testing.T) {
	// get test configuration
	configuration, err := getTestConfig()
	if err != nil {
		t.Fatalf("Get test config: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// init the db package
	Init(configuration.SQL)

	// init db
	err = InitDatabase()
	if err != nil {
		t.Fatalf("Initialise database: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// insert sample datapoint
	var sampleName string = "SENSORNAME"
	var sampleUnit string = "SENSORUNIT"
	var sampleValue float64 = 42
	sampleDataPoint := senml.SenMLRecord{Name: &sampleName, Unit: &sampleUnit, Value: &sampleValue}
	err = InsertDataPoint(sampleDataPoint)
	if err != nil {
		t.Fatalf("Insert sample datapoint: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// get all datapoints after insert
	dataPoints, err := GetDataPoints()
	if err != nil {
		t.Fatalf("Get all datapoints: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// validate datapoints after insert
	if len(dataPoints) != 1 {
		t.Fatalf("After insert, there are " + strconv.Itoa(len(dataPoints)) + " datapoints in db, expected is 1.")
		_ = ClearDatabase()
		return
	}
	dataPoint := dataPoints[0]

	// delete datapoint
	err = dataPoint.Delete()
	if err != nil {
		t.Fatalf("Delete datapoint: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// get all datapoints after delete
	deletedDataPoints, err := GetDataPoints()
	if err != nil {
		t.Fatalf("Get all datapoints: " + err.Error())
		_ = ClearDatabase()
		return
	}

	// validate all datapoints after delete
	if len(deletedDataPoints) != 0 {
		t.Fatalf("After insert, there are " + strconv.Itoa(len(deletedDataPoints)) + " datapoints in db, expected is 0.")
		_ = ClearDatabase()
		return
	}

	// clear db
	err = ClearDatabase()
	if err != nil {
		t.Fatalf("Clear database: " + err.Error())
		_ = ClearDatabase()
		return
	}
}

func getTestConfig() (config.Configuration, error) {
	// get the default path
	path, err := config.GetDefaultPath()
	if err != nil {
		return config.Configuration{}, errors.New("Get default path: " + err.Error())
	}
	path = path + "test.json"

	// read test config and init the db package
	configuration, err := config.ReadConfig(path)
	if err != nil {
		return configuration, errors.New("Reading config: " + err.Error())
	}

	return configuration, nil
}
