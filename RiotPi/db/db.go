package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fu-inet-swp17/team2/RiotPi/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
)

// vars

var log = logging.MustGetLogger("db")

var sqlConfiguration config.SQLSettings

// functions

func Init(settings config.SQLSettings) {
	sqlConfiguration = settings
}

func InitDatabase() error {
	// create connection
	log.Notice("Initialize database")

	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return errors.New("Connecting to server: " + err.Error())
	}
	defer conn.Close()

	// create database
	_, err = conn.Exec("CREATE DATABASE " + sqlConfiguration.DatabaseName)
	if err != nil {
		return errors.New("Creating database: " + err.Error())
	}

	// use database
	_, err = conn.Exec("USE " + sqlConfiguration.DatabaseName)
	if err != nil {
		return errors.New("Using created database: " + err.Error())
	}

	// create device table
	_, err = conn.Exec("CREATE TABLE `" + sqlConfiguration.DeviceTableName + "` (ID INT PRIMARY KEY AUTO_INCREMENT, Address varchar(255) NOT NULL, FailedAttempts INT DEFAULT 0);")
	if err != nil {
		return errors.New("Creating device table: " + err.Error())
	}

	// create data table
	_, err = conn.Exec("CREATE TABLE `" + sqlConfiguration.DataTableName + "` (ID INT PRIMARY KEY AUTO_INCREMENT, SensorType VARCHAR(255) NULL DEFAULT NULL, Date DATE NULL DEFAULT NULL, Unit VARCHAR(20) NULL DEFAULT NULL, Value1 VARCHAR(255) NULL DEFAULT NULL, Value2 VARCHAR(255) NULL DEFAULT NULL, Value3 VARCHAR(255) NULL DEFAULT NULL);")
	if err != nil {
		return errors.New("Creating data table: " + err.Error())
	}

	log.Notice("Database initialized")
	return nil
}

func ClearDatabase() error {
	// create connection
	log.Notice("Clearing database")

	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/", sqlConfiguration.User, sqlConfiguration.Password, sqlConfiguration.Address, sqlConfiguration.Port)
	conn, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return errors.New("Connecting to server: " + err.Error())
	}
	defer conn.Close()

	// drop database
	_, err = conn.Exec("DROP DATABASE " + sqlConfiguration.DatabaseName)
	if err != nil {
		return errors.New("Dropping database: " + err.Error())
	}

	log.Notice("Database cleared")
	return nil
}
