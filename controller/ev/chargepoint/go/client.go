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
// "Common API"
////////////////////////////////////////////////////////////////////////////////

// API Guide: "Use this call to retrieve ChargePoint NOS instances."
func (c *client) GetCPNInstances(ctx context.Context, req *schema.GetCPNInstancesRequest) (*schema.GetCPNInstancesResponse, error) {
	resp := &schema.GetCPNInstancesResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

////////////////////////////////////////////////////////////////////////////////
// "Station Management API"
////////////////////////////////////////////////////////////////////////////////

// API Guide: "Use this call to return a list of stations.  This will not return
// stations that you don't have access rights to.  For example, it will not
// return a public station unless you either own the station or have been
// granted rights by the station's owner.  ...  Up to 500 stations will be
// returned by this method."
func (c *client) GetStations(ctx context.Context, req *schema.GetStationsRequest) (*schema.GetStationsResponse, error) {
	resp := &schema.GetStationsResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// API Guide: "Use this call to retrieve custom station groups for any
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
// "Demand Management API"
////////////////////////////////////////////////////////////////////////////////

// API Guide: "Use this call to shed load for a single port on a station, both
// ports on a multi-port station or a group of stations."
func (c *client) ShedLoad(ctx context.Context, req *schema.ShedLoadRequest) (*schema.ShedLoadResponse, error) {
	resp := &schema.ShedLoadResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// API Guide: "Use this call to clear the shed state from a single station or
// group of stations."
func (c *client) ClearShedState(ctx context.Context, req *schema.ClearShedStateRequest) (*schema.ClearShedStateResponse, error) {
	resp := &schema.ClearShedStateResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// API Guide: "Use this call to retrieve the load and shed state for a single
// station or custom station group.  This method also returns the load for each
// port on a multi-port station."
func (c *client) GetLoad(ctx context.Context, req *schema.GetLoadRequest) (*schema.GetLoadResponse, error) {
	resp := &schema.GetLoadResponse{}
	if err := c.call(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
