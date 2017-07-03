package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nkristek/go-senml"
)

type DataPoint struct {
	Id          int
	Name        string
	Unit        string
	Value       float64
	StringValue string
	BoolValue   bool
	DataValue   string
	Sum         float64
	Time        float64
	Link        string
}

func InsertDataPoint(message senml.SenMLRecord) error {
	// create connection
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port, sqlConfiguration.DatabaseName)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return errors.New("Creating connection to database: " + err.Error())
	}
	defer conn.Close()

	// prepare insert
	stmt, err := conn.Prepare("INSERT " + sqlConfiguration.DataTableName + " SET Name=?,Unit=?,Value=?,StringValue=?,BoolValue=?,DataValue=?,Sum=?,Time=?,Link=?")
	if err != nil {
		return errors.New("Preparing insert statement: " + err.Error())
	}

	// execute insert
	_, err = stmt.Exec(message.Name, message.Unit, message.Value, message.StringValue, message.BoolValue, message.DataValue, message.Sum, message.Time, message.Link)
	if err != nil {
		return errors.New("Executing insert statement: " + err.Error())
	}

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
		var name string
		var unit string
		var value float64
		var stringValue string
		var boolValue bool
		var dataValue string
		var sum float64
		var time float64
		var link string

		err = rows.Scan(&id, &name, &unit, &value, &stringValue, &boolValue, &dataValue, &sum, &time, &link)
		if err != nil {
			return nil, errors.New("Reading values from row: " + err.Error())
		}

		dataPoints = append(dataPoints, DataPoint{Id: id, Name: name, Unit: unit, Value: value, StringValue: stringValue, BoolValue: boolValue, DataValue: dataValue, Sum: sum, Time: time, Link: link})
	}

	return dataPoints, nil
}
