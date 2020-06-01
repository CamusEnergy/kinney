package chargepoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSOAPCall(t *testing.T) {
	// Create an HTTP server for testing.  This binds to a random port on
	// the loopback interface, the address of which can be retrieved as
	// `server.URL`.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if b, err := MarshalEnvelope("respHeader", "respBody"); err != nil {
			// There is no simple way to propogate this error to the
			// calling code, so just panic.
			panic(err)
		} else if _, err := w.Write(b); err != nil {
			panic(err)
		}
	}))
	defer server.Close()

	// Issue a SOAP request against the HTTP server created above, and make
	// sure that the response is umarshalled correctly.
	var respHeader, respBody string
	var httpLog bytes.Buffer
	if err := soapCall(context.Background(), &http.Client{}, server.URL, "reqHeader", "reqBody", &respHeader, &respBody, &httpLog); err != nil {
		t.Errorf("soapCall() = %q; want nil", err)
	} else if respHeader != "respHeader" || respBody != "respBody" {
		t.Errorf("respHeader = %q, respBody = %q; want %q and %q", respHeader, respBody, "respHeader", "respBody")
	}

	var logEntry httpLogEntry
	if err := json.Unmarshal(httpLog.Bytes(), &logEntry); err != nil {
		t.Errorf("json.Unmarshal(%q) = %q; want nil", httpLog.Bytes(), err)
	}

	expectedReqBody, _ := MarshalEnvelope("reqHeader", "reqBody")
	expectedRespBody, _ := MarshalEnvelope("respHeader", "respBody")

	// Validate the log entry that was written to the HTTP log strream.
	expectedLogEntry := httpLogEntry{
		RequestTimestamp:   time.Now(),
		RequestMethod:      "POST",
		RequestURL:         server.URL,
		RequestHeaders:     http.Header{"Content-Type": {"text/xml; charset=utf-8"}},
		RequestBody:        expectedReqBody,
		ResponseTimestamp:  time.Now(),
		ResponseStatusCode: 200,
		ResponseHeaders:    nil,
		ResponseBody:       expectedRespBody,
		Err:                nil,
	}
	if diff := cmp.Diff(&expectedLogEntry, &logEntry, cmp.Options{cmpopts.EquateApproxTime(time.Second)}); diff != "" {
		t.Errorf("log entry mismatch (-want +got):\n%s", diff)
	}
}

const envelopeXML = `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/"><Header><string>foo</string></Header><Body><string>bar</string></Body></Envelope>`

func TestMarshalUnmarshal(t *testing.T) {
	header := "foo"
	body := "bar"

	b, err := MarshalEnvelope(&header, &body)
	if err != nil {
		t.Errorf("marshalEnvelope(%#v, %#v) = %q; want nil", &header, &body, err)
	} else if diff := cmp.Diff(envelopeXML, string(b)); diff != "" {
		t.Errorf("marshalEnvelope(%#v, %#v) mismatch (-want +got):\n%s", &header, &body, diff)
	}

	var parsedHeader, parsedBody string
	if err := UnmarshalEnvelope(b, &parsedHeader, &parsedBody); err != nil {
		t.Errorf("unmarshalEnvelope(%q) = %q; want nil", string(b), err)
	} else if diff := cmp.Diff(header, parsedHeader); diff != "" {
		t.Errorf("unmarshalEnvelope(%q) mismatch (-want +got):\n%s", string(b), diff)
	} else if diff := cmp.Diff(body, parsedBody); diff != "" {
		t.Errorf("unmarshalEnvelope(%q) mismatch (-want +got):\n%s", string(b), diff)
	}
}
