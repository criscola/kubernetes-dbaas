package database

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
)

// SqlserverConn represents a connection to a SQL Server DBMS.
type SqlserverConn struct {
	c *sql.DB
}

// NewSqlserverConn opens a new SQL Server connection from a given dsn.
func NewSqlserverConn(dsn string) (*SqlserverConn, error) {
	dbConn, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}

	conn := SqlserverConn{dbConn}
	return &conn, nil
}

// CreateDb attempts to create a new database as specified in the operation parameter. It returns an OpOutput with the
// result of the call.
func (c *SqlserverConn) CreateDb(operation Operation) OpOutput {
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
// result of the call.
func (c *SqlserverConn) DeleteDb(operation Operation) OpOutput {
	inputParams := getQueryInputs(operation.Inputs)

	_, err := c.c.Exec(operation.Name, inputParams...)
	if err != nil {
		return OpOutput{nil, err}
	}

	return OpOutput{}
}

// Rotate attempts to rotate the credentials of a connection.
func (c *SqlserverConn) Rotate(operation Operation) OpOutput {
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

func (c *SqlserverConn) Ping() error {
	return c.c.Ping()
}
