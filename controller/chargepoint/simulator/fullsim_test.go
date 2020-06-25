package simulator


import (
	"fmt"
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
	var cf = NewChargeFacility(1234, 1,8.0)
	var want = float32(0.0)
	if got := cf.GetLoad(); got != want {
			t.Errorf("cf.GetLoad() = %f, want %f", got, want)
	}
}

func TestPlugin(t *testing.T) {
	fmt.Printf("Empty Session: %#v\n", emptySession)
	var v = NewVehicle(1234, "Someone",75.0, 10.0, 6.0)
	fmt.Printf("%#v\n", v)
	var cf = NewChargeFacility(1234, 1,8.0)
	var want = true
	fmt.Println("******* Port Sessions Before Plugin ******")
	cf.showPorts()
	if got := cf.Plugin(v); got != want {
		t.Errorf("cf.Plugin(v) = %t, want %t", got, want)
	}
	fmt.Println("******* Port Sessions Before Plugin ******")
	cf.showPorts()

	if got := cf.Unplug(v); got != want {
		t.Errorf("cf.UnPlug(v) = %t, want %t", got, want)
	}
	fmt.Println("******* Port Sessions After Unplugin ******")
	cf.showPorts()
}