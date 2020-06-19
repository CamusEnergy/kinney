package main

import (
	"fmt"
	"flag"
	"os"

	"github.com/CamusEnergy/kinney/controller/chargepoint/simulator/sim"
)

var file = flag.String("file", "data.xlsx", "input excel file")
var ver1 = flag.Bool("v", false, "dump sessions summary")
var ver2 = flag.Bool("vv", false, "dump all sessions")

func main() {
	flag.Parse()
	ev := sim.EvChargers{}
	count, err := sim.FileExLoad(file, &ev)
	if (err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}

	var verbose uint8
	if *ver1 {
		verbose |= sim.PrintSessions
	}
	if *ver2 {
		verbose |= sim.PrintAll
		verbose |= sim.PrintSessions
	}
	
	fmt.Print("Loaded ", count, " samples")
	fmt.Println()
	sim.DataPrint(&ev, verbose)
}
