package riot

import (
	"fmt"                                          // sprintf
	"github.com/fu-inet-swp17/team2/RiotPi/config" // configuration
	"github.com/fu-inet-swp17/team2/RiotPi/db"     // register device
	"net"                                          // udp server
	"strings"                                      // hasprefix
	"time"                                         // timer
)

// listening for client announcements

func StartListeningForDevices(configuration config.Configuration, completion chan struct{}) {
	go func() {
		listenToLinkLocalMulticast(configuration)
		completion <- struct{}{}
	}()
}

func listenToLinkLocalMulticast(configuration config.Configuration) {
	// build the address
	udpinfo := fmt.Sprintf("[ff02::1]:%d", configuration.ListeningPort)
	udpAddr, err := net.ResolveUDPAddr("udp", udpinfo)
	if err != nil {
		log.Error("Resolving udp address: ", err)
		return
	}

	// listen
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Error("Listen to udp: ", err)
		return
	}
	defer conn.Close()

	// get all requests
	buffer := make([]byte, 1024)
	for {
		_, remoteAddress, err := conn.ReadFromUDP(buffer[0:])
		if err != nil {
			log.Error("Reading message: ", err)
		} else {
			handleResponse(buffer, remoteAddress.IP.String(), configuration)
		}
	}
}

func handleResponse(response []byte, remoteAddress string, configuration config.Configuration) {
	if strings.HasPrefix(string(response[:]), "riot-swp-2017-se") {
		db.RegisterDevice(remoteAddress)
	}
}

func StartScheduledPolling(configuration config.Configuration, cancel chan struct{}) {
	ticker := time.NewTicker(time.Duration(configuration.PollingInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go GetSensorDataFromRegisteredDevices(configuration)
			case <-cancel:
				ticker.Stop()
				log.Info("Stop listening")
				return
			}
		}
	}()
}
