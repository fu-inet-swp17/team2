package config

import (
	"testing"
	"os"
	"crypto/rand"
    "fmt"
)

func TestReadWriteSampleConfig(t *testing.T) {
	uuid, err := generateUuid()
	if err != nil {
		t.Fatalf("Generate UUID: ", err)
		return
	}

	_, err = WriteSampleConfig(uuid)
	if err != nil {
		t.Fatalf("Write config: ", err)
		_ = os.Remove(uuid)
		return
	}

	configuration, err := ReadConfig(uuid)
	if err != nil {
		t.Fatalf("Read config: ", err)
		_ = os.Remove(uuid)
		return
	}

	if !isEqualConfig(configuration, Default) {
		t.Fatalf("Configuration validation failed")
	}

	_ = os.Remove(uuid)
}

func generateUuid() (uuid string, err error) {
    b := make([]byte, 16)
    _, err = rand.Read(b)
    if err != nil {
        return
    }

    uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

    return
}

func isEqualConfig(a Configuration, b Configuration) bool {
	return a.PollingInterval == b.PollingInterval &&
	a.ListeningPort == b.ListeningPort &&
	a.RiotPort == b.RiotPort &&
	a.LowPanInterfaceName == b.LowPanInterfaceName &&
	a.SQL.Address == b.SQL.Address &&
	a.SQL.Port == b.SQL.Port &&
	a.SQL.User == b.SQL.User &&
	a.SQL.Password == b.SQL.Password &&
	a.SQL.DatabaseName == b.SQL.DatabaseName &&
	a.SQL.DeviceTableName == b.SQL.DeviceTableName &&
	a.SQL.DataTableName == b.SQL.DataTableName
}