/*
	author: Niclas Kristek
	github.com/nkristek
*/
package riot

import (
	"../config" // configuration
	"../db"     // get all devices
	"fmt"       // sprintf
	"net"       // udp server
	"time"      // timer
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
	//udpAddr, err := net.ResolveUDPAddr("udp6", ":" + strconv.Itoa(configuration.ListeningPort))
	udpinfo := fmt.Sprintf("[ff02::1]:%d", configuration.ListeningPort)
	udpAddr, err := net.ResolveUDPAddr("udp", udpinfo)
	if err != nil {
		log.Error("resolving udp address: ", err)
		return
	}

	// listen
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Error("listen to udp: ", err)
		return
	}
	defer conn.Close()

	// get all requests
	buffer := make([]byte, 1024)
	for {
		_, remoteAddr, err := conn.ReadFromUDP(buffer[0:])
		if err != nil {
			log.Error("reading message: ", err)
		} else {
			log.Debugf("got message from " + remoteAddr.IP.String() + " with content:\n" + string(buffer[:]))
		}
	}
}

func StartScheduledPolling(configuration config.Configuration, cancel chan struct{}) {
	ticker := time.NewTicker(time.Duration(configuration.PollingInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go getInformationFromRegisteredDevices(configuration)
			case <-cancel:
				ticker.Stop()
				log.Info("stop listening")
				return
			}
		}
	}()
}

func getInformationFromRegisteredDevices(configuration config.Configuration) {
	devices := db.GetRegisteredDevices()
	for _, device := range devices {
		deviceResources := db.GetRegisteredDeviceResources(device)
		GetSensorData(device, deviceResources, configuration)
	}
}
