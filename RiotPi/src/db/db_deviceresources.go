/*
	author: Niclas Kristek
	github.com/nkristek
*/
package db

import (
	"../data"                          // data for Device struct
	"database/sql"                     // sql
	"fmt"                              // sprintf
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"strconv"                          // atoi/itoa
)

func RegisterDeviceResource(path string, device data.Device) {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("creating connection to database: ", err)
		return
	}
	defer db.Close()

	// check if the device resource is already in the db
	rows, err := db.Query("SELECT * FROM "+sqlConfiguration.DeviceResourceTableName+" WHERE BoardId=? AND ResourcePath=?", device.Id, path)
	if err != nil {
		log.Error("checking if device resource already exists: ", err)
		return
	}
	if rows.Next() {
		log.Warning("device resource already in db")
		return
	}

	// prepare insert
	stmt, err := db.Prepare("INSERT " + sqlConfiguration.DeviceResourceTableName + " SET BoardId=?,ResourcePath=?")
	if err != nil {
		log.Error("error preparing insert statement: ", err)
		return
	}

	// execute insert
	_, err = stmt.Exec(device.Id, path)
	if err != nil {
		log.Error("error executing insert statement: ", err)
		return
	}
}

func UnregisterDeviceResource(path string, device data.Device) {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("creating connection to database: ", err)
		return
	}
	defer db.Close()

	// prepare delete
	stmt, err := db.Prepare("DELETE FROM " + sqlConfiguration.DeviceResourceTableName + " WHERE BoardId=? AND ResourcePath=?")
	if err != nil {
		log.Error("error preparing delete statement: ", err)
		return
	}

	// execute delete
	_, err = stmt.Exec(device.Id, path)
	if err != nil {
		log.Error("error executing delete statement: ", err)
		return
	}
}

func GetRegisteredDeviceResources(device data.Device) []data.DeviceResource {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("creating connection to database: ", err)
		return nil
	}
	defer db.Close()

	// select
	rows, err := db.Query("SELECT * FROM "+sqlConfiguration.DeviceResourceTableName+" WHERE BoardId=?", device.Id)
	if err != nil {
		log.Error("querying from database: ", err)
		return nil
	}

	// iterate over results
	var deviceResources []data.DeviceResource
	for rows.Next() {
		var boardId int
		var path string
		err = rows.Scan(&boardId, &path)
		if err != nil {
			log.Warning("reading values from row: ", err)
			continue
		}

		deviceResources = append(deviceResources, data.DeviceResource{BoardId: boardId, Path: path})
	}

	log.Debugf("got " + strconv.Itoa(len(deviceResources)) + " device resources")

	return deviceResources
}
