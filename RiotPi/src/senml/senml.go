/*
	author: Niclas Kristek
	github.com/nkristek
*/
package senml

/*
	This implements 'draft-ietf-core-senml-08' (https://tools.ietf.org/html/draft-ietf-core-senml-08)
*/

import (
	"encoding/json"            // json formatting
	"errors"                   // throw error when processing failed
	"github.com/op/go-logging" // logging
	"time"                     // get current time
)

// vars

var log = logging.MustGetLogger("senml")

// functions

func TestParseMessage() {
	var testData string = `[
     {"bn":"urn:dev:ow:10e2073a0108006:","bt":1.276020076001e+09,
      "bu":"A","bver":5,
      "n":"voltage","u":"V","v":120.1},
     {"n":"current","t":-5,"v":1.2},
     {"n":"current","t":-4,"v":1.3},
     {"n":"current","t":-3,"v":1.4},
     {"n":"current","t":-2,"v":1.5},
     {"n":"current","t":-1,"v":1.6},
     {"n":"current","v":1.7}
   ]`
	log.Debugf("testing with data:\n", testData)

	message, err := ParseMessage([]byte(testData))
	if err != nil {
		log.Error("testing failed: ", err)
		return
	}

	LogMessage(*message)
}

func LogMessage(message SenMLMessage) {
	if len(message.Records) == 0 {
		log.Debugf("no records")
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

func ParseMessage(payload []byte) (*SenMLMessage, error) {
	var message SenMLMessage

	err := json.Unmarshal(payload, &message.Records)
	if err != nil {
		log.Critical("parsing config file: ", err)
		return nil, err
	}

	ret, err := resolve(message)
	if err != nil {
		log.Critical("processing message: ", err)
		return nil, err
	}

	return ret, nil
}

// populates the base attributes in the regular attributes and deletes the base attributes in the process
// also deletes records with no value and sum (following the guidelines of the rfc)
func resolve(message SenMLMessage) (*SenMLMessage, error) {
	/*
				Each SenML Pack carries a single array that represents a set of
		   		measurements and/or parameters.  This array contains a series of
		   		SenML Records with several attributes described below.  There are two
		   		kind of attributes: base and regular.  The base attributes can be
		   		included in the any SenML Record and they apply to the entries in the
		   		Record.  Each base attribute also applies to all Records after it up
		   		to, but not including, the next Record that has that same base
		   		attribute.  All base attributes are optional.  Regular attributes can
		   		be included in any SenML Record and apply only to that Record.
	*/

	var resolvedMessage SenMLMessage

	var basename string = ""
	var basetime float64 = 0
	var baseunit string = ""
	var basevalue float64 = 0
	var basesum float64 = 0
	var baseversion int = 5 // current version in the draft

	for _, record := range message.Records {
		// get base attributes from current record
		if record.BaseName != nil && len(*record.BaseName) > 0 {
			basename = *record.BaseName
		}
		if record.BaseTime != nil && *record.BaseTime > 0 {
			basetime = *record.BaseTime
		}
		if record.BaseUnit != nil && len(*record.BaseUnit) > 0 {
			baseunit = *record.BaseUnit
		}
		if record.BaseValue != nil && *record.BaseValue > 0 {
			basevalue = *record.BaseValue
		}
		if record.BaseSum != nil && *record.BaseSum > 0 {
			basesum = *record.BaseSum
		}
		if record.Version != nil && *record.Version > baseversion {
			return nil, errors.New("version number is higher than supported")
		}

		// delete base attributes from record
		record.BaseName = nil
		record.BaseTime = nil
		record.BaseUnit = nil
		record.BaseValue = nil
		record.BaseSum = nil

		// 1. prepend the basename
		combinedName := basename
		if record.Name != nil {
			combinedName += *record.Name
		}
		record.Name = &combinedName

		// 2. add base time to every time field and convert time to absolute
		/*
						quote from the draft:
						A time of zero indicates that the sensor does not know the absolute
						time and the measurement was made roughly "now".  A negative value is
			   			used to indicate seconds in the past from roughly "now".  A positive
			   			value is used to indicate the number of seconds, excluding leap
			   			seconds, since the start of the year 1970 in UTC.
		*/
		combinedTime := basetime
		if record.Time != nil {
			combinedTime += *record.Time
		}
		record.Time = &combinedTime
		if *record.Time <= 0 {
			var now int64 = time.Now().UnixNano() / 1000000000.0
			absoluteTime := float64(now) + *record.Time
			record.Time = &absoluteTime
		}

		// 3. populate base unit on empty unit fields
		if record.Unit == nil || len(*record.Unit) <= 0 {
			currentBaseUnit := baseunit
			record.Unit = &currentBaseUnit
		}

		// 4. add base value to every value field
		combinedValue := basevalue
		if record.Value != nil {
			combinedValue += *record.Value
		}
		record.Value = &combinedValue

		// 5. add base sum to every sum field
		combinedSum := basesum
		if record.Sum != nil {
			combinedSum += *record.Sum
		}
		record.Sum = &combinedSum

		// 6. set version to baseversion
		record.Version = &baseversion

		// add the record to the output message if a value is set
		if record.Value != nil || record.StringValue != nil || record.BoolValue != nil || record.DataValue != nil || record.Sum != nil {
			resolvedMessage.Records = append(resolvedMessage.Records, record)
		} else {
			log.Warning("record has no value and won't be appended")
		}
	}

	return &resolvedMessage, nil
}
