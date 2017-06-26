/*
	author: Niclas Kristek
	github.com/nkristek
*/
package db

import (
	"database/sql"                                 // sql
	"fmt"                                          // sprintf
	"github.com/fu-inet-swp17/team2/RiotPi/config" // configuration
	_ "github.com/go-sql-driver/mysql"             // mysql driver
	"github.com/op/go-logging"                     // logging
)

// vars

var log = logging.MustGetLogger("db")

var sqlConfiguration config.SQLSettings

// functions

func InitWithConfiguration(configuration config.Configuration) {
	sqlConfiguration = configuration.SQL
}

func ClearDatabase() {
	// create connection
	log.Notice("clearing database")

	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Error("connecting to server: ", err)
		panic(err)
	}
	defer db.Close()

	// drop database
	_, err = db.Exec("DROP DATABASE " + sqlConfiguration.DatabaseName)
	if err != nil {
		log.Error("dropping database: ", err)
		panic(err)
	}

	log.Notice("database cleared")
}

func InitDatabase() {
	// create connection
	log.Notice("initialize database")

	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port)
	db, err := sql.Open("mysql", dbinfo)
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

	// create device resource table
	_, err = db.Exec("CREATE TABLE `" + sqlConfiguration.DeviceResourceTableName + "` (BoardId INT NOT NULL, ResourcePath varchar(255) NOT NULL);")
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
