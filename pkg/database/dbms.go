// Package database contains all the code related to the interaction with databases. It doesn't contain any application
// state. Actual connections are retained in the pool package.
package database

import (
	"bytes"
	"fmt"
	"k8s.io/apimachinery/pkg/util/json"

	//v1 "github.com/bedag/kubernetes-dbaas/api/v1"
	"text/template"
)

const (
	Sqlserver    			= "sqlserver"
	Psql         			= "psql"
	CreateMapKey 			= "create"
	DeleteMapKey 			= "delete"
	K8sMapKey    			= "k8sName"
	UserMapKey   			= "username"
	PassMapKey   			= "password"
	DbNameMapKey 			= "dbName"
	FqdnMapKey 	 			= "fqdn"
	PortMapKey	 			= "port"
	ErrorOnMissingKeyOption = "missingkey=error"
)

// Driver represents a struct responsible for executing CreateDb and DeleteDb operations on a system it supports. Drivers
// should provide a way to check their current status (i.e. whether it can accept CreateDb and DeleteDb operations at the
// moment of a Ping call
type Driver interface {
	CreateDb(operation Operation) OpOutput
	DeleteDb(operation Operation) OpOutput
	Ping() error
}

// OpValuesd represent the input values of an operation.
type OpValues struct {
	Metadata   map[string]interface{}
	Parameters map[string]string
}

// OpOutput represents the return values of an operation. If the operation generates an error, it must be set in the Err
// field. If Err is nil, the operation is assumed to be successful.
type OpOutput struct {
	Out []string // May be changed to interface{} if typing is needed
	Err error
}

// DbmsConn represents the DBMS connection. See Driver.
type DbmsConn struct {
	Driver
}

// DbmsConfig is a slice containing Dbms structs.
type DbmsConfig []Dbms

// Dbms is the instance associated with a Dbms resource. It contains the Driver responsible for the Operations executed on
// Endpoints.
type Dbms struct {
	Driver     string
	Operations map[string]Operation
	Endpoints  []Endpoint
}

// Endpoint represent the configuration of a DBMS endpoint identified by a name.
type Endpoint struct {
	Name string
	Dsn  Dsn
}

// Operation represents an operation performed on a DBMS identified by name and containing a map of inputs and a map
// of outputs.
type Operation struct {
	Name    string
	Inputs  map[string]string
	Outputs map[string]string
}

// New initializes a Dbms instance based on a map of Operation. It expects a dsn like that:
// driver://username:password@host/instance?param1=value&param2=value
//
// See the individual Driver implementations.
func New(dsn Dsn, ops map[string]Operation) (*DbmsConn, error) {
	var dbmsConn *DbmsConn

	switch dsn.GetDriver() {
	case Sqlserver:
		sqlserverConn, err := NewMssqlConn(dsn)
		if err != nil {
			return nil, err
		}
		dbmsConn = &DbmsConn{sqlserverConn}
	case Psql:
		psqlConn, err := NewPsqlConn(dsn.String())
		if err != nil {
			return nil, err
		}
		dbmsConn = &DbmsConn{psqlConn}
	default:
		return nil, fmt.Errorf("invalid dsn '%s': driver not found", dsn)
	}

	if err := dbmsConn.Ping(); err != nil {
		return nil, err
	}

	return dbmsConn, nil
}

// GetByDriverAndEndpoint gets a Dbms identified by driver and endpoint from a DbmsConfig type.
func (c DbmsConfig) GetByDriverAndEndpoint(driver, endpoint string) (Dbms, error) {
	for _, dbms := range c {
		if dbms.Driver == driver && contains(dbms.Endpoints, endpoint) {
			return dbms, nil
		}
	}
	return Dbms{}, fmt.Errorf("dbms entry not found for driver: %s, endpoint: %s", driver, endpoint)
}

// RenderOperation renders "actions" specified through the use of the Go text/template format. It renders Dbms.Input of
// the receiver. Data to be inserted is taken directly from values. See OpValues. If the rendering is successful, the
// method returns a rendered Operation, if an error is generated, it is returned along with an empty Operation struct.
// Keys which are specified but not found generate an error (i.e. no unreferenced keys are allowed).
func (d Dbms) RenderOperation(opKey string, values OpValues) (Operation, error) {
	// Get inputs
	inputs := d.Operations[opKey].Inputs
	// Transform map[string]string to a single json string
	stringInputs, err := json.Marshal(inputs)
	if err != nil {
		return Operation{}, err
	}
	// Setup the template to be rendered based on the inputs
	tmpl, err := template.New("spParam").Parse(string(stringInputs))
	if err != nil {
		return Operation{}, err
	}
	tmpl.Option(ErrorOnMissingKeyOption)
	// Create a new buffer for the rendering result
	renderedInputsBuf := bytes.NewBufferString("")
	// Render each templated value by taking the values from the OpValues struct
	err = tmpl.Execute(renderedInputsBuf, values)
	if err != nil {
		return Operation{}, err
	}
	var renderedInputs map[string]string
	err = json.Unmarshal([]byte(renderedInputsBuf.String()), &renderedInputs)
	if err != nil {
		return Operation{}, err
	}

	renderedOp := Operation{
		Name:    d.Operations[opKey].Name,
		Inputs:  renderedInputs,
		Outputs: d.Operations[opKey].Outputs,
	}

	return renderedOp, nil
}

// IsNamePresent return true if an endpoint name is not empty, else it returns false.
func (e Endpoint) IsNamePresent() bool {
	return e.Name != ""
}

// IsDsnPresent return true if an endpoint dsn is not empty, else it returns false.
func (e Endpoint) IsDsnPresent() bool {
	return e.Dsn != ""
}

// contains is a very small utility function which returns true if s has been found in list.
func contains(list []Endpoint, s string) bool {
	for _, v := range list {
		if v.Name == s {
			return true
		}
	}
	return false
}
