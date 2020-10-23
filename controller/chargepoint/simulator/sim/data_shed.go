// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package sim

import (
	"fmt"
	"strconv"
	"time"
)

func (e *EVChargers) clearShedLoad(groupID *int32, stationID *string, ports *[]string) bool {
	res := false

	for _, org := range e.facilities {
		for i, g := range org.groups {
			id, err := strconv.ParseInt(i, 10, 32)
			if err != nil {
				continue
			}
			if groupID != nil && *groupID != 0 && *groupID != int32(id) {
				continue
			}
			for k, st := range g.stations {
				success := false
				if stationID != nil && len(*stationID) > 0 && *stationID != k {
					continue
				}
				if len(*ports) > 0 {
					for _, p := range *ports {
						if pr, ok := st.ports[p]; ok {
							success = true
							pr.current_capacity = pr.max_capacity
							pr.shed = nil
						}
					}
				} else {
					for _, pr := range st.ports {
						success = true
						pr.current_capacity = pr.max_capacity
						pr.shed = nil
					}
				}
				if success {
					res = true
					checkStationShed(st)
					// Calculate new charging load according to the new shed state
					e.getNextLoad(nil, g, &k)
				}
			}
		}
	}
	return res
}

func getShedPower(port *chargePort, power *shedPower) float32 {
	if power.shedType == shedTypePercent {
		if port.now != nil {
			return port.now.chargeRate - (float32(power.percentShed) * port.now.chargeRate / 100)
		} else {
			return 0.0
		}
	} else if power.shedType == shedTypeKW {
		return power.allowedKW
	}

	return port.current_capacity
}

func checkStationShed(station *chargeStation) {
	for _, p := range station.ports {
		if p.shed != nil && p.shed.shedType != shedTypeNone {
			return
		}
	}
	station.shedType = shedTypeNone
}

func (e *EVChargers) schedPort(group *chargeGroup, stationID *string, port *chargePort, power *shedPower, minutes int32) {
	newPower := getShedPower(port, power)
	fmt.Println("Shed port from", port.current_capacity, "to", newPower, "for", minutes, "minutes")
	port.current_capacity = newPower
	port.shed = power
	if port.shedTimer != nil {
		if !port.shedTimer.Stop() {
			<-port.shedTimer.C
		}
		port.shedTimer = nil
	}
	if minutes > 0 {
		port.shedTimer = time.AfterFunc(time.Duration(minutes)*time.Minute, func() {
			e.lock.Lock()
			defer e.lock.Unlock()
			port.shedTimer = nil
			/* Restore port capacity */
			port.current_capacity = port.max_capacity
			port.shed = nil
			if st, ok := group.stations[*stationID]; ok {
				checkStationShed(st)
			}
			if group.recalcTimer == nil {
				/* Recalculate vehicles being charged, with the new port capacity */
				group.recalcTimer = time.AfterFunc(time.Duration(2)*time.Second, func() {
					e.lock.Lock()
					defer e.lock.Unlock()
					group.recalcTimer = nil
					e.getNextLoad(nil, group, stationID)
				})
			}
		})
	}
}

func (e *EVChargers) shedLoad(groupID *int32, stationID *string, power *shedPower, ports *map[string]*shedPower, minutes int32) bool {
	res := false
	shedRequested := shedTypeNone

	if power != nil {
		shedRequested = power.shedType
	} else if ports != nil && len(*ports) > 0 {
		for _, p := range *ports {
			shedRequested = p.shedType
			break
		}
	} else {
		return false
	}

	for _, org := range e.facilities {
		for i, g := range org.groups {
			id, err := strconv.ParseInt(i, 10, 32)
			if err != nil {
				continue
			}
			if groupID != nil && *groupID != 0 && *groupID != int32(id) {
				continue
			}
			for k, st := range g.stations {
				if stationID != nil && len(*stationID) > 0 && *stationID != k {
					continue
				}
				/* Check shed mode mismatch for the requested stations */
				if st.shedType != shedTypeNone && st.shedType != shedRequested {
					return false
				}
			}
			if g.recalcTimer != nil {
				if !g.recalcTimer.Stop() {
					<-g.recalcTimer.C
				}
				g.recalcTimer = nil
			}
			for k, st := range g.stations {
				success := false
				if stationID != nil && *stationID != "" && *stationID != k {
					continue
				}
				if st.name != nil {
					fmt.Println("Shed ports in station", *st.name)
				}
				if ports != nil && len(*ports) > 0 {
					for portID, portPower := range *ports {
						if port, ok := st.ports[portID]; ok {
							success = true
							e.schedPort(g, &k, port, portPower, minutes)
						}
					}
				} else if power != nil {
					spower := *power
					if spower.shedType == shedTypeKW {
						spower.allowedKW = spower.allowedKW / float32(len(st.ports))
					}
					for _, port := range st.ports {
						success = true
						e.schedPort(g, &k, port, &spower, minutes)
					}
				}
				if success {
					res = true
					st.shedType = shedRequested
					// Calculate new charging load according to the new shed state
					e.getNextLoad(nil, g, &k)
				}
			}
		}
	}
	return res
}
