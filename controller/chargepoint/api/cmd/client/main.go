// client is a reflection-based command-line interface to the ChargePoint API
// client library.  For example:
//
//     $ ./client \
//           --credentials="credentials.json" \
//           --http_log="http_log.jsonl" \
//           --url="http://example.com/chargepoint/api" \
//           --method="GetStations" \
//           --request='{"SearchQuery": {"StationGroupID": "1:23456"}}'
//
// This tool is intended to simply be illustrative and help with exploring the
// ChargePoint API.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/CamusEnergy/kinney/controller/chargepoint/api"
)

var (
	credentialsFileName = flag.String("credentials", "", "JSON file containing the API credentials (see the `credentials` struct).  [required]")
	httpLogFileName     = flag.String("http_log", "", "File to write the JSONL-formatted HTTP log to.  [required]")
	url                 = flag.String("url", "", "URL of the ChargePoint API.  [required]")
	method              = flag.String("method", "", "API method to call.  This must exactly match (case-sensitively) a method on the `client` type.  [required]")
	request             = flag.String("request", "", "Request to send, specified as a JSON-marshalled instance of the request type.  [required]")

	timeout = flag.Duration("timeout", time.Second, "Request timeout.")
	outfile = flag.String("outfile", "", "Output file.  If empty, will write to stdout.")
)

func main() {
	// Run the main logic inside of a separate function to allow using
	// `return <error>` to interrupt control flow (instead of `log.Fatal()`)
	// to make sure that all cleanup handlers (`defer` statements) are run.
	if err := mainInternal(); err != nil {
		log.Fatal(err)
	}
}

// credentials is the schema that the `--credentials` file should conform to.
// The source file will be parsed according to the `encoding/json` umarshaling
// rules, which means that it should look something like the following:
//
//     {"APIKey": "<api-key>", "APIPassword": "<api-password>"}
type credentials struct {
	APIKey      string
	APIPassword string
}

func mainInternal() error {
	flag.Parse()

	switch {
	case *credentialsFileName == "":
		return errors.New("--credentials is required")
	case *httpLogFileName == "":
		return errors.New("--http_log is required")
	case *url == "":
		return errors.New("--url is required")
	case *method == "":
		return errors.New("--method is required")
	case *request == "":
		return errors.New("--request is required")
	}

	// Parse the file containing the credentials as a JSON-marshalled
	// instance of the `credentials` struct type.
	var creds credentials
	if b, err := ioutil.ReadFile(*credentialsFileName); err != nil {
		return fmt.Errorf("error reading credentials: %w", err)
	} else if err := json.Unmarshal(b, &creds); err != nil {
		return fmt.Errorf("error parsing credentials: %w", err)
	} else if creds.APIKey == "" || creds.APIPassword == "" {
		return errors.New("Credential file must specify both `APIKey` and `APIPassword`")
	}

	// Initialize the HTTP log file.
	httpLog, err := os.Create(*httpLogFileName)
	if err != nil {
		return fmt.Errorf("error creating HTTP log file: %w", err)
	}
	defer httpLog.Close()

	// Create the API client.
	log.Printf("Connection to SOAP endpoint: %s", *url)
	c := api.NewClient(*url, creds.APIKey, creds.APIPassword, httpLog)

	// Get the API method to call.
	methodVal := reflect.ValueOf(c).MethodByName(*method)
	if !methodVal.IsValid() {
		return fmt.Errorf("unsupported method: %s", *method)
	} else if methodVal.Type().NumIn() != 2 {
		return fmt.Errorf("client method has %q input parameters, want 2", methodVal.Type().NumIn())
	} else if methodVal.Type().In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		return fmt.Errorf("client method's first input parameter is %q, want `context.Context`", methodVal.Type().In(0))
	} else if methodVal.Type().NumOut() != 2 {
		return fmt.Errorf("client method has %q output parameters, want 2", methodVal.Type().NumOut())
	}

	// Get the type of the argument to the method on `client`.  Note that
	// the first parameter (index "0") is a `context.Context`, so this
	// parameter has index "1".
	reqType := methodVal.Type().In(1)
	if reqType.Kind() != reflect.Ptr {
		return fmt.Errorf("request type has kind %q, expected 'Ptr'", reqType.Kind())
	}

	// Create a new instance of an object of that type.
	req := reflect.New(reqType.Elem()).Interface()

	// Unmarshal `--request` into `req`.
	if err := json.Unmarshal([]byte(*request), req); err != nil {
		return fmt.Errorf("error unmarshalling JSON request: %w", err)
	}
	log.Printf("Using request: %#v", req)

	// Create a call context with the timeout specified in the flag.
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Call the API method via the reflected client method.
	out := methodVal.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
	respVal, errVal := out[0], out[1]

	// Create the output writer.  If `--outfile` was not specified, then the
	// output will be written to stdout.
	fd := os.Stdout
	if *outfile != "" {
		outfile, err := os.Create(*outfile)
		if err != nil {
			return fmt.Errorf("error creating the output file: %w", err)
		}
		fd = outfile
		defer outfile.Close()
	}

	// Handle the output parameters.
	if !errVal.IsNil() {
		return fmt.Errorf("error calling API: %w", errVal.Interface().(error))
	} else if b, err := json.Marshal(respVal.Interface()); err != nil {
		return fmt.Errorf("error marshalling response as JSON: %w", err)
	} else {
		fmt.Fprintln(fd, string(b))
	}

	return nil
}
