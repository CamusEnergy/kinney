package main

import (
	"fmt"
	"flag"
	"os"
	"net/http"
	
	"github.com/CamusEnergy/kinney/controller/chargepoint/api"	
	"github.com/CamusEnergy/kinney/controller/chargepoint/simulator/sim"
	"github.com/CamusEnergy/kinney/controller/chargepoint/simulator/simapi"
)

var file = flag.String("file", "", "input excel file")
var ver1 = flag.Bool("v", false, "dump sessions summary")
var ver2 = flag.Bool("vv", false, "dump all sessions")
var port = flag.Int("port", 8080, "Web server port")
var crt = flag.String("crt", "", "TLS certificate file")
var key = flag.String("key", "", "TLS key file")
var url = flag.String("url", "/webservices/chargepoint/services/5.0", "API endpoint")

func serverRun(ev *sim.EvChargers, port int, url string, crt, key string) error {
	var err error
	
	mux := http.NewServeMux()
	mux.Handle(url, api.NewHandler(simapi.SimulatorServer{Ev: ev}))
	p := fmt.Sprintf(":%d", port)
	if crt != "" && key != "" {
		err = http.ListenAndServeTLS(p, crt, key, mux)
	} else {
		err = http.ListenAndServe(p, mux)
	}	
	return err
}

func main() {
	flag.Parse()
	ev := sim.EvChargers{}
	var count int
	var err error
	
	if *file != "" {
		count, err = sim.FileExLoad(file, &ev)
		if (err != nil) {
			fmt.Println(err)
			os.Exit(1)
		}
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
	
	err = serverRun(&ev, *port, *url, *crt, *key)
	
	fmt.Println(err)
}
