package schema

import (
	"encoding/xml"
	"math/big"
)

type GetLoadRequest struct {
	XMLName xml.Name `xml:"urn:dictionary:com.chargepoint.webservices getLoad"`

	StationGroupID int32  `xml:"searchQuery>sgID,omitempty"`
	StationID      string `xml:"searchQuery>stationID,omitempty"`
}

type GetLoadResponse struct {
	XMLName xml.Name `xml:"getLoadResponse"`

	commonResponseParameters

	StationGroupID          int32  `xml:"sgID,omitempty"`
	StationGroupNumStations int32  `xml:"numStations,omitempty"`
	StationGroupName        string `xml:"groupName,omitempty"`
	StationGroupLoadKW      string `xml:"sgLoad,omitempty"`

	Stations []struct {
		StationID      string  `xml:"stationID,omitempty"`
		StationName    string  `xml:"stationName,omitempty"`
		StationAddress string  `xml:"Address,omitempty"`
		StationLoadKW  big.Rat `xml:"stationLoad,omitempty"`

		Ports []struct {
			PortNumber string `xml:"portNumber,omitempty"`
			UserID     string `xml:"userID,omitempty"`
			// The empty string means that a contactless credit card
			// was used to start the session.
			CredentialID *string `xml:"credentialID,omitempty"`

			PortLoadKW big.Rat `xml:"portLoad,omitempty"`

			// "1 = Shed, 0 = Not Shed"
			ShedState uint8 `xml:"shedState"`

			// Only one of the following two should be set.
			// Maximum load allowed at the station.
			AllowedLoadKW *big.Rat `xml:"allowedLoad,omitempty"`
			// Percent of load currently being shed (0 - 100).
			PercentShed *uint8 `xml:"percentShed,omitempty"`
		} `xml:"Port,omitempty"`
	} `xml:"stationData,omitempty"`
}
