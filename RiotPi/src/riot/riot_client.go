/*
	author: Niclas Kristek
	github.com/nkristek
*/
package riot

import (
	"../data"                   // data for Device struct
	"github.com/dustin/go-coap" // coap
	//"strconv"                  // atoi/itoa
)

var messageCount int = 0

func GetInformation(device data.Device) []data.SensorData {
	log.Debugf("get information from " + device.Address)

	// create the request
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: uint16(device.MessageCount),
		Payload:   nil,
	}
	req.SetPathString("/se-app/sensors")

	// dial the address
	c, err := coap.Dial("udp", device.Address)
	if err != nil {
		log.Error("dial address: ", err)
		return nil
	}

	// send request and receive response
	resp, err := c.Send(req)
	if err != nil {
		log.Error("sending message: ", err)
		return nil
	}

	// update the message count
	messageCount++
	if messageCount > 65535 {
		messageCount = 0
	}

	// parse response
	if resp != nil {
		log.Debugf("response payload: " + string(resp.Payload[:]))
	} else {
		log.Info("no response")
	}

	return nil
}
