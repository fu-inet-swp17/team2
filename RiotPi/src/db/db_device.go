/*
	author: Niclas Kristek
	github.com/nkristek
*/
package db

import (
	"../data"                          // data for Device struct
	"database/sql"                     // sql
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"fmt"								// sprintf
	"strconv"                          // atoi/itoa
	"time"                             // time
)

func RegisterDevice(address string) {
	// create connection
	sqlConnectString := sqlConfiguration.User + ":" + sqlConfiguration.Password + "@tcp(" + sqlConfiguration.Address + ":" + strconv.Itoa(sqlConfiguration.Port) + ")/" + sqlConfiguration.DatabaseName
	db, err := sql.Open("mysql", sqlConnectString)
	if err != nil {
		log.Error("creating connection to database: ", err)
		return
	}
	defer db.Close()

	// check if the device is already in the db
	rows, err := db.Query("SELECT ID FROM " + sqlConfiguration.DeviceTableName + " WHERE Address=?", address)
	if err != nil {
		log.Error("checking if device already exists: ", err)
		return
	}

	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Debugf("Device already in db with ID=", id)
		} else {
			log.Debugf("Device already in db")
		}
		return
	}

	// prepare insert
	stmt, err := db.Prepare("INSERT " + sqlConfiguration.DeviceTableName + " SET Address=?,LastPing=?")
	if err != nil {
		log.Error("error preparing insert statement: ", err)
		return
	}

	// execute insert
	_, err = stmt.Exec(address, time.Now())
	if err != nil {
		log.Error("error executing insert statement: ", err)
		return
	}
}

func UnregisterDevice(deviceId int) {
	// create connection
	sqlConnectString := sqlConfiguration.User + ":" + sqlConfiguration.Password + "@tcp(" + sqlConfiguration.Address + ":" + strconv.Itoa(sqlConfiguration.Port) + ")/" + sqlConfiguration.DatabaseName
	db, err := sql.Open("mysql", sqlConnectString)
	if err != nil {
		log.Error("creating connection to database: ", err)
		return
	}
	defer db.Close()

	// prepare delete
	stmt, err := db.Prepare("DELETE FROM " + sqlConfiguration.DeviceTableName + " WHERE ID=?")
	if err != nil {
		log.Error("error preparing delete statement: ", err)
		return
	}

	// execute delete
	_, err = stmt.Exec(deviceId)
	if err != nil {
		log.Error("error executing delete statement: ", err)
		return
	}
}

func GetRegisteredDevices() []data.Device {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	log.Debugf("connecting using \n", dbinfo)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("creating connection to database: ", err)
		return nil
	}
	defer db.Close()

	// select
	rows, err := db.Query("SELECT * FROM " + sqlConfiguration.DeviceTableName)
	if err != nil {
		log.Error("querying from database: ", err)
		return nil
	}

	// iterate over results
	var devices []data.Device
	for rows.Next() {
		var id int
		var address string
		var lastPing string
		err = rows.Scan(&id, &address, &lastPing)
		if err != nil {
			log.Warning("reading values from row: ", err)
			continue
		}

		devices = append(devices, data.Device{Id: id, Address: address, LastPing: lastPing})
	}

	log.Debugf("got " + strconv.Itoa(len(devices)) + " devices")

	return devices
}
