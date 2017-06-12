/*
	author: Niclas Kristek
	github.com/nkristek
*/
package db

import (
	"../config"                        // configuration
	"database/sql"                     // sql
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/op/go-logging"         // logging
	"strconv"                          // atoi/itoa
)

// vars

var log = logging.MustGetLogger("dbhandler")

var sqlConfiguration config.SQLSettings

// functions

func InitWithConfiguration(configuration config.Configuration) {
	sqlConfiguration = configuration.SQL
}

func InitDatabase() {
	// create connection
	log.Notice("initialize database")

	sqlConnectString := sqlConfiguration.User + ":" + sqlConfiguration.Password + "@tcp(" + sqlConfiguration.Address + ":" + strconv.Itoa(sqlConfiguration.Port) + ")/"
	db, err := sql.Open("mysql", sqlConnectString)
	if err != nil {
		log.Error("connecting to server: ", err)
		panic(err)
	}
	defer db.Close()

	// create database
	_, err = db.Exec("CREATE DATABASE " + sqlConfiguration.DatabaseName)
	if err != nil {
		log.Error("creating database: ", err)
		panic(err)
	}

	// use database
	_, err = db.Exec("USE " + sqlConfiguration.DatabaseName)
	if err != nil {
		log.Error("using created database: ", err)
		panic(err)
	}

	// create device table
	_, err = db.Exec("CREATE TABLE `" + sqlConfiguration.DeviceTableName + "` (ID INT PRIMARY KEY AUTO_INCREMENT, Address varchar(255) NOT NULL, LastPing DATE NULL DEFAULT NULL);")
	if err != nil {
		log.Error("creating device table: ", err)
		panic(err)
	}

	// create data table
	_, err = db.Exec("CREATE TABLE `" + sqlConfiguration.DataTableName + "` (BoardId INT NOT NULL, SensorType VARCHAR(255) NULL DEFAULT NULL, Date DATE NULL DEFAULT NULL, Unit VARCHAR(20) NULL DEFAULT NULL, Value1 VARCHAR(255) NULL DEFAULT NULL, Value2 VARCHAR(255) NULL DEFAULT NULL, Value3 VARCHAR(255) NULL DEFAULT NULL);")
	if err != nil {
		log.Error("creating data table: ", err)
		panic(err)
	}

	log.Notice("database initialized")
}
