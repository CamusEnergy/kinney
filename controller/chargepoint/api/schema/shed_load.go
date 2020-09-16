package schema

import (
	"encoding/xml"
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
type ShedLoadRequest struct {
	XMLName xml.Name `xml:"urn:dictionary:com.chargepoint.webservices shedLoad"`

	StationGroupID int32 `xml:"shedQuery>shedGroup>sgID,omitempty"`
	// Only one of the following two fields may be set.
	//
	// API Guide (§ 6.1.2): "Maximum allowed load expressed in kW.  This
	// value is an absolute maximum and is not relative to the power being
	// dispensed by the station.  At the group level, this parameter applies
	// to each station, not the total power for the group."
	StationGroupAllowedLoadKW string `xml:"shedQuery>shedGroup>allowedLoadPerStation,omitempty"`
	// API Guide (§ 6.1.2): "Percentage of the power currently being
	// dispensed by the station to shed.  For example, if the station is
	// currently dispensing 10kW, a value of 60% will lower the power being
	// dispensed to 4kW.  At the group level, this value applies to each
	// station.  If a station is not dispensing any power, the output will
	// be set to zero until the shed state is cleared."
	StationGroupPercentShed *int32 `xml:"shedQuery>shedGroup>percentShedPerStation,omitempty"`

	Station *ShedLoadRequest_Station `xml:"shedQuery>shedStation,omitempty"`

	// API Guide (§ 6.1.2): "Time interval in minutes.  A value of 0
	// indicates that there is no specified duration for which the power
	// will be shed."
	TimeInterval int32 `xml:"shedQuery>timeInterval"`
}

type ShedLoadRequest_Station struct {
	StationID string `xml:"stationID"`
	// Only one of the following two fields may be set.
	AllowedLoadKW string `xml:"allowedLoadPerStation,omitempty"`
	PercentShed   *int32 `xml:"percentShedPerStation,omitempty"`

	Ports *ShedLoadRequest_Station_Ports `xml:"Ports,omitempty"`
}

type ShedLoadRequest_Station_Ports struct {
	Ports []ShedLoadRequest_Station_Ports_Port `xml:"Port"`
}

type ShedLoadRequest_Station_Ports_Port struct {
	PortNumber    string `xml:"portNumber"`
	AllowedLoadKW string `xml:"allowedLoadPerPort"`
	PercentShed   *int32 `xml:"percentShedPerPort"`
}

type ShedLoadResponse struct {
	XMLName xml.Name `xml:"urn:dictionary:com.chargepoint.webservices shedLoadResponse"`

	commonResponseParameters

	// API Guide (§ 6.1.3): "A success (1) or failure (0) response code
	// only."
	Success uint8 `xml:"Success,omitempty"`

	StationGroupID int32  `xml:"sgID,omitempty"`
	StationID      string `xml:"stationID,omitempty"`

	AllowedLoadKW string `xml:"allowedLoadPerStation,omitempty"`
	PercentShed   *int32 `xml:"percentShedPerStation,omitempty"`

	Ports []ShedLoadResponse_Port `xml:"Ports>Port,omitempty"`
}

type ShedLoadResponse_Port struct {
	PortNumber    string `xml:"portNumber"`
	AllowedLoadKW string `xml:"allowedLoadPerPort,omitempty"`
	PercentShed   *int32 `xml:"percentShedPerPort,omitempty"`
}
