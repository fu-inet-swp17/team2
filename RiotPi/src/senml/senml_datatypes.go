/*
	author: Niclas Kristek
	github.com/nkristek
*/
package senml

type SenMLRecord struct {
	/*
		This is a string that is prepended to the names found in the entries.
	*/
	BaseName *string `json:"bn,omitempty"`

	/*
		A base time that is added to the time found in an entry.
	*/
	BaseTime *float64 `json:"bt,omitempty"`

	/*
			A base unit that is assumed for all entries, unless
		    otherwise indicated.  If a record does not contain a Unit value,
		    then the Base Unit is used.  Otherwise the value found in the Unit
		    (if any) is used.
	*/
	BaseUnit *string `json:"bu,omitempty"`

	/*
		A base value is added to the value found in an entry, similar to Base Time.
	*/
	BaseValue *float64 `json:"bv:omitempty"`

	/*
		A base sum is added to the sum found in an entry, similar to Base Time.
	*/
	BaseSum *float64 `json:"bs:omitempty"`

	/*
			Version number of media type format.  This attribute is an
		    optional positive integer and defaults to 5 if not present.  [RFC
		    Editor: change the default value to 10 when this specification is
		    published as an RFC and remove this note]
	*/
	Version *int `json:"bver,omitempty"`

	/*
			Name of the sensor or parameter.
			When appended to the Base Name attribute, this must result in a
			globally unique identifier for the resource.  The name is optional,
			if the Base Name is present.  If the name is missing, Base Name must
			uniquely identify the resource.  This can be used to represent a large
			array of measurements from the same sensor without having to repeat its
		    identifier on every measurement.
	*/
	Name *string `json:"n,omitempty"`

	/*
		Units for a measurement value.  Optional.
	*/
	Unit *string `json:"u,omitempty"`

	/*
			Value of the entry.  Optional if a Sum value is present,
		    otherwise required.  Values are represented using basic data
		    types.  This specification defines floating point numbers ("v"
		    field for "Value"), booleans ("vb" for "Boolean Value"), strings
		    ("vs" for "String Value") and binary data ("vd" for "Data Value").
		    Exactly one value field MUST appear unless there is Sum field in
		    which case it is allowed to have no Value field.
	*/
	Value       *float64 `json:"v,omitempty"`
	StringValue *string  `json:"vs,omitempty"`
	BoolValue   *bool    `json:"vb,omitempty"`
	DataValue   *string  `json:"vd,omitempty"`

	/*
			Integrated sum of the values over time.  Optional.  This
		    attribute is in the units specified in the Unit value multiplied
		    by seconds.
	*/
	Sum *float64 `json:"s,omitempty"`

	/*
		Time when value was recorded.  Optional.
	*/
	Time *float64 `json:"t,omitempty"`

	/*
	 	An optional time in seconds that represents the maximum
	    time before this sensor will provide an updated reading for a
	    measurement.  This can be used to detect the failure of sensors or
	    communications path from the sensor.
	*/
	UpdateTime *float64 `json:"ut,omitempty"`

	/*

	 */
	Link *string `json:"l,omitempty"`
}

type SenMLMessage struct {
	Records []SenMLRecord
}
