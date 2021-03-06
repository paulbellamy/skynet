package skylib

import (
	"os"
	"time"
	"fmt"
	"rpc"
	"container/vector"
)


// RpcService is a struct that represents a remotly 
// callable function.  It is intented to be part of 
// an array or collection of RpcServices.  It contains
// a member "Provides" which is the name of the service the
// remote call provides, and a Client pointer which is a pointer
// to an RPC client connected to this service.
type RpcService struct {
	Provides string
}


func (r *RpcService) parseError(err string) {
	panic(&Error{err, r.Provides})
}


// A Generic struct to represent any service in the SkyNet system.
type Service struct {
	IPAddress string
	Name      string
	Port      int
	Provides  string
}


// A HeartbeatRequest is the struct that is sent for ping checks.
type HeartbeatRequest struct {
	Timestamp int64
}

// HeartbeatResponse is the struct that is returned on a ping check.
type HeartbeatResponse struct {
	Timestamp int64
	Ok        bool
}

// HealthCheckRequest is the struct that is sent on a more advanced heartbeat request.
type HealthCheckRequest struct {
	Timestamp int64
}


// HealthCheckResponse is the struct that is sent back to the HealthCheckRequest-er
type HealthCheckResponse struct {
	Timestamp int64
	Load      float64
}


// A Route represents an ordered list of RPC calls that should be made for 
// a request.  Routes are versioned and named.  Names should correspond to 
// Service names- which makes me wonder if the route should be stored right there 
// in the Service struct??
type Route struct {
	Name        string
	RouteList   *vector.Vector
	Revision    int64
	LastUpdated int64
}

// The struct that is stored in the Route
// Async delineates whether it's ok to call this and not
// care about the response.
// OkToRetry delineates whether it's ok to call this service
// more than once.
type RpcCall struct {
	Service   string
	Async     bool
	OkToRetry bool
	ErrOnFail bool
}

// Parent struct for the configuration
type NetworkServers struct {
	Services []*Service
}

type ServerConfig interface {
	Equal(that interface{}) bool
}

// Exported RPC method for the health check
func (hc *Service) Ping(hr *HeartbeatRequest, resp *HeartbeatResponse) (err os.Error) {

	resp.Timestamp = time.Seconds()

	return nil
}

// Exported RPC method for the advanced health check
func (hc *Service) PingAdvanced(hr *HealthCheckRequest, resp *HealthCheckResponse) (err os.Error) {

	resp.Timestamp = time.Seconds()
	resp.Load = 0.1 //todo
	return nil
}

// Method to register the heartbeat of each skynet
// client with the healthcheck exporter.
func RegisterHeartbeat() {
	r := NewService("Service.Ping")
	rpc.Register(r)
}


func (r *Service) Equal(that *Service) bool {
	var b bool
	b = false
	if r.Name != that.Name {
		return b
	}
	if r.IPAddress != that.IPAddress {
		return b
	}
	if r.Port != that.Port {
		return b
	}
	if r.Provides != that.Provides {
		return b
	}
	b = true
	return b
}

// Utility function to return a new Service struct
// pre-populated with the data on the command line.
func NewService(provides string) *Service {
	r := &Service{
		Name:      *Name,
		Port:      *Port,
		IPAddress: *BindIP,
		Provides:  provides,
	}

	return r
}


type Error struct {
	Msg     string
	Service string
}

func (e *Error) String() string { return fmt.Sprintf("Service %s had error: %s", e.Service, e.Msg) }

func NewError(msg string, service string) (err *Error) {
	err = &Error{Msg: msg, Service: service}
	return
}
