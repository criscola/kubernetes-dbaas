package database

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
)

// MssqlConn represents a connection to a SQL Server DBMS.
type MssqlConn struct {
	c *sql.DB
}

// NewMssqlConn constructs a new SQL Server connection from a given dsn.
func NewMssqlConn(dsn Dsn) (*MssqlConn, error) {
	dbConn, err := sql.Open("sqlserver", dsn.String())
	if err != nil {
		return nil, err
	}

	conn := MssqlConn{dbConn}
	return &conn, nil
}

// CreateDb attempts to create a new database as specified in the operation parameter. It returns an OpOutput with the
// result of the call.
func (c *MssqlConn) CreateDb(operation Operation) OpOutput {
	var username string
	var password string
	var dbName string
	var fqdn string
	var port string

	var inputParams []interface{}
	for k, v := range operation.Inputs {
		inputParams = append(inputParams, sql.Named(k, v))
	}
	inputParams = append(inputParams,
		sql.Named(operation.Outputs[UserMapKey], sql.Out{Dest: &username}),
		sql.Named(operation.Outputs[PassMapKey], sql.Out{Dest: &password}),
		sql.Named(operation.Outputs[DbNameMapKey], sql.Out{Dest: &dbName}),
		sql.Named(operation.Outputs[FqdnMapKey], sql.Out{Dest: &fqdn}),
		sql.Named(operation.Outputs[PortMapKey], sql.Out{Dest: &port}))

	_, err := c.c.Exec(operation.Name, inputParams...)
	if err != nil {
		return OpOutput{nil, err}
	}

	return OpOutput{[]string{username, password, dbName, fqdn, port}, nil}
}

// DeleteDb attemps to delete a database instance as specified in the operation parameter. It returns an OpOutput with the
// result of the call.
func (c *MssqlConn) DeleteDb(operation Operation) OpOutput {
	var inputParams []interface{}
	for k, v := range operation.Inputs {
		inputParams = append(inputParams, sql.Named(k, v))
	}

	_, err := c.c.Exec(operation.Name, inputParams...)
	if err != nil {
		return OpOutput{nil, err}
	}

	return OpOutput{nil, nil}
}

// Ping returns an error if a connection cannot be established with the DBMS, else it returns nil.
func (c *MssqlConn) Ping() error {
	return c.c.Ping()
}
