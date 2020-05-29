package schema

// XML Schema 1.0 "time" built-in datatype.
//
// https://www.w3.org/TR/xmlschema-2/#time
//
// TODO(james): Add validation of the value.
type xsdTime string

// XML Schema 1.0 "dateTime" built-in datatype.
//
// https://www.w3.org/TR/xmlschema-2/#dateTime
//
// TODO(james): Add validation of the value.
type xsdDateTime string

// API Guide (§ 2.5): "All responses of the ChargePoint Web Services API contain
// the following parameters."
type commonResponseParameters struct {
	// API Guide (§ 2.5): "Code indicating success or failure for the API
	// call."
	//
	// Everything but "100" represents an error response.
	ResponseCode string `xml:"responseCode"`

	// API Guide (§ 2.5): "If an error occurs, this field contains a
	// description of the error.  If no error occurred, this field will be
	// blank."
	ResponseText string `xml:"responseText,omitempty"`
}

// TODO(james): Figure out the CRS that the api uses for these coordinates.
type Coordinate struct {
	Latitude  string `xml:"Lat"`
	Longitude string `xml:"Long"`
}
