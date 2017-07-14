package config

type Configuration struct {
	PollingInterval     int
	ListeningPort       int
	RiotPort            int
	LowPanInterfaceName string
	SQL                 SQLSettings
}

type SQLSettings struct {
	Address         string
	Port            int
	User            string
	Password        string
	DatabaseName    string
	DeviceTableName string
	DataTableName   string
}
