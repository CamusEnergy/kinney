package schema

import (
	"encoding/xml"
	"math/big"
)

// API Guide: "Use this call to shed load for a single port on a station, both
// ports on a multi-port station or a group of stations.  Only one of these
// three options may be used in a request as follows:
// - Group: Include the shedGroup element.
// - Station: Include the shedStation element and either the
//   allowedLoadPerStation or percentShedPerStation parameters within that
//   element; omit the Ports array.
// - Port: Include the shedStation element and the Ports array; set the
//   allowedLoadPerStation and percentShedPerStation parameters in the
//   shedStation element to a null value or omit them from the request."
type ShedLoadRequest struct {
	XMLName xml.Name `xml:"urn:dictionary:com.chargepoint.webservices shedLoad"`

	ShedQuery struct {
		ShedStationGroup *struct {
			StationGroupID int32 `xml:"sgID"`

			// These parameters apply *per station* in the station
			// group; not to the group in aggregate.
			AllowedLoadKW *big.Rat `xml:"allowedLoadPerStation,omitempty"`
			PercentShed   *int32   `xml:"percentShedPerStation,omitempty"`
		} `xml:"shedGroup,omitempty"`

		ShedStation *struct {
			StationID     string   `xml:"stationID"`
			AllowedLoadKW *big.Rat `xml:"allowedLoadPerStation,omitempty"`
			PercentShed   *int32   `xml:"percentShedPerStation,omitempty"`

			Ports []struct {
				PortNumber    string   `xml:"portNumber"`
				AllowedLoadKW *big.Rat `xml:"allowedLoadPerPort"`
				PercentShed   *int32   `xml:"percentShedPerPort"`
			} `xml:"Ports>Port,omitempty"`
		} `xml:"shedStation,omitempty"`

		// API Guide: "Time interval in minutes.  A value of 0 indicates
		// that there is no specified duration for which the power will
		// be shed.
		TimeInterval int32 `xml:"timeInterval"`
	} `xml:"shedQuery"`
}

type ShedLoadResponse struct {
	XMLName xml.Name `xml:"shedLoadResponse"`

	commonResponseParameters

	// API Guide: "A success (1) or failure (0) response code only.
	//
	// This field has `type="xsd:int"` in the WSDL, but is semantically
	// boolean, and can safely be parsed as such.
	Success bool `xml:"Success,omitempty"`

	StationGroupID int32  `xml:"sgID,omitempty"`
	StationID      string `xml:"stationID,omitempty"`

	AllowedLoadKW *big.Rat `xml:"allowedLoadPerStation,omitempty"`
	PercentShed   *int32   `xml:"percentShedPerStation,omitempty"`

	Ports []struct {
		PortNumber    string   `xml:"portNumber"`
		AllowedLoadKW *big.Rat `xml:"allowedLoadPerPort,omitempty"`
		PercentShed   *int32   `xml:"percentShedPerPort,omitempty"`
	} `xml:"Ports>Port,omitempty"`
}
