package riot

import (
	"errors"                                       // error handling
	"fmt"                                          // sprintf
	"github.com/dustin/go-coap"                    // coap
	"github.com/fu-inet-swp17/team2/RiotPi/config" // config for Configuration struct
	"github.com/fu-inet-swp17/team2/RiotPi/data"   // data for Device struct
	"github.com/fu-inet-swp17/team2/RiotPi/db"     // get all devices
	"github.com/nkristek/go-senml"                 // SenML parsing
	"strings"                                      // string formatting
)

var messageCount int = 0

const coapResourcePath string = "/.well-known/core"

func GetCoAPResources(device data.Device, configuration config.Configuration) (deviceResources []string, err error) {
	log.Debugf("Get device resources from address: " + device.Address)

	// create the request
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 0,
		Payload:   nil,
	}
	req.SetPathString(coapResourcePath)

	// dial the address
	coapinfo := fmt.Sprintf("[%s%%%s]:%d", device.Address, configuration.LowPanInterfaceName, configuration.RiotPort)
	c, err := coap.Dial("udp", coapinfo)
	if err != nil {
		err = errors.New("Dial address: " + err.Error())
		return
	}

	// send request and receive response
	resp, err := c.Send(req)
	if err != nil {
		err = errors.New("Sending message: " + err.Error())
		return
	}

	// parse response
	if resp == nil {
		err = errors.New("No response received")
		return
	}

	if len(resp.Payload) <= 0 {
		err = errors.New("Response payload is empty")
		return
	}

	response := string(resp.Payload[:])

	for _, resourcePath := range strings.Split(response, ",") {
		resourcePath = strings.TrimPrefix(resourcePath, "<")
		resourcePath = strings.TrimSuffix(resourcePath, ">")
		deviceResources = append(deviceResources, resourcePath)
	}

	return
}

func GetSensorData(device data.Device, deviceResources []string, configuration config.Configuration) (data []data.SensorData) {
	log.Debugf("Get sensor data from: " + device.Address)
	for _, resourcePath := range deviceResources {
		log.Debugf("Get resource: " + resourcePath)

		// create the request
		req := coap.Message{
			Type:      coap.Confirmable,
			Code:      coap.GET,
			MessageID: uint16(messageCount),
			Payload:   nil,
		}
		req.SetPathString(resourcePath)

		// dial the address
		coapinfo := fmt.Sprintf("[%s%%%s]:%d", device.Address, configuration.LowPanInterfaceName, configuration.RiotPort)
		c, err := coap.Dial("udp", coapinfo)
		if err != nil {
			log.Error("Dial address: ", err)
			continue
		}

		// send request and receive response
		resp, err := c.Send(req)
		if err != nil {
			log.Error("Sending message: ", err)
			continue
		}

		// update the message count
		messageCount++
		if messageCount > 65535 {
			messageCount = 0
		}

		// parse response
		if resp == nil {
			log.Warning("No response received")
			continue
		}

		if len(resp.Payload) <= 0 {
			log.Warning("Response payload is empty")
			continue
		}

		log.Debugf("Response payload: " + string(resp.Payload[:]))

		message, err := senml.Decode(resp.Payload, senml.JSON)
		if err != nil {
			log.Error("Parsing SenML: ", err)
			continue
		}

		resolvedMessage, err := message.Resolve()
		if err != nil {
			log.Error("Resolving SenML: ", err)
			continue
		}

		log.Debugf("Successfully parsed SenML")
		logMessage(resolvedMessage)
	}

	return
}

func GetSensorDataFromRegisteredDevices(configuration config.Configuration) {
	devices := db.GetRegisteredDevices()
	for _, device := range devices {
		deviceResources, err := GetCoAPResources(device, configuration)
		if err != nil {
			log.Error("Get device resources: ", err)
			continue
		}
		_ = GetSensorData(device, deviceResources, configuration)
	}
}

func logMessage(message senml.SenMLMessage) {
	if len(message.Records) == 0 {
		log.Debugf("No records")
		return
	}

	for _, record := range message.Records {
		log.Debugf("-----")
		if record.BaseName != nil {
			log.Debugf("BaseName: %s", *record.BaseName)
		}
		if record.BaseTime != nil {
			log.Debugf("BaseTime: %f", *record.BaseTime)
		}
		if record.BaseUnit != nil {
			log.Debugf("BaseUnit: %s", *record.BaseUnit)
		}
		if record.BaseValue != nil {
			log.Debugf("BaseValue: %f", *record.BaseValue)
		}
		if record.BaseSum != nil {
			log.Debugf("BaseSum: %f", *record.BaseSum)
		}
		if record.Version != nil {
			log.Debugf("Version: %d", *record.Version)
		}
		if record.Name != nil {
			log.Debugf("Name: %s", *record.Name)
		}
		if record.Unit != nil {
			log.Debugf("Unit: %s", *record.Unit)
		}
		if record.Value != nil {
			log.Debugf("Value: %f", *record.Value)
		}
		if record.StringValue != nil {
			log.Debugf("StringValue: %s", *record.StringValue)
		}
		if record.BoolValue != nil {
			log.Debugf("BoolValue: %t", *record.BoolValue)
		}
		if record.DataValue != nil {
			log.Debugf("DataValue: %s", *record.DataValue)
		}
		if record.Sum != nil {
			log.Debugf("Sum: %f", *record.Sum)
		}
		if record.Time != nil {
			log.Debugf("Time: %f", *record.Time)
		}
		if record.UpdateTime != nil {
			log.Debugf("UpdateTime: %f", *record.UpdateTime)
		}
		if record.Link != nil {
			log.Debugf("Link: %s", *record.Link)
		}
		log.Debugf("-----")
	}
}
