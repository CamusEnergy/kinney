package simulator

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"
)

// TODO malini
// may want to add a unit time to fast forward .. maintaining an internal clock and from there reading time
// define a little fast_forward_function (n int, units seconds|minutes|hours|days)


type vehicle struct {
	ID                      		int
	ownerID                        	string
	capacity                        float32
	currCharge                     	float32
	chargeRate                     	float32
}

func NewVehicle(ID int, ownerID string, capacity, currCharge, chargeRate float32) vehicle {
	if (0 == ID)  || ("" == ownerID)  || (0.0 == capacity) || (currCharge < 0.0) || (chargeRate <= 0.0) {
		log.Fatal("Invalid input provided to initialize a vehicle ID =  %d," +
			" ownerID = %s, capacity %f, currCharge: %f, chargeRate = %f",
			ID, ownerID, capacity, currCharge, chargeRate)
	}
	return vehicle{ID, ownerID, capacity, currCharge, chargeRate}
}

type chargeSession struct {
	ID                      int
	vehicle                 *vehicle
	start                   time.Time
	lastComputed			time.Time
	end                     time.Time
	totalCharge    		float32
	fullPortID           	string
}

func newChargeSession(ID int, v *vehicle, fullPortID string) *chargeSession {
	var now = time.Now()
	fmt.Printf("Starting charge session for vehicle: %+v\n", v)
	return &chargeSession{ID, v, now, now, now, 0.0, fullPortID }
}

type port struct {
	ID                      int
	maxCapacity    			float32
	capacity                float32
	shed                    bool
	session                 *chargeSession
}

type station struct {
	ID              string
	shed            bool
	ports           [2]port
}

// TODO -- make variable size using lists .. unsure if one can use a slice in a struct
type stationGroup struct {
	sgID            int
	numStations		int
	stations        [10]station
	shed        	bool
}

type chargeFacility struct {
	sg 					stationGroup
	address             string
	created				time.Time
	now					time.Time
	completedSessions 	list.List
	lastSessionID  		int
	m  					sync.Mutex
}

func NewChargeFacility(sgID, numStations int,  maxCapacity float32, address string) chargeFacility {
	var stations = [10]station{}
	for i:= 0; i < numStations; i++ {
		var ports = [2]port{}
		for j:= 0; j < 2; j++ {
			ports[j] = port{ID: j, maxCapacity: maxCapacity, capacity: maxCapacity, shed: false, session: nil}
		}
		// "1:NNNN" would be a US station ID
		stations[i] = station{fmt.Sprintf("1:%d", i), false, ports}
	}
	var sg = stationGroup{sgID: sgID, numStations: numStations, stations: stations, shed: false}
	var cf =  chargeFacility{sg: sg, address: address}
	return cf
}

func uniquePortID(sgID int, station_id string,  port_id int) string {
	return fmt.Sprintf("%d*%s*%d", sgID, station_id, port_id)
}

func (cf *chargeFacility) showPorts(numStations int, msg string) {
	fmt.Printf("-------------%s------------\n", msg)
	for i:=0; i < numStations; i++ {
		var s = cf.sg.stations[i]
		for j, p := range s.ports {
			fmt.Printf("Station-Port[%d, %d] = (%#v), capacity: %f, max_capacity %f\n", i, j, p.session, p.capacity, p.maxCapacity)
			fmt.Printf("%#v\n", p)
		}
	}
}

func (cf *chargeFacility) Plugin(v *vehicle) bool {
	cf.m.Lock()
	defer cf.m.Unlock()
	for i := 0; i < cf.sg.numStations; i++ {
		var s = cf.sg.stations[i]
		for j := 0; j < 2; j++ {
			var p = s.ports[j]
			if cf.sg.stations[i].ports[j].session == nil {
				cf.lastSessionID = cf.lastSessionID + 1 // obtain the next new session ID
				cf.sg.stations[i].ports[j].session = newChargeSession(cf.lastSessionID, v, uniquePortID(cf.sg.sgID, s.ID, p.ID))
				return true
			}
		}
	}
	fmt.Println("No vacant port found for vehicle plugin")
	return false
}

func (cf *chargeFacility) Unplugin(v *vehicle) bool {
	cf.m.Lock()
	defer cf.m.Unlock()
	var chargeSession *chargeSession
	for i := 0; i < cf.sg.numStations; i++ {
		for j := 0; j < 2; j++ {
			if cf.sg.stations[i].ports[j].session!= nil {
				if  cf.sg.stations[i].ports[j].session.vehicle == v {
					chargeSession = cf.sg.stations[i].ports[j].session
					// update that session is complete
					chargeSession.end = time.Now()
					cf.completedSessions.PushFront(chargeSession)
					cf.sg.stations[i].ports[j].session = nil
					return true
				}
			}
		}
	}
	fmt.Println("No such vehicle currently plugged in")
	return false
}

// The vehicle charge rate limits the charge rate if it is less than the port load capacity
// and if port_capacity is less than vehicle charge rate, it is the limiting factor
// Is the vehicle fully charged?
//if (now - session-start) * charge_rate + current_charge
// TODO -- return in the ChargePoint getLoad struct, full details, total and per port across all stations
func (cf *chargeFacility) GetLoad() float32 {
	cf.m.Lock()
	defer cf.m.Unlock()
	var now = time.Now()
	var totalLoad = float32(0.0)
	for i := 0; i < cf.sg.numStations; i++ {
		var s = cf.sg.stations[i]
		for _, p := range s.ports {
			var portLoad = float32(0.0)
			// is a vehicle connected at this port?
			if p.session != nil {
				var v = p.session.vehicle
				// what is the vehicle's charge rate? the lower of the port capacity and the vehicle's charge rate
				var vehicleChargeRate = v.chargeRate
				if vehicleChargeRate < p.capacity {
					portLoad = vehicleChargeRate
				} else{
					portLoad = p.capacity
				}
				// is the vehicle fully charged?
				var lastCharge = v.currCharge
				var duration = now.Second() - p.session.lastComputed.Second()
				var amount = portLoad * float32(duration) / float32(60 * 60)
				if (amount + lastCharge) > v.capacity {
					portLoad = 0.0
					v.currCharge = v.capacity
				} else {
					v.currCharge = amount + lastCharge
				}
				// update port session
				p.session.lastComputed = now
				totalLoad = portLoad + totalLoad
			}
		}
	}
    return totalLoad
}

func (cf *chargeFacility) Shed(amount float32, percent bool) {
	fmt.Printf("enter shed with amount = %f\n", amount)
	cf.m.Lock()
	defer cf.m.Unlock()
	cf.sg.shed = true
	var s station
	for i := 0; i < cf.sg.numStations; i++ {
		s = cf.sg.stations[i]
		s.shed = true
		for _, p := range s.ports {
			if percent {
				p.capacity = p.maxCapacity * amount * 0.01
			} else {
				p.capacity = amount
			}
		}
	}
}

func (cf *chargeFacility) Clear() {
	cf.m.Lock()
	defer cf.m.Unlock()
	cf.sg.shed = false
	for i := 0; i < cf.sg.numStations; i++  {
		var s = cf.sg.stations[i]
		s.shed = false
		for j, p := range s.ports {
			p.shed = false
			p.capacity = p.maxCapacity
			fmt.Printf("Cleared shed on station[%d] port[%d]\n", i, j)
		}
	}
}
