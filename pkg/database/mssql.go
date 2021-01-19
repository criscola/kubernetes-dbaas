package database

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
)

type MssqlConn struct {
	c *sql.DB
}

func NewMssqlConn(dsn Dsn) (*MssqlConn, error) {
	dbConn, err := sql.Open("sqlserver", dsn.String())
	if err != nil {
		return nil, err
	}

	conn := MssqlConn{dbConn}
	return &conn, nil
}

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

func (c *MssqlConn) Ping() error {
	return c.c.Ping()
}
