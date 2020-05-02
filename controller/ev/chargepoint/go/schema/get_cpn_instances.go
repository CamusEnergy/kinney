package schema

import "encoding/xml"

type GetCPNInstancesRequest struct {
	XMLName xml.Name `xml:"urn:dictionary:com.chargepoint.webservices getCPNInstances"`
}

type GetCPNInstancesResponse struct {
	XMLName xml.Name `xml:"getCPNInstancesResponse"`

	ChargePointNetworks []struct {
		ChargePointNetworkID          string `xml:"cpnID,omitempty"`
		ChargePointNetworkName        string `xml:"cpnName,omitempty"`
		ChargePointNetworkDescription string `xml:"cpnDescription,omitempty"`
	} `xml:"CPN,omitempty"`
}
