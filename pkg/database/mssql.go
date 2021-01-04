package database

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
)

type MssqlConn struct {
	c          *sql.DB
	operations map[string]Operation
}

func NewMssqlConn(dsn Dsn, ops map[string]Operation) (*MssqlConn, error) {
	dbConn, err := sql.Open("sqlserver", dsn.String())
	if err != nil {
		return nil, err
	}

	conn := MssqlConn{dbConn, ops}
	return &conn, nil
}

func (c *MssqlConn) CreateDb(name string) QueryOutput {
	var username string
	var password string
	var dbName string

	operation := c.operations[CreateMapKey]

	_, err := c.c.Exec(operation.Name,
		sql.Named(operation.Inputs[K8sMapKey], name),
		sql.Named(operation.Outputs[UserMapKey], sql.Out{Dest: &username}),
		sql.Named(operation.Outputs[PassMapKey], sql.Out{Dest: &password}),
		sql.Named(operation.Outputs[DbNameMapKey], sql.Out{Dest: &dbName}),
	)
	if err != nil {
		return QueryOutput{nil, err}
	}
	return QueryOutput{[]string{username, password, dbName}, nil}
}

func (c *MssqlConn) DeleteDb(name string) QueryOutput {
	operation := c.operations[DeleteMapKey]

	_, err := c.c.Exec(operation.Name,
		sql.Named(operation.Inputs[K8sMapKey], name),
	)
	if err != nil {
		return QueryOutput{nil, err}
	}
	return QueryOutput{nil, nil}
}

func (c *MssqlConn) Ping() error {
	return c.c.Ping()
}
