// Package chargepoint provides a native client library wrapping the
// "ChargePoint Web Services API Version 5.0".
//
// Note that any comments quoted from the API Guide are indicated as such, and
// will be quoted verbatim.
//
// TODO(james): Add more detail to the package comment, including examples of
// how to use the library.
package chargepoint

import (
	"context"
	"io"
	"net/http"

	"github.com/CamusEnergy/kinney/controller/ev/chargepoint/go/schema"
)

// client provides an interface for communicating with a ChargePoint API server.
type client struct {
	url            string
	securityHeader *securityHeader
	httpClient     *http.Client
	httpLogWriter  io.Writer
}

// NewCliant returns a new client instance that communicates with the given SOAP
// Endpoint (`url`) with the given ChargePoint API credentials.  A record of
// every HTTP exchange will be written to `httpLogWriter`, in JSONL format.
//
// Note that the ChargePoint Web Services API is a SOAP v1.1 endpoint protected
// by a WS-Security "UsernameToken".  The username is the API license key, and
// the password is the API password.
func NewClient(url, apiKey, apiPassword string, httpLogWriter io.Writer) *client {
	return &client{
		url:            url,
		securityHeader: newSecurityHeader(apiKey, apiPassword),
		httpClient:     http.DefaultClient,
		httpLogWriter:  httpLogWriter,
	}
}

// call simply wraps `soapCall`, filling in the common parameters from the
// client's configuration.
func (c *client) call(ctx context.Context, req, resp interface{}) error {
	return soapCall(ctx, c.httpClient, c.url, c.securityHeader, req, nil, resp, c.httpLogWriter)
}

////////////////////////////////////////////////////////////////////////////////
// API Guide (§ 4): "Common API"
////////////////////////////////////////////////////////////////////////////////

// API Guide (§ 4.1): "Use this call to retrieve ChargePoint NOS instances."
//
// TODO(james): Figure out what "NOS" means, as it is not documented in the API
// Guide.
func (c *client) GetCPNInstances(ctx context.Context, req *schema.GetCPNInstancesRequest) (*schema.GetCPNInstancesResponse, error) {
	resp := &schema.GetCPNInstancesResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

////////////////////////////////////////////////////////////////////////////////
// API Guide (§ 8): "Station Management API"
////////////////////////////////////////////////////////////////////////////////

// API Guide (§ 8.1): "Use this call to return a list of stations.  This will
// not return stations that you don't have access rights to.  For example, it
// will not return a public station unless you either own the station or have
// been granted rights by the station's owner."
//
// API Guide (§ 8.1.1): "Up to 500 stations will be returned by this method."
func (c *client) GetStations(ctx context.Context, req *schema.GetStationsRequest) (*schema.GetStationsResponse, error) {
	resp := &schema.GetStationsResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// API Guide (§ 8.3): "Use this call to retrieve custom station groups for any
// organization.  It returns an array of groups for a given organization and
// lists the stations included in each group."
func (c *client) GetStationGroups(ctx context.Context, req *schema.GetStationGroupsRequest) (*schema.GetStationGroupsResponse, error) {
	resp := &schema.GetStationGroupsResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

////////////////////////////////////////////////////////////////////////////////
// API Guide (§ 6): "Demand Management API"
////////////////////////////////////////////////////////////////////////////////

// API Guide (§ 6.1): "Use this call to shed load for a single port on a
// station, both ports on a multi-port station or a group of stations."
func (c *client) ShedLoad(ctx context.Context, req *schema.ShedLoadRequest) (*schema.ShedLoadResponse, error) {
	resp := &schema.ShedLoadResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// API Guide (§ 6.2): "Use this call to clear the shed state from a single
// station or group of stations."
func (c *client) ClearShedState(ctx context.Context, req *schema.ClearShedStateRequest) (*schema.ClearShedStateResponse, error) {
	resp := &schema.ClearShedStateResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// API Guide (§ 6.3): "Use this call to retrieve the load and shed state for a
// single station or custom station group.  This method also returns the load
// for each port on a multi-port station."
func (c *client) GetLoad(ctx context.Context, req *schema.GetLoadRequest) (*schema.GetLoadResponse, error) {
	resp := &schema.GetLoadResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
