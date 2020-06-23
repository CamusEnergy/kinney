package simulator

import (

)

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"time"
)

// TODO test without SOAP etc, making directly the method calls
// use sleep to move in time .. but that would be slow
// may want to add a unit time to fast forward .. maintaining an internal clock and from there reading time -- nice!
// define a little fast_forward_function (n int, units seconds|minutes|hours|days)

var (
	empty_session chargeSession
)

type vehicle struct {
	ID                      		int
	ownerID                        	string
	capacity                        float32
	curr_charge                     float32
	charge_rate                     float32
}

func newVehicle(ID int, ownerID string, capacity, curr_charge, charge_rate float32) vehicle {
	if (0 == ID)  || ("" == ownerID)  || (0.0 == capacity) || (curr_charge < 0.0) || (charge_rate <= 0.0) {
		log.Fatal("Invalid input provided to initialize a vehicle ID =  %d," +
			" ownerID = %s, capacity %v, curr_charge: %v, charge_rate = %v",
			ID, ownerID, capacity, curr_charge, charge_rate)
	}
	return vehicle{ID, ownerID, capacity, curr_charge, charge_rate}
}

type chargeSession struct {
	ID                      int
	vehicle                 vehicle
	start                   time.Time
	last_computed			time.Time
	end                     time.Time
	total_charge    		float32
	full_port_ID           	string
}

func newChargeSession(ID int, v vehicle, full_port_id string) chargeSession {
	now = time.Now()
	return chargeSession{ID, v, now, now, now, 0.0,  full_port_id}
}

type port struct {
	ID                      int
	max_capacity    		float32
	capacity                float32
	shed                    bool
	session                 chargeSession
}

type station struct {
	ID              string
	shed            bool
	ports           [2]port
}

type stationGroup struct {
	sgID            int
	stations        [10]station
	num_stations	int
	shed        	bool
}

type chargeFacility struct {
	sg 					stationGroup
	created				time.Time
	now					time.Time
	completed_sessions 	list.List
	last_sessionID  	int
}

func NewChargeFacility(sgID, num_stations int,  max_capacity float32) chargeFacility {
	var stations = [10]station{}
	for i:= 0; i < num_stations; i++ {
		var ports = [2]port{}
		for j:= 0; j < 2; j++ {
			ports[j] = port{ID: j, max_capacity: max_capacity, capacity: max_capacity, shed: false, session: empty_session}
		}
		stations[i] = station{i, false, ports}
	}
	var sg = stationGroup{sgID: sgID, stations: stations, shed: false}
	return chargeFacility{sg: sg, last_sessionID: 0}
}

func uniquePortID(sgID int, station_id string,  port_id int) string {
	return fmt.Sprintf("%d*%s*%d", sgID, station_id, port_id)
}


func (cf chargeFacility) plugin(v vehicle) bool {
	for _, s := range cf.sg.stations {
		var m  sync.Mutex
		m.Lock() // only on go thread at a time can modify this shared data structure
		for _, p := range s.ports {
			if p.session == empty_session {
				cf.last_sessionID = cf.last_sessionID +1 // obtain the next new session ID
				p.session = newChargeSession(cf.last_sessionID, v, uniquePortID(cf.sg.sgID, s.ID, p.ID))
				fmt.Println("Started new charge_session: %v", p.session)
				return true
			}
		}
		m.Unlock()
	}
	fmt.Println("No vacant port found for vehicle plugin")
	return false
}

func (cf chargeFacility) unplug(v vehicle) bool {
	var charge_session = empty_session
	for _, s := range cf.sg.stations {
		for _, p := range s.ports {
			if (p.session != empty_session) && (p.session.vehicle == v) {
				charge_session = p.session
				cf.completed_sessions.PushFront(charge_session)
				p.session = empty_session
				return true
			}
		}

	}
	return false
}

// The vehicle charge rate limits the charge rate if it is less than the port load capacity
// and if port_capacity is less than vehicle charge rate, it is the limiting factor
// Is the vehicle fully charged?
//if (now - session-start) * charge_rate + current_charge
// TODO -- return in the ChargePoint getLoad struct
func (cf chargeFacility) getLoad() float32 {
	var now = time.Now()
	var total_load  float32
	for i := 0; i < cf.sg.num_stations; i++ {
		var s = cf.sg.stations[i]
		for j, p := range s.ports {
			var port_load = float32(0.0)
			// is a vehicle connected at this port?
			if (p.session != empty_session) {
				var v = p.session.vehicle
				// what is the vehicle's charge rate? the lower of the port capacity and the vehicle's charge rate
				var vehicle_charge_rate = v.charge_rate
				if (vehicle_charge_rate < p.capacity) {
					port_load = vehicle_charge_rate
				} else{
					port_load = p.capacity
				}
				// is the vehicle fully charged?
				var last_charge = v.curr_charge
				var duration = now.Second() - p.session.last_computed.Second()
				var amount = port_load * float32(duration) / float32(60 * 60)
				if ((amount + last_charge) > v.capacity) {
					port_load = 0.0
					v.curr_charge = v.capacity
				} else {
					v.curr_charge = amount + last_charge
				}
				// update port session
				p.session.last_computed = now
				total_load = port_load
			}
			fmt.Println("station[%d] Port [%j] load = %f", i, j, port_load)
		}
	}
    return total_load
}

func (cf chargeFacility) shed(amount float32, percent bool) {
	cf.sg.shed = true
	var s station
	for i := 0; i < cf.sg.num_stations; i++ {
		s = cf.sg.stations[i]
		fmt.Printf("station[%d] = %s\n", i, s)
		s.shed = true
		for j, p := range s.ports {
			fmt.Println("Prior to shed port[%d] capacity = %f", j, p.capacity)
			if percent {
				p.capacity = p.max_capacity * amount * 0.01
			} else {
				p.capacity = amount
			}
			fmt.Println("Post shed port[%d] capacity = %f", j, p.capacity)
		}
	}
}

func (cf chargeFacility)clear() {
	cf.sg.shed = false
	for i, s := range cf.sg.stations {
		s.shed = false
		for j, p := range s.ports {
			p.shed = false
			p.capacity = p.max_capacity
			fmt.Println("Cleared shed on station[%d] port[%d]", i, j)
		}
	}
}

