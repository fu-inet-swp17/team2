/*
	author: Niclas Kristek
	github.com/nkristek
*/
package db

import (
	"database/sql"                               // sql
	"fmt"                                        // sprintf
	"github.com/fu-inet-swp17/team2/RiotPi/data" // data for Device struct
	_ "github.com/go-sql-driver/mysql"           // mysql driver
)

func InsertSensorData(sensordata data.SensorData) {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("creating connection to database: ", err)
		return
	}
	defer db.Close()

	// prepare insert
	stmt, err := db.Prepare("INSERT " + sqlConfiguration.DataTableName + " SET BoardId=?,SensorType=?,Date=?,Unit=?,Value1=?,Value2=?,Value3=?")
	if err != nil {
		log.Error("error preparing insert statement: ", err)
		return
	}

	// execute insert
	_, err = stmt.Exec(sensordata.BoardId, sensordata.SensorType, sensordata.Date, sensordata.Unit, sensordata.Value1, sensordata.Value2, sensordata.Value3)
	if err != nil {
		log.Error("error executing insert statement: ", err)
		return
	}
}

func DeleteSensorData(boardId int) {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("creating connection to database: ", err)
		return
	}
	defer db.Close()

	// prepare delete
	stmt, err := db.Prepare("DELETE FROM " + sqlConfiguration.DataTableName + " WHERE BoardId=?")
	if err != nil {
		log.Error("error preparing delete statement: ", err)
		return
	}

	// execute delete
	_, err = stmt.Exec(boardId)
	if err != nil {
		log.Error("error executing delete statement: ", err)
		return
	}
}

func GetSensorData() []data.SensorData {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("creating connection to database: ", err)
		return nil
	}
	defer db.Close()

	// select
	rows, err := db.Query("SELECT * FROM " + sqlConfiguration.DataTableName)
	if err != nil {
		log.Error("querying from database: ", err)
		return nil
	}

	// iterate over results
	var sensorData []data.SensorData
	for rows.Next() {
		var boardId int
		var sensorType string
		var date string
		var unit string
		var value1 string
		var value2 string
		var value3 string
		err = rows.Scan(&boardId, &sensorType, &date, &unit, &value1, &value2, &value3)
		if err != nil {
			log.Warning("reading values from row: ", err)
			continue
		}

		log.Debugf("BoardId: ", boardId)
		log.Debugf("SensorType: ", sensorType)
		sensorData = append(sensorData, data.SensorData{BoardId: boardId, SensorType: sensorType, Date: date, Unit: unit, Value1: value1, Value2: value2, Value3: value3})
	}

	return sensorData
}
