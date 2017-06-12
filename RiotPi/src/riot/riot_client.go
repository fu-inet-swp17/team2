/*
	author: Niclas Kristek
	github.com/nkristek
*/
package riot

import (
	"../config"                 // config for Configuration struct
	"../data"                   // data for Device struct
	"fmt"                       // sprintf
	"github.com/dustin/go-coap" // coap
)

var messageCount int = 0

func GetInformation(device data.Device, deviceResources []data.DeviceResource, configuration config.Configuration) []data.SensorData {
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
			if len(resp.Payload) > 0 {
				log.Notice("response payload: " + string(resp.Payload[:]))
			} else {
				log.Warning("response payload is empty")
			}
		} else {
			log.Warning("no response")
		}
	}
	return nil
}
