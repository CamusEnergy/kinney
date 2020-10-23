// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package sim

import (
	"fmt"
	"math/rand"
	"time"
)

type getLoadRandom struct {
	// Probability of a port to have a charge session, in %
	portLoad int

	// Maximum capacity of a vehicle battery
	maxVehicleBattery int
}

func (g getLoadRandom) newChargeSession(port *chargePort, t time.Time) {
	// Delete the old session, if any
	port.now = nil
	// generate new session with the given probability
	if rand.Intn(100) > g.portLoad {
		return
	}

	port.now = &currentCharge{
		vehicle: &vehicle{
			driverId:     randString(1, ""),
			credentialID: randString(1, ""),
			capacity:     float32(randParam(g.maxVehicleBattery/2, g.maxVehicleBattery)),
		},
		start:        t,
		lastComputed: t,
	}
	// Initial battery charge level, random value in interval (0, battery capacity / 2)
	port.now.vehicle.currCharge = float32(randParam(0, int(port.now.vehicle.capacity/2)))

	// Battery charge rate, in KW / KWh of unfilled battery capacity
	// random value/10 in the interval(battery capacity/40, battery capacity/10)
	// example: 80KW battery -> charge rate, random value in the interval [0.2, 0.8], -> 0.3 KW / KWh
	//			unfilled battery capacity 60KW -> max charge rate 18KWh
	//			unfilled battery capacity 20KW -> max charge rate 6KWh
	port.now.vehicle.chargeRate = float32(randParam(int(port.now.vehicle.capacity/40), int(port.now.vehicle.capacity/10))) / 10

	// Current charge rate, limited to max port capacity
	port.now.chargeRate = port.now.vehicle.chargeRate * (port.now.vehicle.capacity - port.now.vehicle.currCharge)
	if port.now.chargeRate > port.current_capacity {
		port.now.chargeRate = port.current_capacity
	}

	fmt.Printf("New vehicle [%s]\n", *port.now.vehicle.driverId)
}

func (g getLoadRandom) calcNextLoad(port *chargePort, t time.Time) {
	chargeTime := time.Since(port.now.lastComputed).Minutes()

	port.now.vehicle.currCharge += (port.now.chargeRate * float32(chargeTime)) / 60
	fmt.Printf("Vehicle [%s] %f/%f KW charged, ", *port.now.vehicle.driverId, port.now.vehicle.currCharge, port.now.vehicle.capacity)

	if port.now.vehicle.currCharge >= port.now.vehicle.capacity {
		// Vehicle is fully charged, unplug it
		port.now = nil
	} else {
		port.now.chargeRate = port.now.vehicle.chargeRate * (port.now.vehicle.capacity - port.now.vehicle.currCharge)
		if port.now.chargeRate > port.current_capacity {
			port.now.chargeRate = port.current_capacity
		}
		port.now.lastComputed = time.Now()
	}
	fmt.Printf("charging @ %f KWh\n", port.now.chargeRate)
}

func (g getLoadRandom) calcTime(group *chargeGroup) time.Time {
	return time.Now()
}

func (g getLoadRandom) getPortLoad(port *chargePort, t time.Time) (float32, *vehicle, error) {
	if port.now == nil {
		g.newChargeSession(port, t)
	}

	if port.now != nil {
		g.calcNextLoad(port, t)
		if port.now != nil {
			return port.now.chargeRate, port.now.vehicle, nil
		}
	}

	return 0.0, nil, nil
}

func (g getLoadRandom) printGetLoadParams() {
	fmt.Println("\t\tProbability of a port to have a charge session", g.portLoad, "%")
}
