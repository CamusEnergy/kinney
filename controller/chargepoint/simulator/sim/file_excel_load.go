package sim

import (
	"fmt"
	"strconv"
	"time"
	"strings"
	
	ex "github.com/360EntSecGroup-Skylar/excelize/v2"
)

const vmGroup = "1:19400001"
const vmName = "VMware"

const geoLat = "42.63228390329662"
const geoLong = "23.378210952553545"

const cpnName = "Virtual"
const cpnDesc = "EV charger simulator"

const threshold = 60 * time.Minute

func FileExLoad(file *string, e *EvChargers) (int, error) {
	var lastSession *chargeSession
	var lastTime time.Time
	var samples int

	e.Vehicles = make(map[string]int)
	e.StOrgs = make(map[string]evOrg)
	e.StCPNs = make(map[string]evCPN)
	
	f, err := ex.OpenFile(*file)
	if err != nil {
		return samples, err
	}
	
	e.StCPNs["1"] = evCPN{
				Name:cpnName,
				Desc:cpnDesc,
			}
	e.StOrgs[vmGroup] = evOrg{
				Name:vmName,
				StGrpoups: make(map[string]evStGroup),
			}
	
	sheets := f.GetSheetMap() 
	for _, name := range sheets {
		rows, err := f.GetRows(name)
		if err != nil {
				fmt.Println(err)
				continue
		}

		for _, col := range rows[1:] {
			var diff time.Duration
			var ev chargeSession
			var s chargeSample
			var i int
			
			if len(col) < 4 || len(col) > 6 {
				continue
			}
/*
			col[0]: Timestamp,		"1583550049.76136"
			col[1]: VehicleID,		"HNA3BC734CE51"
			col[2]: Charge,			"5.661"
			col[3]: Full-Port-ID	"238421*1:569591*2"
*/
			for i = 0; i < 4; i++ {
				if col[i] == "" {
					break
				}
			}
			if i < 4 {
				continue
			}
			ev.VehicleId = strings.TrimSpace(col[1])
			ev.PortId = strings.TrimSpace(col[3])

			times := strings.SplitN(col[0], ".", -1)
			t1, err := strconv.ParseInt(times[0], 10, 64)
			if err != nil {
				continue
			}
			t2, err := strconv.ParseInt(times[1], 10, 64)
			if err != nil {
				continue
			}
			p, err := strconv.ParseFloat(col[2], 32)
			if err != nil {
				continue
			}
			
			var gr *evStGroup
			ids := strings.SplitN(ev.PortId, "*", -1)
			if ids[0] != "" {
				if g, ok := e.StOrgs[vmGroup].StGrpoups[ids[0]]; !ok {
					g = evStGroup{Stations: make(map[string]evStation)}
					e.StOrgs[vmGroup].StGrpoups[ids[0]] = g;
					gr = &g
				} else {
					gr = &g
				}
			} else {
				continue
			}

			var st *evStation
			if ids[1] != "" {
				if g, ok := gr.Stations[ids[1]]; !ok {
					g = evStation{ 
							Geo: evGeo{Lat:geoLat, Long:geoLong} ,
							Ports: make(map[string]evPort),
					}
					gr.Stations[ids[1]] = g
					st = &g
				} else {
					st = &g
				}
				
			} else {
				continue
			}
			
			if ids[2] != "" {
				if _, ok := st.Ports[ids[2]]; !ok {
					st.Ports[ids[2]] = evPort{Charge: make([]*chargeSession, 0)}
				}
			} else {
				continue
			}

			s.Time = time.Unix(t1, t2)
			s.Power = float32(p)
			diff = s.Time.Sub(lastTime)
			if lastSession != nil &&
			   lastSession.VehicleId == ev.VehicleId &&
			   lastSession.PortId == ev.PortId &&
			   diff < threshold {
				   	lastSession.Samples = append(lastSession.Samples, s)
				   	lastTime = s.Time
			} else {
					ev.Samples = append(ev.Samples, s)
					if prt,ok := st.Ports[ids[2]]; ok {
						prt.Charge = append(prt.Charge, &ev)
						st.Ports[ids[2]] = prt
					}
					lastSession = &ev
					lastTime = s.Time
			}
			e.Vehicles[ev.VehicleId]++
			samples++
		}
	}
	
	return samples, nil
}

