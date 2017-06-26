/*
	author: Niclas Kristek
	github.com/nkristek
*/
package riot

import (
	"fmt"                                          // sprintf
	"github.com/dustin/go-coap"                    // coap
	"github.com/fu-inet-swp17/team2/RiotPi/config" // config for Configuration struct
	"github.com/fu-inet-swp17/team2/RiotPi/data"   // data for Device struct
	"github.com/nkristek/go-senml"                 // SenML parsing
)

var messageCount int = 0

func GetSensorData(device data.Device, deviceResources []data.DeviceResource, configuration config.Configuration) (data []data.SensorData) {
	log.Debugf("Get sensor data from: " + device.Address)
	for _, resource := range deviceResources {
		// create the request
		req := coap.Message{
			Type:      coap.Confirmable,
			Code:      coap.GET,
			MessageID: uint16(messageCount),
			Payload:   nil,
		}
		req.SetPathString(resource.Path)

		// dial the address
		coapinfo := fmt.Sprintf("[%s%%%s]:%d", device.Address, configuration.LowPanInterfaceName, configuration.RiotPort)
		c, err := coap.Dial("udp", coapinfo)
		if err != nil {
			log.Error("dial address: ", err)
			continue
		}

		// send request and receive response
		resp, err := c.Send(req)
		if err != nil {
			log.Error("sending message: ", err)
			continue
		}

		// update the message count
		messageCount++
		if messageCount > 65535 {
			messageCount = 0
		}

		// parse response
		if resp == nil {
			log.Warning("no response received")
			continue
		}

		if len(resp.Payload) <= 0 {
			log.Warning("response payload is empty")
			continue
		}

		log.Debugf("response payload: " + string(resp.Payload[:]))

		message, err := senml.Decode(resp.Payload, senml.JSON)
		if err != nil {
			log.Error("parsing SenML: ", err)
			continue
		}

		resolvedMessage, err := message.Resolve()
		if err != nil {
			log.Error("resolving SenML: ", err)
			continue
		}

		log.Debugf("successfully parsed SenML")
		logMessage(resolvedMessage)
	}

	return
}

func logMessage(message senml.SenMLMessage) {
	if len(message.Records) == 0 {
		log.Debugf("no records")
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
