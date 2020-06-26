package simulator

import (
	"fmt"
	"testing"
)
/*
func TestVehicleCreation(t *testing.T) {
	var v = NewVehicle(1234,"Someone",75.0, 10.0, 6.0)
	var want = float32(75.0)
	if got := v.capacity; got != want {
		t.Errorf("v.capacity = %f, want %f", got, want)
	}
}
func TestChargeFacilityCreation(t *testing.T) {
	var cf = NewChargeFacility(1234, 1,8.0, "somewhere groovy")
	var want = float32(0.0)
	if got := cf.GetLoad(); got != want {
			t.Errorf("cf.GetLoad() = %f, want %f", got, want)
	}
}

func TestPlugin(t *testing.T) {
	var numStations = 1
	var v1 = NewVehicle(11,"Someone One",75.0, 10.0, 6.5)
	var v2 = NewVehicle(22,"Someone Two",45.0, 20.0, 6.0)
	var v3 = NewVehicle(33,"Someone Three",50.0, 40.0, 5.0)
	var cf = NewChargeFacility(1234, numStations,8.0, "somewhere groovy")
	//cf.showPorts(numStations, "Before plugin")
	var want = true
	if got := cf.Plugin(&v1); got != want {
		t.Errorf("cf.Plugin() = %t, want %t", got, want)
	}
	//cf.showPorts(numStations, "After plugin 1")
	// second vehicle plugin
	want = true
	if got := cf.Plugin(&v2); got != want {
		t.Errorf("cf.Plugin() = %t, want %t", got, want)
	}
	//cf.showPorts(numStations, "After plugin 2")
	// third vehicle plugin
	want = false
	if got := cf.Plugin(&v3); got != want {
		t.Errorf("cf.Plugin() = %t, want %t", got, want)
	}
	//cf.showPorts(numStations, "After plugin 3 -- should not change anything because full")
	want = true
	if got := cf.Unplugin(&v1); got != want {
		t.Errorf("cf.Unplugin() = %t, want %t", got, want)
	}
	want = true
	if got := cf.Unplugin(&v2); got != want {
		t.Errorf("cf.Unplugin() = %t, want %t", got, want)
	}
	//cf.showPorts(numStations, "After unplugin v2")
	want = false
	if got := cf.Unplugin(&v3); got != want {
		t.Errorf("cf.Unplugin() for v3 = %t, want %t", got, want)
	}
	//cf.showPorts(numStations, "After unplugin v3")
	want = true
}*/

func TestShed(t *testing.T) {
	var cf = NewChargeFacility(1234, 2, 8.0, "somewhere groovy")
	//cf.showPorts(2, "============ pre shed ============")
	//var valWant = float32(3.0)
	//cf.Shed(valWant, false)
	cf.sg.shed = true
	for i := 0; i < cf.sg.numStations; i++ {
		var s = cf.sg.stations[i]
		s.shed = true
		for j, p := range s.ports {
			fmt.Printf("Before Station-Port[%d, %d] = (%#v), capacity: %f, max_capacity %f\n", i, j, p.session, p.capacity, p.maxCapacity)
			p.capacity = float32(3.0)
			fmt.Printf("After Station-Port[%d, %d] = (%#v), capacity: %f, max_capacity %f\n", i, j, p.session, p.capacity, p.maxCapacity)
		}
	}
	//cf.showPorts(2, "============ post shed ============")
	cf.sg.shed = false
	for i := 0; i < cf.sg.numStations; i++ {
		var s = cf.sg.stations[i]
		s.shed = false
		for j, p := range s.ports {
			fmt.Printf("Before clear Station-Port[%d, %d] = (%#v), capacity: %f, max_capacity %f\n", i, j, p.session, p.capacity, p.maxCapacity)
			p.capacity = p.maxCapacity
			p.shed = false
			fmt.Printf("After clear Station-Port[%d, %d] = (%#v), capacity: %f, max_capacity %f\n", i, j, p.session, p.capacity, p.maxCapacity)


		}
	}
	cf.showPorts(2, "============ post clear ============")
/*
	var want = false
	if got := cf.sg.shed; got != want {
		t.Errorf("Pre shed  flag = %t, want %t", got, want)
	}
	var p = cf.sg.stations[0].ports[0]
	var valWant = p.maxCapacity
	var valGot = p.capacity
	if  valGot != valWant {
		t.Errorf("Pre shed capacity %f does not equal max capacity %f\n",
		valGot, valWant)
	}
	want = true
	valWant = float32(3.0)
	cf.Shed(valWant, false)
	cf.showPorts(2, "============ post shed ============")
	// shedding load to valWant
	// checking shed flag first
	if got := cf.sg.shed; got != want {
		t.Errorf("shed flag = %t, want %t", got, want)
	}
	// check if capacity reduced
	valGot = cf.sg.stations[0].ports[0].capacity
	if valGot != valWant {
		t.Errorf("shed capacity = %f want %f", valGot, valWant)
	}*/
}
/*
func TestClear(t *testing.T) {
	var cf = NewChargeFacility(1234, 1, 8.0, "somewhere groovy")
	cf.Shed(float32(5.0), false)
	var want = true
	if got := cf.sg.shed; got != want {
		t.Errorf("cf.shed = %t, want %t", got, want)
	}
	var p = cf.sg.stations[0].ports[0]
	if p.capacity > p.maxCapacity {
		t.Errorf("Pre clear capacity not less than max_capacity!")
	}
	cf.Clear()
	want = false
	if got := cf.sg.shed; got != want {
		t.Errorf("cf.shed = %t, want %t", got, want)
	}
	p = cf.sg.stations[0].ports[0]
	if p.capacity != p.maxCapacity{
		t.Errorf("Post clear, capacity must equal max_capacity!")
	}
}*/