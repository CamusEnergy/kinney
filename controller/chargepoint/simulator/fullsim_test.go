package simulator

import (
	"testing"
)

func TestVehicleCreation(t *testing.T) {
	var v = NewVehicle(1234,"Someone",75.0, 10.0, 6.0)
	var want = float32(75.0)
	if got := v.capacity; got != want {
		t.Errorf("v.capacity = %f, want %f", got, want)
	}
}
func TestChargeFacilityCreation(t *testing.T) {
	var cf = NewChargeFacility(1234, 1,2, 8.0, "somewhere groovy")
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
	var cf = NewChargeFacility(1234, numStations,2, 8.0, "somewhere groovy")
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
	if got := cf.Unplug(&v1); got != want {
		t.Errorf("cf.Unplug() = %t, want %t", got, want)
	}
	want = true
	if got := cf.Unplug(&v2); got != want {
		t.Errorf("cf.Unplug() = %t, want %t", got, want)
	}
	//cf.showPorts(numStations, "After unplug v2")
	want = false
	if got := cf.Unplug(&v3); got != want {
		t.Errorf("cf.Unplug() for v3 = %t, want %t", got, want)
	}
	//cf.showPorts(numStations, "After unplug v3")
	want = true
}

func TestShed(t *testing.T) {
	var cf = NewChargeFacility(1234, 2, 2, 8.0, "somewhere groovy")
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
	// shedding load to valWant
	// checking shed flag
	if got := cf.sg.shed; got != want {
		t.Errorf("shed flag = %t, want %t", got, want)
	}
	// check if capacity curtailed down from maxCapacity to valWant
	valGot = cf.sg.stations[0].ports[0].capacity
	if valGot != valWant {
		t.Errorf("shed capacity = %f want %f", valGot, valWant)
	}
}


func TestClear(t *testing.T) {
	var cf = NewChargeFacility(1234, 1, 2, 8.0, "somewhere groovy")
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
}