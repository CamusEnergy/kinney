// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package sim

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

type locGeo struct {
	lat  string
	long string
}

type chargeSample struct {
	time  time.Time
	power float32
}

type vehicle struct {
	driverId     *string
	credentialID *string
	capacity     float32 //KW
	currCharge   float32 // KW
	chargeRate   float32 // KW / KWh
}

type currentCharge struct {
	vehicle      *vehicle
	start        time.Time
	lastComputed time.Time
	chargeRate   float32 // KWh
}

type chargeSession struct {
	vehicle vehicle
	samples []*chargeSample
}

const (
	shedTypeNone = iota
	shedTypeKW
	shedTypePercent
)

type shedPower struct {
	shedType    int
	allowedKW   float32
	percentShed int32
}

type chargePort struct {
	max_capacity     float32 // KWh
	current_capacity float32 // KWh
	shed             *shedPower
	recorded         []*chargeSession
	now              *currentCharge
	shedTimer        *time.Timer
}

type chargeStation struct {
	geo      *locGeo
	name     *string
	address  *string
	ports    map[string]*chargePort
	shedType int
}

type chargeGetLoad interface {
	calcTime(group *chargeGroup) time.Time
	getPortLoad(port *chargePort, t time.Time) (float32, *vehicle, error)
	printGetLoadParams()
}

type chargeGroup struct {
	name        *string
	stations    map[string]*chargeStation
	getLoad     chargeGetLoad
	recalcTimer *time.Timer
}

type chargeFacility struct {
	name   *string
	groups map[string]*chargeGroup
}

type chargeNetwork struct {
	name *string
	desc *string
}

type EVChargers struct {
	lock       sync.Mutex
	drivers    map[string]int
	facilities map[string]*chargeFacility
	networks   map[string]*chargeNetwork
}

const (
	vmwareOrganizationID   = "1:19400001"
	vmwareOrganizationName = "VMware"

	geoDefLat  = "42.63228390329662"
	geoDefLong = "23.378210952553545"

	cpnDefName = "Virtual"
	cpnDefDesc = "EV charger simulator"

	threshold = 60 * time.Minute
)

var (
	printSummary = flag.Bool("print_summary", false, "Print summary of sessions for debugging.")
	printDetail  = flag.Bool("print_detail", false, "Print session details for debugging.")
)

// addChargeProbe adds new recorded charge probe
func (e *EVChargers) addChargeProbe(faciltyID, groupID, stationID, portID, driver, credential *string, t time.Time, power float32) {
	var s *chargeSession
	sample := chargeSample{time: t, power: power}

	port := e.getChargePort(faciltyID, groupID, stationID, portID)
	if len(port.recorded) > 0 {
		lastCharge := port.recorded[len(port.recorded)-1]
		if *lastCharge.vehicle.driverId == *driver && len(lastCharge.samples) > 0 {
			lastSample := lastCharge.samples[len(lastCharge.samples)-1]
			diff := t.Sub(lastSample.time)
			if diff < threshold {
				s = lastCharge
			}
		}
	}

	if s != nil {
		s.samples = append(s.samples, &sample)
	} else {
		s = &chargeSession{
			vehicle: vehicle{
				driverId:     driver,
				credentialID: credential,
			},
			samples: make([]*chargeSample, 0),
		}
		s.samples = append(s.samples, &sample)
		port.recorded = append(port.recorded, s)
		e.drivers[*driver]++
	}
}

func (e *EVChargers) getChargePort(faciltyID, groupID, stationID, id *string) *chargePort {
	s := e.getChargeStation(faciltyID, groupID, stationID, nil, nil, nil)
	if _, ok := s.ports[*id]; ok {
		return s.ports[*id]
	}

	s.ports[*id] = &chargePort{recorded: make([]*chargeSession, 0)}
	return s.ports[*id]
}

func (e *EVChargers) getChargeStation(faciltyID, groupID, id, name, address *string, loc *locGeo) *chargeStation {
	g := e.getChargeGroup(faciltyID, groupID, nil, nil)
	if _, ok := g.stations[*id]; ok {
		return g.stations[*id]
	}

	g.stations[*id] = &chargeStation{
		geo:     loc,
		name:    name,
		address: address,
		ports:   make(map[string]*chargePort),
	}
	return g.stations[*id]
}

func (e *EVChargers) getChargeGroup(faciltyID, id, name *string, gload chargeGetLoad) *chargeGroup {
	f := e.getChargeFacility(faciltyID, nil)
	if _, ok := f.groups[*id]; ok {
		return f.groups[*id]
	}

	f.groups[*id] = &chargeGroup{
		name:     name,
		stations: make(map[string]*chargeStation),
		getLoad:  gload,
	}
	return f.groups[*id]
}

func (e *EVChargers) getChargeFacility(id, name *string) *chargeFacility {
	if _, ok := e.facilities[*id]; ok {
		return e.facilities[*id]
	}

	e.facilities[*id] = &chargeFacility{
		name:   name,
		groups: make(map[string]*chargeGroup),
	}
	return e.facilities[*id]
}

func (e *EVChargers) getCPNetwork(id, name, descriprtion *string) *chargeNetwork {
	if _, ok := e.networks[*id]; ok {
		return e.networks[*id]
	}

	e.networks[*id] = &chargeNetwork{
		name: name,
		desc: descriprtion,
	}
	return e.networks[*id]
}

func NewEvChargers() *EVChargers {
	return &EVChargers{
		drivers:    make(map[string]int),
		facilities: make(map[string]*chargeFacility),
		networks:   make(map[string]*chargeNetwork),
	}
}

// DataPrint prints the entire database with all charger ports
func DataPrint(e *EVChargers) {
	var nGroups int
	var nStations int
	var nPorts int

	fmt.Println("Vehicles:", len(e.drivers))
	for i, v := range e.drivers {
		fmt.Println("\t", i, "charged ", v, "times")

	}
	fmt.Println()

	for o, org := range e.facilities {
		fmt.Println("Organization", *org.name, o)
		nGroups += len(org.groups)
		for i, gr := range org.groups {
			fmt.Printf("\tStation group %s [%s]\n", i, *gr.name)

			gr.getLoad.printGetLoadParams()
			nStations += len(gr.stations)
			for j, st := range gr.stations {
				fmt.Printf("\t\tStation %s [%s], %d ports\n", j, *st.name, len(st.ports))
				nPorts += len(st.ports)
				if !*printSummary && !*printDetail {
					continue
				}
				for k, pr := range st.ports {
					fmt.Printf("\t\t\tPort %s, charges: %d", k, len(pr.recorded))
					if !*printSummary && !*printDetail {
						continue
					}
					for _, k := range pr.recorded {
						chargeTime := k.samples[len(k.samples)-1].time.Sub(k.samples[0].time)
						fmt.Println("\t Driver [", *k.vehicle.driverId, *k.vehicle.credentialID, "] probes",
							len(k.samples), "Total time:", chargeTime)
						if !*printDetail {
							continue
						}
						for _, s := range k.samples {
							fmt.Println("\t\t", s.time, s.power)
						}
					}
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}
	fmt.Printf("\nTotal %d facilities, %d groups, %d stations, %d ports, %d recorded vehicles\n",
		len(e.facilities), nGroups, nStations, nPorts, len(e.drivers))
}
