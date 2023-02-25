// Code generated by go-swagger; DO NOT EDIT.

package waypoint

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewWaypointGetWorkspaceParams creates a new WaypointGetWorkspaceParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWaypointGetWorkspaceParams() *WaypointGetWorkspaceParams {
	return &WaypointGetWorkspaceParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewWaypointGetWorkspaceParamsWithTimeout creates a new WaypointGetWorkspaceParams object
// with the ability to set a timeout on a request.
func NewWaypointGetWorkspaceParamsWithTimeout(timeout time.Duration) *WaypointGetWorkspaceParams {
	return &WaypointGetWorkspaceParams{
		timeout: timeout,
	}
}

// NewWaypointGetWorkspaceParamsWithContext creates a new WaypointGetWorkspaceParams object
// with the ability to set a context for a request.
func NewWaypointGetWorkspaceParamsWithContext(ctx context.Context) *WaypointGetWorkspaceParams {
	return &WaypointGetWorkspaceParams{
		Context: ctx,
	}
}

// NewWaypointGetWorkspaceParamsWithHTTPClient creates a new WaypointGetWorkspaceParams object
// with the ability to set a custom HTTPClient for a request.
func NewWaypointGetWorkspaceParamsWithHTTPClient(client *http.Client) *WaypointGetWorkspaceParams {
	return &WaypointGetWorkspaceParams{
		HTTPClient: client,
	}
}

/*
WaypointGetWorkspaceParams contains all the parameters to send to the API endpoint

	for the waypoint get workspace operation.

	Typically these are written to a http.Request.
*/
type WaypointGetWorkspaceParams struct {

	// WorkspaceWorkspace.
	WorkspaceWorkspace string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the waypoint get workspace params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetWorkspaceParams) WithDefaults() *WaypointGetWorkspaceParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the waypoint get workspace params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetWorkspaceParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the waypoint get workspace params
func (o *WaypointGetWorkspaceParams) WithTimeout(timeout time.Duration) *WaypointGetWorkspaceParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the waypoint get workspace params
func (o *WaypointGetWorkspaceParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the waypoint get workspace params
func (o *WaypointGetWorkspaceParams) WithContext(ctx context.Context) *WaypointGetWorkspaceParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the waypoint get workspace params
func (o *WaypointGetWorkspaceParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the waypoint get workspace params
func (o *WaypointGetWorkspaceParams) WithHTTPClient(client *http.Client) *WaypointGetWorkspaceParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the waypoint get workspace params
func (o *WaypointGetWorkspaceParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithWorkspaceWorkspace adds the workspaceWorkspace to the waypoint get workspace params
func (o *WaypointGetWorkspaceParams) WithWorkspaceWorkspace(workspaceWorkspace string) *WaypointGetWorkspaceParams {
	o.SetWorkspaceWorkspace(workspaceWorkspace)
	return o
}

// SetWorkspaceWorkspace adds the workspaceWorkspace to the waypoint get workspace params
func (o *WaypointGetWorkspaceParams) SetWorkspaceWorkspace(workspaceWorkspace string) {
	o.WorkspaceWorkspace = workspaceWorkspace
}

// WriteToRequest writes these params to a swagger request
func (o *WaypointGetWorkspaceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param workspace.workspace
	if err := r.SetPathParam("workspace.workspace", o.WorkspaceWorkspace); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}