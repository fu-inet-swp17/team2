package riot

import (
	"errors"
	"fmt"
	"strings"
	"github.com/dustin/go-coap"
	"github.com/nkristek/go-senml"
	"github.com/fu-inet-swp17/team2/SenPi/config"
	"github.com/fu-inet-swp17/team2/SenPi/db"
)

func ExecuteCoAPGetRequest(address string, port int, interfaceName string, path string) ([]byte, error) {
	// dial the address
	coapInfo := fmt.Sprintf("[%s%%%s]:%d", address, interfaceName, port)
	conn, err := coap.Dial("udp", coapInfo)
	if err != nil {
		return nil, errors.New("Dial address: " + err.Error())
	}

	// create the request
	request := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 0,
		Payload:   nil,
	}
	request.SetPathString(path)

	// send request and receive response
	response, err := conn.Send(request)
	if err != nil {
		return nil, errors.New("Sending message: " + err.Error())
	}

	// validate response
	if response == nil {
		return nil, errors.New("No response received")
	}
	if len(response.Payload) <= 0 {
		return nil, errors.New("Response payload is empty")
	}

	return response.Payload, nil
}

func GetCoAPResources(address string, port int, interfaceName string) ([]string, error) {
	// get the resources formatted like <PATH1>,<PATH2>,<PATH3>
	responsePayload, err := ExecuteCoAPGetRequest(address, port, interfaceName, coapWellKnownResourcesPath)
	if err != nil {
		return nil, err
	}

	// parse response
	return ParseCoAPResourceResponse(responsePayload), nil
}

func ParseCoAPResourceResponse(response []byte) (deviceResources []string) {
	for _, resourcePath := range strings.Split(string(response[:]), ",") {
		resourcePath = strings.TrimPrefix(resourcePath, "<")
		resourcePath = strings.TrimSuffix(resourcePath, ">")
		deviceResources = append(deviceResources, resourcePath)
	}
	return
}

func GetDataPoints(address string, port int, interfaceName string, deviceResources []string) ([]senml.SenMLMessage, error) {
	var messages []senml.SenMLMessage
	for _, resourcePath := range deviceResources {
		// get the SenML formatted response
		responsePayload, err := ExecuteCoAPGetRequest(address, port, interfaceName, resourcePath)
		if err != nil {
			return nil, err
		}

		// parse as SenML
		message, err := senml.Decode(responsePayload, senml.JSON)
		if err != nil {
			log.Warning("Parsing SenML (will skip resource): " + err.Error())
			continue
		}
		resolvedMessage, err := message.Resolve()
		if err != nil {
			log.Warning("Resolving SenML (will skip resource): " + err.Error())
			continue
		}

		messages = append(messages, resolvedMessage)
	}

	return messages, nil
}

func GetDataPointsFromDevice(device db.Device, configuration config.Configuration) ([]senml.SenMLMessage, error) {
	deviceResources, err := GetCoAPResources(device.Address, configuration.RiotPort, configuration.LowPanInterfaceName)
	if err != nil {
		// increase failed attempts
		device.FailedAttempts++
		updateErr := device.Update()
		if updateErr != nil {
			// just log the error from the db
			log.Warning(updateErr)
		}

		return nil, errors.New("Get device resources: " + err.Error())
	}

	return GetDataPoints(device.Address, configuration.RiotPort, configuration.LowPanInterfaceName, deviceResources)
}
