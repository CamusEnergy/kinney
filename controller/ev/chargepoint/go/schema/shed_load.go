package schema

import (
	"encoding/xml"
	"math/big"
)

// API Guide (§ 6.1): "Use this call to shed load for a single port on a
// station, both ports on a multi-port station or a group of stations.  Only one
// of these three options may be used in a request as follows:
// - Group: Include the shedGroup element.
// - Station: Include the shedStation element and either the
//   allowedLoadPerStation or percentShedPerStation parameters within that
//   element; omit the Ports array.
// - Port: Include the shedStation element and the Ports array; set the
//   allowedLoadPerStation and percentShedPerStation parameters in the
//   shedStation element to a null value or omit them from the request."
type ShedLoadRequestParams struct {
	ShedQuery struct {
		ShedStationGroup *struct {
			StationGroupID int32 `xml:"sgID"`

			// Only one of the following two fields may be set.

			// API Guide (§ 6.1.2): "Maximum allowed load expressed
			// in kW.  This value is an absolute maximum and is not
			// relative to the power being dispensed by the station.
			// At the group level, this parameter applies to each
			// station, not the total power for the group."
			AllowedLoadKW *big.Rat `xml:"allowedLoadPerStation,omitempty"`

			// API Guide (§ 6.1.2): "Percentage of the power
			// currently being dispensed by the station to shed.r
			// For example, if the station is currently dispensing
			// 10kW, a value of 60% will lower the power being
			// dispensed to 4kW.  At the group level, this value
			// applies to each station.  If a station is not
			// dispensing any power, the output will be set to zero
			// until the shed state is cleared."
			PercentShed *int32 `xml:"percentShedPerStation,omitempty"`
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

		// API Guide (§ 6.1.2): "Time interval in minutes.  A value of 0
		// indicates that there is no specified duration for which the
		// power will be shed."
		TimeInterval int32 `xml:"timeInterval"`
	} `xml:"shedQuery"`
}

type ShedLoadRequest struct {
	XMLName xml.Name `xml:"urn:dictionary:com.chargepoint.webservices shedLoad"`
	ShedLoadRequestParams
}

type ShedLoadResponse struct {
	XMLName xml.Name `xml:"shedLoadResponse"`

	commonResponseParameters

	// API Guide (§ 6.1.3): "A success (1) or failure (0) response code
	// only."
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
