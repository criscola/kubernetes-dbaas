package database

import (
	"fmt"
)

const (
	Sqlserver    = "sqlserver"
	Psql         = "psql"
	CreateMapKey = "create"
	DeleteMapKey = "delete"
	K8sMapKey    = "k8sName"
	UserMapKey   = "username"
	PassMapKey   = "password"
	DbNameMapKey = "dbName"
)

// Driver represents a struct responsible for executing CreateDb and DeleteDb operations on a system it supports. Drivers
// should provide a way to check their current status (i.e. whether it can accept CreateDb and DeleteDb operations at the
// moment of a Ping call
type Driver interface {
	CreateDb(name string) OpOutput
	DeleteDb(name string) OpOutput
	Ping() error
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

type Endpoint struct {
	Name string
	Dsn  Dsn
}

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
		sqlserverConn, err := NewMssqlConn(dsn, ops)
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

func (e Endpoint) IsNamePresent() bool {
	return e.Name != ""
}

func (e Endpoint) IsDsnPresent() bool {
	return e.Dsn != ""
}
