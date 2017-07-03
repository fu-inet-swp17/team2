package riot

import (
	"errors"
	"fmt"
	"github.com/fu-inet-swp17/team2/RiotPi/config"
	"github.com/fu-inet-swp17/team2/RiotPi/db"
	"github.com/nkristek/go-senml"
	"net"
	"strings"
	"time"
)

func StartListeningForAnnouncements(configuration config.Configuration, completion chan struct{}) {
	err := listenToLinkLocalMulticast(configuration.ListeningPort)
	if err != nil {
		log.Error(err)
	}
	completion <- struct{}{}
}

func listenToLinkLocalMulticast(port int) error {
	// build the address
	udpinfo := fmt.Sprintf("[ff02::1]:%d", port)
	udpAddr, err := net.ResolveUDPAddr("udp", udpinfo)
	if err != nil {
		return errors.New("Resolving udp address: " + err.Error())
	}

	// listen
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return errors.New("Listen to udp: " + err.Error())
	}
	defer conn.Close()

	// handle requests
	buffer := make([]byte, 1024)
	for {
		_, remoteAddress, err := conn.ReadFromUDP(buffer[0:])
		if err != nil {
			log.Error("Reading message: " + err.Error())
			continue
		}

		err = handleRequest(remoteAddress.IP.String(), buffer)
		if err != nil {
			log.Error("Process message: " + err.Error())
			continue
		}
	}

	return nil
}

func handleRequest(remoteAddress string, response []byte) error {
	if strings.HasPrefix(string(response[:]), "riot-swp-2017-se") {
		return db.RegisterDevice(remoteAddress)
	}
	return nil
}

func StartScheduledPolling(ticker *time.Ticker, configuration config.Configuration, cancel chan struct{}) {
	for {
		select {
		case <-ticker.C:
			go pollRegisteredDevices(configuration)
		case <-cancel:
			ticker.Stop()
			log.Info("Stop listening")
			return
		}
	}
}

func pollRegisteredDevices(configuration config.Configuration) {
	// get all devices
	devices, err := db.GetRegisteredDevices()
	if err != nil {
		log.Error(err)
		return
	}

	log.Notice("Polling " + strconv.Itoa(len(devices)) + " devices.")

	// collect sensor data from all registered devices
	var messages []senml.SenMLMessage
	for _, device := range devices {
		currentMessages, err := GetDataPointsFromDevice(device, configuration)
		if err != nil {
			log.Error(err)
			continue
		}

		messages = append(messages, currentMessages...)
	}

	// save the datapoints
	for _, message := range messages {
		for _, dataPoint := range message.Records {
			err = db.InsertDataPoint(dataPoint)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
