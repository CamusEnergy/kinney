package main

import (
	"errors"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/CamusEnergy/kinney/controller/chargepoint/simulator"

)

var (
	sgID = flag.Int("sgID", 1, "sgID is the station group ID, defaults to 1")
	address = flag.String("address", "Somewhere, SomeState, SomeCountry", "address is the address of the EV charging facility")
	capacity = flag.Float64("capacity", 8.0, "Maximum charge capacity of the ports")
	numStations = flag.Int("numStations", 1, "numStations is the number of stations in this station group, defaults to 1")
    simPort = flag.Int("simPort", 8089, "port on which to run simulator service, defaults to 8089")
)

// TODO -- flesh the methods out using the Chargepoint structs, leverage Tzvetomir's code for all the SOAP stuff
// Alternately  work at a layer above SOAP, taking and returning Go structs, past SOAP input/output digestion
// Add support simulating different vehicle arrival and departure patterns based on location (office|home|mall) 
func main() {
	// runs the main simulator logic
	fmt.Println("Simulator: collecting parameters!")
	if err := mainInternal(); err != nil {
		log.Fatal(err)
	}
	var chargeFacility = simulator.NewChargeFacility(*sgID, *numStations, float32(8.0), *address)
	fmt.Printf("Created a Charge Facility: %#v\n", chargeFacility)
	// TODO handle the API specific operations
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TODO - integrate with the SOAP handling code")
	})

	fmt.Printf("Simulator: running on port %d\n", *simPort)
	log.Fatal(http.ListenAndServe(":8089", nil))

}

func mainInternal() error {
	flag.Parse()
	switch {
	case 0 == *sgID:
		return errors.New("non-zero --sgID is required")
	case 0 == *numStations:
		return errors.New("no-zero --numStations is required")
	}
	fmt.Printf("sgId = %d, numStations = %d numPorts = %d\n", *sgID, *numStations, *numPorts)
	return nil
}


