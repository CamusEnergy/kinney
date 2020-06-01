package schema

import (
	"encoding/xml"
	"math/big"
)

// API Guide (§ 6.3): "Use this call to retrieve the load and shed state for a
// single station or custom station group.  This method also returns the load
// for each port on a multi-port station."
type GetLoadRequestParams struct {
	StationGroupID int32  `xml:"searchQuery>sgID,omitempty"`
	StationID      string `xml:"searchQuery>stationID,omitempty"`
}

type GetLoadRequest struct {
	XMLName xml.Name `xml:"urn:dictionary:com.chargepoint.webservices getLoad"`
	GetLoadRequestParams
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
			// API Guide (§ 6.3.3): "Identifier of the credential
			// used to start the session.  If it was a ChargePoint
			// RFID card, it is the printed serial number on the
			// card.  If it was the ChargePoint Mobile App, it will
			// be the identifier displayed in the user’s mobile app.
			// Contactless credit cards will be obviously be
			// displayed as blank."
			CredentialID *string `xml:"credentialID,omitempty"`

			PortLoadKW big.Rat `xml:"portLoad,omitempty"`

			// API Guide (§ 6.3.3): "1 = Shed, 0 = Not Shed"
			ShedState uint8 `xml:"shedState"`

			// Only one of the following two fields should be set.

			// API Guide (§ 6.3.3): "Maximum load allowed at the
			// station (kW).  If percentShed was used in the last
			// shedLoad call to this station, this parameter will be
			// zero."
			AllowedLoadKW *big.Rat `xml:"allowedLoad,omitempty"`
			// API Guide (§ 6.3.3): "Percent of load currently being
			// shed.  If allowedLoad was used in the last shedLoad
			// call to this station, this parameter will be zero."
			//
			// Percent of load currently being shed (0 - 100).
			PercentShed *uint8 `xml:"percentShed,omitempty"`
		} `xml:"Port,omitempty"`
	} `xml:"stationData,omitempty"`
}
