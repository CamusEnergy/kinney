// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package sim

import (
	"strconv"

	"github.com/CamusEnergy/kinney/controller/chargepoint/api/schema"
)

type SimulatorServer struct {
	Ev *EVChargers
}

func (s SimulatorServer) GetLoad(req *schema.GetLoadRequest) (*schema.GetLoadResponse, error) {
	resp := &schema.GetLoadResponse{}
	var group *chargeGroup
	var err error
	sgid := ""

	if req.StationGroupID != 0 {
		sgid = strconv.Itoa(int(req.StationGroupID))
	}

	for _, f := range s.Ev.facilities {
		if sgid != "" {
			if g, ok := f.groups[sgid]; !ok {
				continue
			} else {
				group = g
			}
		} else {
			for i, g := range f.groups {
				if id, err := strconv.Atoi(i); err == nil {
					req.StationGroupID = int32(id)
					group = g
				}
				break
			}
		}
		break
	}

	err = s.Ev.getNextLoad(resp, group, &req.StationID)
	if err != nil {
		resp.ResponseCode = "102"
		resp.ResponseText = err.Error()
	} else {
		resp.ResponseCode = "100"
		resp.ResponseText = "OK"
		resp.StationGroupNumStations = int32(len(group.stations))
		resp.StationGroupID = req.StationGroupID
		if group.name != nil {
			resp.StationGroupName = *group.name
		}
	}

	return resp, nil
}

func (s SimulatorServer) GetStations(req *schema.GetStationsRequest) (*schema.GetStationsResponse, error) {
	resp := &schema.GetStationsResponse{}

	for j, org := range s.Ev.facilities {
		if req.OrganizationID != "" && req.OrganizationID != j {
			continue
		}
		if req.OrganizationName != "" && req.OrganizationName != *org.name {
			continue
		}

		for i, g := range org.groups {
			if req.StationGroupID != "" && req.StationGroupID != i {
				continue
			}
			if req.StationGroupName != "" && req.StationGroupName != *g.name {
				continue
			}

			for k, st := range g.stations {
				if req.StationID != "" && req.StationID != k {
					continue
				}

				station := schema.GetStationsResponse_Station{
					OrganizationID: j,
					StationGroupID: i,
					StationID:      k,
					NumPorts:       int32(len(st.ports)),
				}
				if org.name != nil {
					station.OrganizationName = *org.name
				}
				if st.address != nil {
					station.Address = *st.address
				}
				for l := range st.ports {
					port := schema.GetStationsResponse_Station_Port{
						PortNumber: l,
						Coordinate: &schema.Coordinate{
							Latitude:  st.geo.lat,
							Longitude: st.geo.long,
						},
					}
					if st.name != nil {
						port.StationName = *st.name
					}
					station.Ports = append(station.Ports, port)
				}
				resp.Stations = append(resp.Stations, station)
			}
		}
	}

	if len(resp.Stations) > 0 {
		resp.ResponseCode = "100"
		resp.ResponseText = "OK"

	} else {
		resp.ResponseCode = "102"
		resp.ResponseText = "No stations found"
	}

	return resp, nil
}

func (s SimulatorServer) GetStationGroups(req *schema.GetStationGroupsRequest) (*schema.GetStationGroupsResponse, error) {
	resp := &schema.GetStationGroupsResponse{}

	for j, org := range s.Ev.facilities {
		if req.OrganizationID != "" && req.OrganizationID != j {
			continue
		}
		for i, g := range org.groups {
			gid, err := strconv.ParseInt(i, 10, 32)
			if err != nil {
				continue
			}
			r := schema.GetStationGroupsResponse_StationGroup{
				OrganizationID: j,
				StationGroupID: int32(gid),
			}
			if org.name != nil {
				r.OrganizationName = *org.name
			}
			if g.name != nil {
				r.StationGroupName = *g.name
			}
			for k, st := range g.stations {
				station := schema.GetStationGroupsResponse_StationGroup_Station{
					StationID: k,
					Coordinate: &schema.Coordinate{
						Latitude:  st.geo.lat,
						Longitude: st.geo.long,
					},
				}
				r.Stations = append(r.Stations, station)
			}
			resp.StationGroups = append(resp.StationGroups, r)
		}
	}

	if len(resp.StationGroups) > 0 {
		resp.ResponseCode = "100"
		resp.ResponseText = "OK"

	} else {
		resp.ResponseCode = "102"
		resp.ResponseText = "No station groups found"
	}

	return resp, nil
}

func (s SimulatorServer) GetCPNInstances(req *schema.GetCPNInstancesRequest) (*schema.GetCPNInstancesResponse, error) {
	resp := &schema.GetCPNInstancesResponse{}

	for i, n := range s.Ev.networks {
		nw := schema.GetCPNInstancesResponse_ChargePointNetwork{
			ID: i,
		}
		if n.name != nil {
			nw.Name = *n.name
		}
		if n.desc != nil {
			nw.Description = *n.desc
		}
		resp.ChargePointNetworks = append(resp.ChargePointNetworks, nw)
	}

	return resp, nil
}

func (s SimulatorServer) ShedLoad(req *schema.ShedLoadRequest) (*schema.ShedLoadResponse, error) {
	resp := &schema.ShedLoadResponse{}

	resp.ResponseCode = "188"
	resp.ResponseText = "Not implemented yet "

	return resp, nil
}

func (s SimulatorServer) ClearShedState(req *schema.ClearShedStateRequest) (*schema.ClearShedStateResponse, error) {
	resp := &schema.ClearShedStateResponse{}

	resp.Success = false
	if req.StationGroupID != nil {
		resp.StationGroupID = *req.StationGroupID
	}
	if req.StationID != nil {
		resp.StationID = *req.StationID
	}
	for _, org := range s.Ev.facilities {
		for i, g := range org.groups {
			id, err := strconv.ParseInt(i, 10, 32)
			if err != nil {
				continue
			}
			if req.StationGroupID != nil && *req.StationGroupID != int32(id) {
				continue
			}

			for k, st := range g.stations {
				if req.StationID != nil && *req.StationID != k {
					continue
				}
				resp.Success = true
				for _, pr := range st.ports {
					pr.sched = 0
				}
			}
		}
	}
	return resp, nil
}
