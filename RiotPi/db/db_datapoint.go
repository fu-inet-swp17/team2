package db

import (
	"errors"
	"fmt"
	//"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nkristek/go-senml"
)

type DataPoint struct {
	Id         int
	SensorType string
	Date       string
	Unit       string
	Value1     string
	Value2     string
	Value3     string
}

// TODO: Parse date and save full SenMLRecord
func InsertDataPoint(message senml.SenMLRecord) error {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return errors.New("Creating connection to database: " + err.Error())
	}
	defer conn.Close()

	// prepare insert
	_, err = conn.Prepare("INSERT " + sqlConfiguration.DataTableName + " SET SensorType=?,Date=?,Unit=?,Value1=?,Value2=?,Value3=?")
	if err != nil {
		return errors.New("Preparing insert statement: " + err.Error())
	}

	// execute insert
	/*
		var date time.Time
		_, err = stmt.Exec(message.Name, date, message.Unit, dataPoint.Value1, dataPoint.Value2, dataPoint.Value3)
		if err != nil {
			return errors.New("Executing insert statement: " + err.Error())
		}
	*/

	return nil
}

func (dataPoint DataPoint) Delete() error {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return errors.New("Creating connection to database: " + err.Error())
	}
	defer conn.Close()

	// prepare delete
	stmt, err := conn.Prepare("DELETE FROM " + sqlConfiguration.DataTableName + " WHERE Id=?")
	if err != nil {
		return errors.New("Preparing delete statement: " + err.Error())
	}

	// execute delete
	_, err = stmt.Exec(dataPoint.Id)
	if err != nil {
		return errors.New("Executing delete statement: " + err.Error())
	}

	return nil
}

func GetDataPoints() ([]DataPoint, error) {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return nil, errors.New("Creating connection to database: " + err.Error())
	}
	defer conn.Close()

	// select
	rows, err := conn.Query("SELECT * FROM " + sqlConfiguration.DataTableName)
	if err != nil {
		return nil, errors.New("Querying from database: " + err.Error())
	}

	// iterate over results
	var dataPoints []DataPoint
	for rows.Next() {
		var id int
		var sensorType string
		var date string
		var unit string
		var value1 string
		var value2 string
		var value3 string

		err = rows.Scan(&id, &sensorType, &date, &unit, &value1, &value2, &value3)
		if err != nil {
			return nil, errors.New("Reading values from row: " + err.Error())
		}

		dataPoints = append(dataPoints, DataPoint{Id: id, SensorType: sensorType, Date: date, Unit: unit, Value1: value1, Value2: value2, Value3: value3})
	}

	return dataPoints, nil
}
