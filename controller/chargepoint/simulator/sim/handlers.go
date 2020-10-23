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

	s.Ev.lock.Lock()
	defer s.Ev.lock.Unlock()

	err = s.Ev.getNextLoad(resp, group, &req.StationID)
	if err != nil {
		resp.ResponseCode = "102"
		resp.ResponseText = "No load recorded"
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
	var shedKW *string
	var shedPercent *int32
	var portShed map[string]*shedPower
	var shed *shedPower
	portShedKW := false
	portSchedPercent := false
	success := false

	resp := &schema.ShedLoadResponse{Success: 0,
		StationGroupID: req.StationGroupID,
		StationID:      req.StationID,
	}
	resp.ResponseCode = "171"

	if len(req.StationGroupAllowedLoadKW) > 0 {
		shedKW = &req.StationGroupAllowedLoadKW
	}
	if req.StationGroupPercentShed != nil {
		if shedKW != nil {
			resp.ResponseText = "shed KW requested, invalid StationGroupPercentShed"
			return resp, nil
		}
		shedPercent = req.StationGroupPercentShed
	}

	if len(req.StationID) > 0 && (shedKW != nil || shedPercent != nil) {
		resp.ResponseText = "shed requested for all stations in the group, invalid StationID"
		return resp, nil
	}

	if len(req.StationAllowedLoadKW) > 0 {
		if shedKW != nil || shedPercent != nil {
			resp.ResponseText = "shed already requested, invalid StationAllowedLoadKW"
			return resp, nil
		}
		shedKW = &req.StationAllowedLoadKW
	}
	if req.StationPercentShed != nil {
		if shedKW != nil || shedPercent != nil {
			resp.ResponseText = "shed already requested, invalid StationPercentShed"
			return resp, nil
		}
		shedPercent = req.StationPercentShed
	}

	if len(req.Ports) > 0 {
		if shedKW != nil || shedPercent != nil {
			resp.ResponseText = "station shed already requested, invalid Ports"
			return resp, nil
		}
		portShed = make(map[string]*shedPower)
		for _, p := range req.Ports {
			if len(p.AllowedLoadKW) > 0 {
				if portSchedPercent {
					resp.ResponseText = "port % shed already requested, invalid AllowedLoadKW"
					return resp, nil
				}
				if power, err := strconv.ParseFloat(p.AllowedLoadKW, 32); err == nil {
					portShedKW = true
					portShed[p.PortNumber] = &shedPower{
						shedType:  shedTypeKW,
						allowedKW: float32(power),
					}
				}
			}
			if p.PercentShed != nil {
				if portShedKW {
					resp.ResponseText = "port KW shed already requested, invalid PercentShed"
					return resp, nil
				}
				portSchedPercent = true
				portShed[p.PortNumber] = &shedPower{
					shedType:    shedTypePercent,
					percentShed: *p.PercentShed,
				}
			}
		}
	}

	s.Ev.lock.Lock()
	defer s.Ev.lock.Unlock()

	if shedKW != nil {
		if power, err := strconv.ParseFloat(*shedKW, 32); err == nil {
			shed = &shedPower{
				shedType:  shedTypeKW,
				allowedKW: float32(power),
			}
		}
	} else if shedPercent != nil {
		shed = &shedPower{
			shedType:    shedTypePercent,
			percentShed: *shedPercent,
		}
	}

	if shed != nil || portShed != nil {
		success = s.Ev.shedLoad(&req.StationGroupID, &req.StationID, shed, &portShed, req.TimeInterval)
	}

	if success {
		resp.Success = 1
		resp.ResponseCode = "100"
		resp.ResponseText = "OK"
		if shedKW != nil {
			resp.AllowedLoadKW = *shedKW
		}
		resp.PercentShed = shedPercent
		for _, p := range req.Ports {
			resp.Ports = append(resp.Ports, schema.ShedLoadResponse_Port{
				PortNumber:    p.PortNumber,
				AllowedLoadKW: p.AllowedLoadKW,
				PercentShed:   p.PercentShed,
			})
		}
	} else {
		resp.ResponseCode = "126"
		resp.ResponseText = "Failed to shed requested load"
	}

	return resp, nil
}

func (s SimulatorServer) ClearShedState(req *schema.ClearShedStateRequest) (*schema.ClearShedStateResponse, error) {
	resp := &schema.ClearShedStateResponse{Success: false}

	if req.StationGroupID != nil {
		resp.StationGroupID = *req.StationGroupID
	}
	if req.StationID != nil {
		resp.StationID = *req.StationID
	}

	if len(req.PortNumbers) > 0 && req.StationID == nil {
		resp.ResponseCode = "119"
		resp.ResponseText = "Invalid Station ID"
		return resp, nil
	}

	s.Ev.lock.Lock()
	defer s.Ev.lock.Unlock()

	resp.Success = s.Ev.clearShedLoad(req.StationGroupID, req.StationID, &req.PortNumbers)

	if resp.Success {
		resp.ResponseCode = "100"
		resp.ResponseText = "OK"
	} else {
		resp.ResponseCode = "127"
		resp.ResponseText = "Failed to restore the load"
	}

	return resp, nil
}
