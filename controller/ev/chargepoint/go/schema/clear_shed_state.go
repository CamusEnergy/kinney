package schema

import "encoding/xml"

type ClearShedStateRequest struct {
	XMLName xml.Name `xml:"urn:dictionary:com.chargepoint.webservices clearShedState"`

	ShedQuery struct {
		StationGroup *struct {
			StationGroupID int32 `xml:"sgID"`
		} `xml:"shedGroup,omitempty"`

		Station *struct {
			StationID string `xml:"stationID"`

			Ports *struct {
				PortNumbers []string `xml:"Port>portNumber"`
			} `xml:"Ports,omitempty"`
		} `xml:"shedStation,omitempty"`
	} `xml:"shedQuery"`
}

type ClearShedStateResponse struct {
	XMLName xml.Name `xml:"clearShedStateResponse"`

	commonResponseParameters

	// API Guide: "A success (1) or failure (0) response code only.
	//
	// This field has `type="xsd:int"` in the WSDL, but is semantically
	// boolean, and can safely be parsed as such.
	Success bool `xml:"Success,omitempty"`

	StationGroupID int32  `xml:"sgID,omitempty"`
	StationID      string `xml:"stationID,omitempty"`
}
