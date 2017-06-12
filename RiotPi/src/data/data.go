/*
	author: Niclas Kristek
	github.com/nkristek
*/
package data

type SensorData struct {
	BoardId    int
	SensorType string
	Date       string
	Unit       string
	Value1     string
	Value2     string
	Value3     string
}

type Device struct {
	Id           int
	Address      string
	MessageCount int
	LastPing     string
}
