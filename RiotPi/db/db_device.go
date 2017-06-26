package db

import (
	"database/sql"                               // sql
	"fmt"                                        // sprintf
	"github.com/fu-inet-swp17/team2/RiotPi/data" // data for Device struct
	_ "github.com/go-sql-driver/mysql"           // mysql driver
	"time"                                       // time
)

func RegisterDevice(address string) {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("Creating connection to database: ", err)
		return
	}
	defer db.Close()

	// check if the device is already in the db
	rows, err := db.Query("SELECT ID FROM "+sqlConfiguration.DeviceTableName+" WHERE Address=?", address)
	if err != nil {
		log.Error("Checking if device already exists: ", err)
		return
	}
	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		return
	}

	// -> device not in database
	log.Notice("Register device with address: " + address)

	// prepare insert
	stmt, err := db.Prepare("INSERT " + sqlConfiguration.DeviceTableName + " SET Address=?,LastPing=?")
	if err != nil {
		log.Error("Preparing insert statement: ", err)
		return
	}

	// execute insert
	_, err = stmt.Exec(address, time.Now())
	if err != nil {
		log.Error("Executing insert statement: ", err)
		return
	}
}

func UnregisterDevice(deviceId int) {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("Creating connection to database: ", err)
		return
	}
	defer db.Close()

	// prepare delete
	stmt, err := db.Prepare("DELETE FROM " + sqlConfiguration.DeviceTableName + " WHERE ID=?")
	if err != nil {
		log.Error("Preparing delete statement: ", err)
		return
	}

	// execute delete
	_, err = stmt.Exec(deviceId)
	if err != nil {
		log.Error("Executing delete statement: ", err)
		return
	}
}

func GetRegisteredDevices() []data.Device {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("Creating connection to database: ", err)
		return nil
	}
	defer db.Close()

	// select
	rows, err := db.Query("SELECT * FROM " + sqlConfiguration.DeviceTableName)
	if err != nil {
		log.Error("Querying from database: ", err)
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
			log.Warning("Reading values from row: ", err)
			continue
		}

		devices = append(devices, data.Device{Id: id, Address: address, LastPing: lastPing})
	}

	return devices
}
