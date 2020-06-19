package sim

import (
	"fmt"
	"time"
)

type evGeo struct {
		Lat string
		Long string
}

type chargeSample struct {
	Time		time.Time
	Power		float32	
}

type chargeSession struct {
	VehicleId	string
	PortId		string
	Samples		[]chargeSample
}


type evPort struct {
	Sched int
	Charge []*chargeSession
}

type evStation struct {
	Geo evGeo
	Ports  map[string]evPort
}

type evStGroup struct {
	Name string
	Stations map[string]evStation
}

type evOrg struct {
	Name string
	StGrpoups map[string]evStGroup
}

type evCPN struct {
	Name string
	Desc string
}

type EvChargers struct {
	Vehicles map[string]int
	StOrgs map[string]evOrg
	StCPNs map[string]evCPN
}

const (
	PrintSessions uint8 = 1 << iota
	PrintAll
)

func DataPrint(e *EvChargers, verbose uint8) {
	fmt.Println("Vechicles:", len(e.Vehicles))
	for i,v := range e.Vehicles {
		fmt.Println("\t", i,"charged ", v, "times")
		
	}
	fmt.Println()
	
	for o,org := range e.StOrgs {
		fmt.Println("Organization", org.Name, o)
		for i,gr := range org.StGrpoups {
			fmt.Println("\tStation group", i)
			for j,st := range gr.Stations {
				fmt.Println("\t\tStation", j)
				for k,pr := range st.Ports {
					fmt.Println("\t\t\tPort", k, ", charges:", len(pr.Charge))
					if (verbose & PrintSessions) == 0 {
						continue;
					}
					for _,k := range pr.Charge {
						chargeTime := k.Samples[len(k.Samples)-1].Time.Sub(k.Samples[0].Time)
						fmt.Println("\t Vehicle [", k.VehicleId, "] probes", len(k.Samples), "Total time:", chargeTime)
						if (verbose & PrintAll) == 0 {
							continue;
						}
						for _,s := range k.Samples {
							fmt.Println("\t\t", s.Time, s.Power)
						}
					}
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}
}
