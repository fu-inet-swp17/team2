package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nkristek/go-senml"
)

type DataPoint struct {
	Id          int
	Name        *string
	Unit        *string
	Value       *float64
	StringValue *string
	BoolValue   *bool
	DataValue   *string
	Sum         *float64
	Time        *float64
	Link        *string
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

	// prepare data

	var name sql.NullString
	if message.Name != nil {
		name.String = *message.Name
		name.Valid = true
	}

	var unit sql.NullString
	if message.Unit != nil {
		unit.String = *message.Unit
		unit.Valid = true
	}

	var value sql.NullFloat64
	if message.Value != nil {
		value.Float64 = *message.Value
		value.Valid = true
	}

	var stringValue sql.NullString
	if message.StringValue != nil {
		stringValue.String = *message.StringValue
		stringValue.Valid = true
	}

	var boolValue sql.NullBool
	if message.BoolValue != nil {
		boolValue.Bool = *message.BoolValue
		boolValue.Valid = true
	}

	var dataValue sql.NullString
	if message.DataValue != nil {
		dataValue.String = *message.DataValue
		dataValue.Valid = true
	}

	var sum sql.NullFloat64
	if message.Sum != nil {
		sum.Float64 = *message.Sum
		sum.Valid = true
	}

	var timeStamp time.Time
	if message.Time != nil {
		timeStamp = time.Unix(int64(*message.Time), 0)
	} else {
		timeStamp = time.Now()
	}

	var link sql.NullString
	if message.Link != nil {
		link.String = *message.Link
		link.Valid = true
	}

	// execute insert
	_, err = stmt.Exec(name, unit, value, stringValue, boolValue, dataValue, sum, timeStamp.Format(time.RFC3339), link)
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
		var name sql.NullString
		var unit sql.NullString
		var value sql.NullFloat64
		var stringValue sql.NullString
		var boolValue sql.NullBool
		var dataValue sql.NullString
		var sum sql.NullFloat64
		var time sql.NullFloat64
		var link sql.NullString

		err = rows.Scan(&id, &name, &unit, &value, &stringValue, &boolValue, &dataValue, &sum, &time, &link)
		if err != nil {
			return nil, errors.New("Reading values from row: " + err.Error())
		}

		var nameString *string
		if name.Valid {
			nameString = &name.String
		}
		var unitString *string
		if unit.Valid {
			unitString = &unit.String
		}
		var valueFloat *float64
		if value.Valid {
			valueFloat = &value.Float64
		}
		var stringValueString *string
		if stringValue.Valid {
			stringValueString = &stringValue.String
		}
		var boolValueBool *bool
		if boolValue.Valid {
			boolValueBool = &boolValue.Bool
		}
		var dataValueString *string
		if dataValue.Valid {
			dataValueString = &dataValue.String
		}
		var sumFloat *float64
		if sum.Valid {
			sumFloat = &sum.Float64
		}
		var timeFloat *float64
		if time.Valid {
			timeFloat = &time.Float64
		}
		var linkString *string
		if link.Valid {
			linkString = &link.String
		}

		dataPoints = append(dataPoints, DataPoint{Id: id, Name: nameString, Unit: unitString, Value: valueFloat, StringValue: stringValueString, BoolValue: boolValueBool, DataValue: dataValueString, Sum: sumFloat, Time: timeFloat, Link: linkString})
	}

	return dataPoints, nil
}
