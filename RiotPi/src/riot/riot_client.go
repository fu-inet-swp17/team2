/*
	author: Niclas Kristek
	github.com/nkristek
*/
package riot

import (
	"../config"                 // config for Configuration struct
	"../data"                   // data for Device struct
	"../senml"                  // SenML parsing
	"errors"                    // create errors
	"fmt"                       // sprintf
	"github.com/dustin/go-coap" // coap
)

var messageCount int = 0

func GetSensorData(device data.Device, deviceResources []data.DeviceResource, configuration config.Configuration) ([]data.SensorData, error) {
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
			return nil, err
		}

		// send request and receive response
		resp, err := c.Send(req)
		if err != nil {
			log.Error("sending message: ", err)
			return nil, err
		}

		// update the message count
		messageCount++
		if messageCount > 65535 {
			messageCount = 0
		}

		// parse response
		if resp == nil {
			log.Warning("no response received")
			return nil, errors.New("no response received")
		}

		if len(resp.Payload) <= 0 {
			log.Warning("response payload is empty")
			return nil, errors.New("response payload is empty")
		}

		log.Debugf("response payload: " + string(resp.Payload[:]))

		parsedResponse, err := senml.ParseMessage(resp.Payload)
		if err != nil {
			log.Error("parsing SenML: ", err)
			return nil, err
		}

		log.Debugf("successfully parsed SenML")
		senml.LogMessage(*parsedResponse)
	}

	return nil, nil
}
