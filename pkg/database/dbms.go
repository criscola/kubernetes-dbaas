package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

const (
	Sqlserver               = "sqlserver"
	Psql                    = "psql"
	CreateMapKey            = "create"
	DeleteMapKey            = "delete"
	OperationsConfigKey     = "operations"
	K8sMapKey               = "k8sName"
	UserMapKey              = "username"
	PassMapKey              = "password"
	DbNameMapKey            = "dbName"
	FqdnMapKey              = "fqdn"
	PortMapKey              = "port"
	ErrorOnMissingKeyOption = "missingkey=error"
	DbmsConfigKey           = "dbms"
)

// Driver represents a struct responsible for executing CreateDb and DeleteDb operations on a system it supports. Drivers
// should provide a way to check their current status (i.e. whether it can accept CreateDb and DeleteDb operations at the
// moment of a Ping call
type Driver interface {
	CreateDb(operation Operation) OpOutput
	DeleteDb(operation Operation) OpOutput
	Ping() error
}

// DbmsConn represents the DBMS connection. See Driver.
type DbmsConn struct {
	Driver
}

// +kubebuilder:object:generate=true
// Operation represents an operation performed on a DBMS identified by name and containing a map of inputs and a map
// of outputs.
type Operation struct {
	Name    string            `json:"name,omitempty"`
	Inputs  map[string]string `json:"inputs,omitempty"`
	Outputs map[string]string `json:"outputs,omitempty"`
}

// OpOutput represents the return values of an operation. If the operation generates an error, it must be set in the Err
// field. If Err is nil, the operation is assumed to be successful.
type OpOutput struct {
	Out []string // May be changed to interface{} if typing is needed
	Err error
}

// OpValues represent the input values of an operation.
type OpValues struct {
	Metadata   map[string]interface{}
	Parameters map[string]string
}

// +kubebuilder:object:generate=true
// Dbms is the instance associated with a Dbms resource. It contains the Driver responsible for the Operations executed on
// Endpoints.
type Dbms struct {
	DatabaseClassName string     `json:"databaseClassName"`
	Endpoints         []Endpoint `json:"endpoints"`
}

// +kubebuilder:object:generate=true
// DbmsList is a slice containing Dbms structs.
type DbmsList []Dbms

// +kubebuilder:object:generate=true
// +kubebuilder:kubebuilder:validation:MinItems=1
// Endpoint represent the configuration of a DBMS endpoint identified by a name.
type Endpoint struct {
	Name string `json:"name"`
	Dsn  Dsn    `json:"dsn"`
}

// New initializes a Dbms instance based on a map of Operation. It expects a dsn like that:
// driver://username:password@host/instance?param1=value&param2=value
//
// See the individual Driver implementations.
func New(dsn Dsn) (*DbmsConn, error) {
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

// RenderOperation renders "actions" specified through the use of the Go text/template format. It renders Input of
// the receiver. Data to be inserted is taken directly from values. See OpValues. If the rendering is successful, the
// method returns a rendered Operation, if an error is generated, it is returned along with an empty Operation struct.
// Keys which are specified but not found generate an error (i.e. no unreferenced keys are allowed).
func (op *Operation) RenderOperation(values OpValues) (Operation, error) {
	// Get inputs
	inputs := op.Inputs
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
		Name:    op.Name,
		Inputs:  renderedInputs,
		Outputs: op.Outputs,
	}

	return renderedOp, nil
}

func (c DbmsList) GetDatabaseClassNameByEndpointName(endpointName string) (string, error) {
	for _, dbms := range c {
		if contains(dbms.Endpoints, endpointName) {
			return dbms.DatabaseClassName, nil
		}
	}
	return "", fmt.Errorf("could not find any DatabaseClass for endpoint '%s'", endpointName)
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
