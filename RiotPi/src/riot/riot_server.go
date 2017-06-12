/*
	author: Niclas Kristek
	github.com/nkristek
*/
package riot

import (
	"../config" // configuration
	"../db"     // get all devices
	//"github.com/dustin/go-coap" // coap
	"net"  // udp server
	"time" // timer
	//"strconv"	// atoi/itoa
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
	udpAddr, err := net.ResolveUDPAddr("udp", "[ff02::1]:2017")
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
			log.Error("reading: ", err)
		} else {
			log.Debugf("got message from " + remoteAddr.IP.String() + " with content: " + string(buffer[:]))
		}
	}
}

func StartScheduledPolling(pollingInterval int, cancel chan struct{}) {
	ticker := time.NewTicker(time.Duration(pollingInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go getInformationFromRegisteredDevices()
			case <-cancel:
				ticker.Stop()
				log.Info("stop listening")
				return
			}
		}
	}()
}

func getInformationFromRegisteredDevices() {
	log.Debugf("polling now")

	devices := db.GetRegisteredDevices()
	for _, device := range devices {
		GetInformation(device)
	}
}
