package database

import (
	"fmt"
)

// Driver represents a struct responsible for executing CreateDb and DeleteDb operations on a system it supports. Drivers
// should provide a way to check their current status (i.e. whether it can accept CreateDb and DeleteDb operations at the
// moment of a Ping call
type Driver interface {
	CreateDb(name string) QueryOutput
	DeleteDb(name string) QueryOutput
	Ping() error
}

// QueryOutput represents the return of a
type QueryOutput struct {
	Out []string // May be changed to interface{} if typing is needed
	Err error
}

type DbmsConn struct {
	Driver
}

type DbmsConfig []Dbms

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

const (
	Sqlserver    = "sqlserver"
	Psql         = "psql"
	CreateMapKey = "create"
	DeleteMapKey = "delete"
	K8sMapKey    = "overrideName"
	UserMapKey   = "username"
	PassMapKey   = "password"
	DbNameMapKey = "dbname"
)

// Expects a dsn like that sqlserver://username:password@host/instance?param1=value&param2=value
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
