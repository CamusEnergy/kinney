package schema

import (
	"encoding/xml"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func newInt32(x int32) *int32 { return &x }

func newString(x string) *string { return &x }

func TestMarshal(t *testing.T) {
	envelopeRE := regexp.MustCompile(`\s*<`)

	for name, tt := range map[string]struct {
		Value    interface{}
		Expected string
	}{
		"ClearShedStateRequest_StationGroupID": {
			&ClearShedStateRequest{
				StationGroupID: newInt32(1234),
			},
			`<clearShedState xmlns="urn:dictionary:com.chargepoint.webservices">
			  <shedQuery>
			    <shedGroup>
			      <sgID>1234</sgID>
			    </shedGroup>
			    <shedStation>
			      <Ports>
				<Port></Port>
			      </Ports>
			    </shedStation>
			  </shedQuery>
			</clearShedState>`,
		},
		"ClearShedStateRequest_StationID": {
			&ClearShedStateRequest{
				StationID: newString("station id"),
			},
			`<clearShedState xmlns="urn:dictionary:com.chargepoint.webservices">
			  <shedQuery>
			    <shedStation>
			      <stationID>station id</stationID>
			      <Ports>
				<Port></Port>
			      </Ports>
			    </shedStation>
			  </shedQuery>
			</clearShedState>`,
		},
		"ClearShedStateRequest_PortNumbers": {
			&ClearShedStateRequest{
				StationID:   newString("station id"),
				PortNumbers: []string{"0", "1"},
			},
			`<clearShedState xmlns="urn:dictionary:com.chargepoint.webservices">
			  <shedQuery>
			    <shedStation>
			      <stationID>station id</stationID>
			      <Ports>
				<Port>
				  <portNumber>0</portNumber>
				  <portNumber>1</portNumber>
				</Port>
			      </Ports>
			    </shedStation>
			  </shedQuery>
			</clearShedState>`,
		},

		"GetCPNInstancesResponse": {
			&GetCPNInstancesResponse{
				ChargePointNetworks: []GetCPNInstancesResponse_ChargePointNetwork{
					{ID: "network id", Name: "network name", Description: "description"},
				},
			},
			`<getCPNInstancesResponse xmlns="urn:dictionary:com.chargepoint.webservices">
			  <CPN>
			    <cpnID>network id</cpnID>
			    <cpnName>network name</cpnName>
			    <cpnDescription>description</cpnDescription>
			  </CPN>
			</getCPNInstancesResponse>`,
		},

		"GetLoadResponse": {
			&GetLoadResponse{
				StationGroupID:          int32(1234),
				StationGroupNumStations: int32(1),
				StationGroupName:        "station group name",
				StationGroupLoadKW:      "0.5",
				Stations: []GetLoadResponse_Station{
					{
						StationID:      "station id",
						StationName:    "station name",
						StationAddress: "station address",
						StationLoadKW:  "0.5",
						Ports: []GetLoadResponse_Station_Port{
							{PortNumber: "0", UserID: "user id"},
						},
					},
				},
			},
			`<getLoadResponse xmlns="urn:dictionary:com.chargepoint.webservices">
			  <responseCode></responseCode>
			  <sgID>1234</sgID>
			  <numStations>1</numStations>
			  <groupName>station group name</groupName>
			  <sgLoad>0.5</sgLoad>
			  <stationData>
			    <stationID>station id</stationID>
			    <stationName>station name</stationName>
			    <Address>station address</Address>
			    <stationLoad>0.5</stationLoad>
			    <Port>
			      <portNumber>0</portNumber>
			      <userID>user id</userID>
			    </Port>
			  </stationData>
			</getLoadResponse>`,
		},

		"GetOrgsAndStationGroupsRequest": {
			&GetOrgsAndStationGroupsRequest{
				OrganizationID: "org id",
			},
			`<getOrgsAndStationGroups xmlns="urn:dictionary:com.chargepoint.webservices">
			  <searchQuery>
			    <orgID>org id</orgID>
			  </searchQuery>
			</getOrgsAndStationGroups>`,
		},

		"GetOrgsAndStationGroupsResponse": {
			&GetOrgsAndStationGroupsResponse{
				Organizations: []GetOrgsAndStationGroupsResponse_Organization{
					{
						OrganizationID:   "org id",
						OrganizationName: "org name",
						StationGroups: []GetOrgsAndStationGroupsResponse_Organization_StationGroup{
							{StationGroupID: 1234, StationGroupName: "sg name"},
						},
					},
				},
			},
			`<getOrgsAndStationGroupsResponse xmlns="urn:dictionary:com.chargepoint.webservices">
			  <responseCode></responseCode>
			  <orgData>
			    <orgID>org id</orgID>
			    <organizationName>org name</organizationName>
			    <sgData>
			      <sgID>1234</sgID>
			      <sgName>sg name</sgName>
			    </sgData>
			  </orgData>
			</getOrgsAndStationGroupsResponse>`,
		},

		"GetStationGroupsResponse": {
			&GetStationGroupsResponse{
				StationGroups: []GetStationGroupsResponse_StationGroup{
					{
						OrganizationID: "organization id",
						Stations: []GetStationGroupsResponse_StationGroup_Station{
							{StationID: "station id"},
						},
					},
				},
			},
			`<getStationGroupsResponse xmlns="urn:dictionary:com.chargepoint.webservices">
			  <responseCode></responseCode>
			  <groupData>
			    <orgID>organization id</orgID>
			    <stationData>
			      <stationID>station id</stationID>
			    </stationData>
			  </groupData>
			</getStationGroupsResponse>`,
		},

		"GetStationsRequest": {
			&GetStationsRequest{
				StationID: "station id",
				PricingSession: &GetStationsRequest_PricingSession{
					StartTime: "some time",
				},
				DemoStationSerialNumbers: []string{"serial number"},
			},
			`<getStations xmlns="urn:dictionary:com.chargepoint.webservices">
			  <searchQuery>
			    <stationID>station id</stationID>
			    <Pricing>
			      <startTime>some time</startTime>
			      <Duration>0</Duration>
			      <energyRequired>0</energyRequired>
			      <vehiclePower>0</vehiclePower>
			    </Pricing>
			    <demoSerialNumber>
			      <serialNumber>serial number</serialNumber>
			    </demoSerialNumber>
			  </searchQuery>
			</getStations>`,
		},

		"GetStationsResponse": {
			&GetStationsResponse{
				Stations: []GetStationsResponse_Station{
					{
						StationID: "station id",
						Ports: []GetStationsResponse_Station_Port{
							{PortNumber: "port number"},
						},
						PricingSpecification: []GetStationsResponse_Station_PricingSpecification{
							{Type: "Session"},
						},
					},
				},
			},
			`<getStationsResponse xmlns="urn:dictionary:com.chargepoint.webservices">
			  <responseCode></responseCode>
			  <stationData>
			    <stationID>station id</stationID>
			    <orgID></orgID>
			    <organizationName></organizationName>
			    <Port>
			      <portNumber>port number</portNumber>
			    </Port>
			    <Pricing>
			      <Type>Session</Type>
			    </Pricing>
			  </stationData>
			</getStationsResponse>`,
		},

		"GetUsersRequest": {
			&GetUsersRequest{
				UserID: "user_id",
				Connection: &GetUsersRequest_Connection{
					Status: "APPROVED",
					CustomInfo: &GetUsersRequest_CustomInfo{
						Key:   "foo",
						Value: "bar",
					},
				},
			},
			`<getUsers xmlns="urn:dictionary:com.chargepoint.webservices">
			  <searchQuery>
			    <userID>user_id</userID>
			    <Connection>
			      <Status>APPROVED</Status>
			      <customInfo>
			        <Key>foo</Key>
				<Value>bar</Value>
			      </customInfo>
			    </Connection>
			  </searchQuery>
			</getUsers>`,
		},

		"GetUsersResponse": {
			&GetUsersResponse{
				Users: []GetUsersResponse_User{
					{
						UserID: "user_id",
						Connection: &GetUsersResponse_Connection{
							Status: "APPROVED",
							CustomInfo: []GetUsersResponse_CustomInfo{
								{Key: "foo", Value: "bar"},
							},
						},
					},
				},
			},
			`<getUsersResponse xmlns="urn:dictionary:com.chargepoint.webservices">
			  <responseCode></responseCode>
			  <users>
			    <user>
			      <userID>user_id</userID>
			      <Connection>
			        <Status>APPROVED</Status>
				<customInfos>
				  <customInfo>
				    <Key>foo</Key>
				    <Value>bar</Value>
				  </customInfo>
				</customInfos>
			      </Connection>
			      <credentialIDs></credentialIDs>
			    </user>
			  </users>
			</getUsersResponse>`,
		},

		"ShedLoadRequest_StationGroupID": {
			&ShedLoadRequest{
				StationGroupID:            1234,
				StationGroupAllowedLoadKW: "0.5",
			},
			`<shedLoad xmlns="urn:dictionary:com.chargepoint.webservices">
			  <shedQuery>
			    <shedGroup>
			      <sgID>1234</sgID>
			      <allowedLoadPerStation>0.5</allowedLoadPerStation>
			    </shedGroup>
			    <timeInterval>0</timeInterval>
			  </shedQuery>
			</shedLoad>`,
		},
		"ShedLoadRequest_StationID": {
			&ShedLoadRequest{
				Station: &ShedLoadRequest_Station{
					StationID:     "station id",
					AllowedLoadKW: "0.5",
				},
			},
			`<shedLoad xmlns="urn:dictionary:com.chargepoint.webservices">
			  <shedQuery>
			    <shedGroup></shedGroup>
			    <shedStation>
			      <stationID>station id</stationID>
			      <allowedLoadPerStation>0.5</allowedLoadPerStation>
			    </shedStation>
			    <timeInterval>0</timeInterval>
			  </shedQuery>
			</shedLoad>`,
		},
		"ShedLoadRequest_Ports": {
			&ShedLoadRequest{
				Station: &ShedLoadRequest_Station{
					StationID: "station id",
					Ports: &ShedLoadRequest_Station_Ports{
						Ports: []ShedLoadRequest_Station_Ports_Port{
							{PortNumber: "port number", AllowedLoadKW: "0.5"},
						},
					},
				},
			},
			`<shedLoad xmlns="urn:dictionary:com.chargepoint.webservices">
			  <shedQuery>
			    <shedGroup></shedGroup>
			    <shedStation>
			      <stationID>station id</stationID>
			      <Ports>
			        <Port>
			          <portNumber>port number</portNumber>
			          <allowedLoadPerPort>0.5</allowedLoadPerPort>
			        </Port>
			      </Ports>
			    </shedStation>
			    <timeInterval>0</timeInterval>
			  </shedQuery>
			</shedLoad>`,
		},

		"ShedLoadResponse": {
			&ShedLoadResponse{
				Success:        1,
				StationGroupID: 1234,
				Ports: []ShedLoadResponse_Port{
					{PortNumber: "port number", AllowedLoadKW: "0.5"},
				},
			},
			`<shedLoadResponse xmlns="urn:dictionary:com.chargepoint.webservices">
			  <responseCode></responseCode>
			  <Success>1</Success>
			  <sgID>1234</sgID>
			  <Ports>
			    <Port>
			      <portNumber>port number</portNumber>
			      <allowedLoadPerPort>0.5</allowedLoadPerPort>
			    </Port>
			  </Ports>
			</shedLoadResponse>`,
		},
	} {
		expected := envelopeRE.ReplaceAllString(tt.Expected, "<")

		if b, err := xml.Marshal(tt.Value); err != nil {
			t.Errorf("%s: xml.Marshal() = %q; want nil", name, err)
		} else if diff := cmp.Diff(expected, string(b)); diff != "" {
			t.Errorf("%s: marshaling mismatch (-want +got):\n%s", name, diff)
		}
	}
}
