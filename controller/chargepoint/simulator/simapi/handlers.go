package simapi

import (
	"fmt"
	"strconv"
	
	"github.com/CamusEnergy/kinney/controller/chargepoint/api/schema"
	"github.com/CamusEnergy/kinney/controller/chargepoint/simulator/sim"
)

type SimulatorServer struct {
	Ev *sim.EvChargers
}

func (s SimulatorServer) GetLoad (req *schema.GetLoadRequest) (*schema.GetLoadResponse, error) {
	resp := &schema.GetLoadResponse{}

	resp.ResponseCode = "188";
	resp.ResponseText = "Not implemented yet "
	
	fmt.Println(req)
	
	return resp, nil
}

func (s SimulatorServer) GetStations (req *schema.GetStationsRequest) (*schema.GetStationsResponse, error) {
	resp := &schema.GetStationsResponse{}
	
	for j,org := range s.Ev.StOrgs {
		if req.OrganizationID != "" && req.OrganizationID != j {
			continue
		}
		if req.OrganizationName != "" && req.OrganizationName != org.Name {
			continue
		}
	
		for i,g := range org.StGrpoups {
			if req.StationGroupID != "" && req.StationGroupID != i {
				continue
			}
			if req.StationGroupName != "" && req.StationGroupName != g.Name {
				continue
			}

			for k,st := range g.Stations {
				if req.StationID != "" && req.StationID != k {
					continue
				}

				station := schema.GetStationsResponse_Station {
								OrganizationID: j,
								OrganizationName: org.Name,
								StationGroupID: i,
								StationID: k,
								NumPorts: int32(len(st.Ports)),
						}
				for l, _ := range st.Ports {
					port := schema.GetStationsResponse_Station_Port {
						PortNumber: l,
						Coordinate: &schema.Coordinate {
								Latitude: st.Geo.Lat,
								Longitude: st.Geo.Long,
						},						
					}
					station.Ports = append(station.Ports, port)
				}
				resp.Stations = append(resp.Stations, station)
			}
		}
	}
	
	if len(resp.Stations) > 0 {
		resp.ResponseCode = "100";
		resp.ResponseText = "OK"
		
	} else {
		resp.ResponseCode = "102";
		resp.ResponseText = "No stations found"
	}
	
	return resp, nil	
}

func (s SimulatorServer) GetStationGroups (req *schema.GetStationGroupsRequest) (*schema.GetStationGroupsResponse, error) {
	resp := &schema.GetStationGroupsResponse{}
	
	for j,org := range s.Ev.StOrgs {
		if req.OrganizationID != "" && req.OrganizationID != j {
			continue
		}
		for i,g := range org.StGrpoups {
			gid,err := strconv.ParseInt(i, 10, 32)
			if err != nil {
				continue
			}
			r := schema.GetStationGroupsResponse_StationGroup{
							OrganizationID: j,
							OrganizationName: org.Name,
							StationGroupID: int32(gid),
							StationGroupName: g.Name,
							}
			for k,st := range g.Stations {
				station := schema.GetStationGroupsResponse_StationGroup_Station {
								StationID: k,
								Coordinate: &schema.Coordinate {
										Latitude: st.Geo.Lat,
										Longitude: st.Geo.Long,
								},
						}
				r.Stations = append(r.Stations, station)
			}
			resp.StationGroups = append(resp.StationGroups, r)
		}
	}
	
	if len(resp.StationGroups) > 0 {
		resp.ResponseCode = "100";
		resp.ResponseText = "OK"
		
	} else {
		resp.ResponseCode = "102";
		resp.ResponseText = "No station groups found"
	}
	
	return resp, nil
}

func (s SimulatorServer) GetCPNInstances (req *schema.GetCPNInstancesRequest) (*schema.GetCPNInstancesResponse, error) {
	resp := &schema.GetCPNInstancesResponse{}
	
	for i,n := range s.Ev.StCPNs {
		resp.ChargePointNetworks = append(resp.ChargePointNetworks,
						schema.GetCPNInstancesResponse_ChargePointNetwork{
							ID:i,
							Name:n.Name,
							Description:n.Desc,})
	}

	return resp, nil	
}

func (s SimulatorServer) ShedLoad (req *schema.ShedLoadRequest) (*schema.ShedLoadResponse, error) {
	resp := &schema.ShedLoadResponse{}
	
	resp.ResponseCode = "188";
	resp.ResponseText = "Not implemented yet "
	
	fmt.Println(req)
		
	return resp, nil	
}

func (s SimulatorServer) ClearShedState (req *schema.ClearShedStateRequest) (*schema.ClearShedStateResponse, error) {
	resp := &schema.ClearShedStateResponse{}
	
	for _,org := range s.Ev.StOrgs {
		for i,g := range org.StGrpoups {
			id,err := strconv.ParseInt(i, 10, 32)
			if err != nil {
				continue
			}			
			if req.StationGroupID != nil && *req.StationGroupID != int32(id) {
				continue
			}

			for k,st := range g.Stations {
				if req.StationID != nil && *req.StationID != k {
					continue
				}
				for _, pr := range st.Ports {
					pr.Sched = 0
				}
			}
		}
	}		
	return resp, nil	
}
