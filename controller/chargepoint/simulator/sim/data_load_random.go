// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

package sim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sethvargo/go-diceware/diceware"
)

type randomParams struct {
	MaxCPNs           int `json:"maxCPNs"`
	MaxFacilities     int `json:"maxFacilities"`
	MaxChargeGroups   int `json:"maxChargeGroups"`
	MaxChargeStations int `json:"maxChargeStations"`
	MaxChargePorts    int `json:"maxChargePorts"`
	MaxVehicleBattery int `json:"maxVehicleBattery"`
	NumCPNs           int `json:"CPNs"`
	NumFacilities     int `json:"Facilities"`
	NumChargeGroups   int `json:"ChargeGroups"`
	NumChargeStations int `json:"ChargeStations"`
	NumChargePorts    int `json:"ChargePorts"`
	PortLoad          int `json:"PortLoad"`
	RandSeed          int `json:"RandomSeed"`
}

const (
	defCPNs           = 1
	defFacilities     = 4
	defChargeGroups   = 4
	defChargeStations = 8
	defChargePorts    = 2
	defVehicleBattery = 80  // KW
	defPortCapacity   = 8.0 // KW

	// default probability of a port to have a charging session, in %
	defPortLoad = 50
)

func randParam(min int, max int) int {
	if max <= min {
		return min
	}
	return min + rand.Intn(max-min)
}

func randString(words int, sep string) *string {
	l, err := diceware.Generate(words)
	str := ""
	if err != nil {
		return &str
	}
	str = strings.Join(l, sep)
	return &str
}

func genRandom(param *randomParams, e *EVChargers) {

	nCPNs := param.NumCPNs
	nFacilities := param.NumFacilities
	nChargeGroups := param.NumChargeGroups
	nChargeStations := param.NumChargeStations
	nChargePorts := param.NumChargePorts

	if param.RandSeed > 0 {
		rand.Seed(int64(param.RandSeed))
	} else {
		rand.Seed(int64(time.Now().Nanosecond()))
	}

	if param.MaxCPNs > 0 {
		nCPNs = randParam(1, param.MaxCPNs)
	}
	if param.MaxFacilities > 0 {
		nFacilities = randParam(1, param.MaxFacilities)
	}
	if param.MaxChargeGroups > 0 {
		nChargeGroups = randParam(1, param.MaxChargeGroups)
	}
	if param.MaxChargeStations > 0 {
		nChargeStations = randParam(1, param.MaxChargeStations)
	}
	if param.MaxChargePorts > 0 {
		nChargePorts = randParam(1, param.MaxChargeStations)
	}

	// Generate up to param.MaxCPNs random networks
	for c := 0; c <= nCPNs; c++ {
		str := strconv.FormatInt(int64(c), 10)
		e.getCPNetwork(&str, randString(1, ""), randString(5, ","))
	}

	// Generate up to param.MaxFacilities random facilities
	for f := 0; f < nFacilities; f++ {
		facility := fmt.Sprintf("1:%d", randParam(10000000, 99999999))
		e.getChargeFacility(&facility, randString(1, ""))

		// Generate up to param.MaxChargeGroups random charge groups in each facilitie
		for g := 0; g < nChargeGroups; g++ {
			group := fmt.Sprintf("%d", randParam(100000, 999999))
			e.getChargeGroup(&facility, &group, randString(1, ""),
				&getLoadRandom{
					portLoad:          param.PortLoad,
					maxVehicleBattery: param.MaxVehicleBattery,
				})

			// Generate up to param.MaxChargeStations random charge stations in each group
			for s := 0; s < nChargeStations; s++ {
				station := fmt.Sprintf("1:%d", randParam(100000, 999999))
				e.getChargeStation(&facility, &group, &station, randString(1, ""),
					randString(5, ","),
					&locGeo{lat: geoDefLat, long: geoDefLong})

				// Generate up to param.MaxChargePorts random charge ports in each station
				for p := 0; p < nChargePorts; p++ {
					str := strconv.FormatInt(int64(p), 10)
					port := e.getChargePort(&facility, &group, &station, &str)
					port.max_capacity = defPortCapacity
					port.current_capacity = defPortCapacity
				}
			}
		}
	}

}

func DataLoadRandom(config *string, e *EVChargers) error {
	var err error
	var file *os.File
	var data []byte
	params := randomParams{
		NumCPNs:           defCPNs,
		NumFacilities:     defFacilities,
		NumChargeGroups:   defChargeGroups,
		NumChargeStations: defChargeStations,
		NumChargePorts:    defChargePorts,
		MaxVehicleBattery: defVehicleBattery,
		PortLoad:          defPortLoad,
	}

	if file, err = os.Open(*config); err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	if data, err = ioutil.ReadAll(file); err != nil {
		fmt.Println(err)
		return err
	}

	if err = json.Unmarshal(data, &params); err != nil {
		fmt.Println(err)
		return err
	}

	genRandom(&params, e)
	return nil
}
