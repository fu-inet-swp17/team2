/*
	author: Niclas Kristek
	github.com/nkristek
*/
package data

type Device struct {
	Id       int
	Address  string
	LastPing string
}

type DeviceResource struct {
	BoardId int
	Path    string
}

type SensorData struct {
	BoardId    int
	SensorType string
	Date       string
	Unit       string
	Value1     string
	Value2     string
	Value3     string
}
