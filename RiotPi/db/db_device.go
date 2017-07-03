package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Device struct {
	Id             int
	Address        string
	FailedAttempts int
}

func RegisterDevice(address string) error {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return errors.New("Creating connection to database: " + err.Error())
	}
	defer conn.Close()

	// check if the device is already in the db
	rows, err := conn.Query("SELECT ID FROM "+sqlConfiguration.DeviceTableName+" WHERE Address=?", address)
	if err != nil {
		return errors.New("Checking if device already exists: " + err.Error())
	}
	if rows.Next() {
		return errors.New("Already in database")
	}

	// -> device not in database
	log.Notice("Register device with address: " + address)

	// prepare insert
	stmt, err := conn.Prepare("INSERT " + sqlConfiguration.DeviceTableName + " SET Address=?")
	if err != nil {
		return errors.New("Preparing insert statement: " + err.Error())
	}

	// execute insert
	_, err = stmt.Exec(address)
	if err != nil {
		return errors.New("Executing insert statement: " + err.Error())
	}

	return nil
}

func (device Device) Unregister() error {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return errors.New("Creating connection to database: " + err.Error())
	}
	defer conn.Close()

	// prepare delete
	stmt, err := conn.Prepare("DELETE FROM " + sqlConfiguration.DeviceTableName + " WHERE ID=?")
	if err != nil {
		return errors.New("Preparing delete statement: " + err.Error())
	}

	// execute delete
	_, err = stmt.Exec(device.Id)
	if err != nil {
		return errors.New("Executing delete statement: " + err.Error())
	}

	return nil
}

func (device Device) Update() error {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return errors.New("Creating connection to database: " + err.Error())
	}
	defer conn.Close()

	// prepare update
	stmt, err := conn.Prepare("UPDATE " + sqlConfiguration.DeviceTableName + " SET Address=?,FailedAttempts=? WHERE Id=?")
	if err != nil {
		return errors.New("Preparing update statement: " + err.Error())
	}

	// execute update
	_, err = stmt.Exec(device.Address, device.FailedAttempts, device.Id)
	if err != nil {
		return errors.New("Executing update statement: " + err.Error())
	}

	return nil
}

func GetRegisteredDevices() ([]Device, error) {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return nil, errors.New("Creating connection to database: " + err.Error())
	}
	defer conn.Close()

	// select
	rows, err := conn.Query("SELECT * FROM " + sqlConfiguration.DeviceTableName)
	if err != nil {
		return nil, errors.New("Querying from database: " + err.Error())
	}

	// iterate over results
	var devices []Device
	for rows.Next() {
		var id int
		var address string
		var failedAttempts int

		err = rows.Scan(&id, &address, &failedAttempts)
		if err != nil {
			return nil, errors.New("Reading values from row: " + err.Error())
		}

		devices = append(devices, Device{Id: id, Address: address, FailedAttempts: failedAttempts})
	}

	return devices, nil
}
