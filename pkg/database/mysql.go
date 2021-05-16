package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// MysqlConn represents a connection to a MySQL DBMS.
type MysqlConn struct {
	c *sql.DB
}

// NewMysqlConn opens a new SQL Server connection from a given dsn.
func NewMysqlConn(dsn string) (*MysqlConn, error) {
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	conn := MysqlConn{dbConn}
	return &conn, nil
}

// CreateDb attempts to create a new database as specified in the operation parameter. It returns an OpOutput with the
// result of the call.
func (c *MysqlConn) CreateDb(operation Operation) OpOutput {
	inputParams := getQueryInputs(operation.Inputs)

	rows, err := c.c.Query(operation.Name, inputParams...)
	if err != nil {
		return OpOutput{Result: nil, Err: err}
	}

	var key string
	var value string
	result := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(&key, &value)
		if err != nil {
			return OpOutput{nil, err}
		}
		result[key] = value
	}

	return OpOutput{result, nil}
}

// DeleteDb attempts to delete a database instance as specified in the operation parameter. It returns an OpOutput with the
// result of the call if present.
func (c *MysqlConn) DeleteDb(operation Operation) OpOutput {
	inputParams := getQueryInputs(operation.Inputs)

	_, err := c.c.Exec(operation.Name, inputParams...)
	if err != nil {
		return OpOutput{nil, err}
	}

	return OpOutput{}
}

// Rotate attempts to rotate the credentials of a connection.
func (c *MysqlConn) Rotate(operation Operation) OpOutput {
	inputParams := getQueryInputs(operation.Inputs)

	_, err := c.c.Exec(operation.Name, inputParams...)
	if err != nil {
		return OpOutput{nil, err}
	}

	return OpOutput{}
}

// Ping returns an error if a connection cannot be established with the DBMS, else it returns nil.
func (c *MysqlConn) Ping() error {
	return c.c.Ping()
}

