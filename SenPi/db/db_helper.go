package db

import (
	"github.com/nkristek/go-senml"
)

// Helper method, will disappear after debugging
func logMessage(message senml.SenMLMessage) {
	if len(message.Records) == 0 {
		log.Debugf("No records")
		return
	}

	for _, record := range message.Records {
		log.Debugf("-----")
		if record.BaseName != nil {
			log.Debugf("BaseName: %s", *record.BaseName)
		}
		if record.BaseTime != nil {
			log.Debugf("BaseTime: %f", *record.BaseTime)
		}
		if record.BaseUnit != nil {
			log.Debugf("BaseUnit: %s", *record.BaseUnit)
		}
		if record.BaseValue != nil {
			log.Debugf("BaseValue: %f", *record.BaseValue)
		}
		if record.BaseSum != nil {
			log.Debugf("BaseSum: %f", *record.BaseSum)
		}
		if record.Version != nil {
			log.Debugf("Version: %d", *record.Version)
		}
		if record.Name != nil {
			log.Debugf("Name: %s", *record.Name)
		}
		if record.Unit != nil {
			log.Debugf("Unit: %s", *record.Unit)
		}
		if record.Value != nil {
			log.Debugf("Value: %f", *record.Value)
		}
		if record.StringValue != nil {
			log.Debugf("StringValue: %s", *record.StringValue)
		}
		if record.BoolValue != nil {
			log.Debugf("BoolValue: %t", *record.BoolValue)
		}
		if record.DataValue != nil {
			log.Debugf("DataValue: %s", *record.DataValue)
		}
		if record.Sum != nil {
			log.Debugf("Sum: %f", *record.Sum)
		}
		if record.Time != nil {
			log.Debugf("Time: %f", *record.Time)
		}
		if record.UpdateTime != nil {
			log.Debugf("UpdateTime: %f", *record.UpdateTime)
		}
		if record.Link != nil {
			log.Debugf("Link: %s", *record.Link)
		}
		log.Debugf("-----")
	}
}
