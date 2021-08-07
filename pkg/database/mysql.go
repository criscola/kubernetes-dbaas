package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
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
	sp, err := GetMysqlOpQuery(operation)
	if err != nil {
		return OpOutput{
			Result: nil,
			Err:    err,
		}
	}

	rows, err := c.c.Query(sp)
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
	sp, err := GetMysqlOpQuery(operation)
	if err != nil {
		return OpOutput{
			Result: nil,
			Err:    err,
		}
	}
	_, err = c.c.Exec(sp)
	if err != nil {
		return OpOutput{nil, err}
	}

	return OpOutput{}
}

// Rotate attempts to rotate the credentials of a connection.
func (c *MysqlConn) Rotate(operation Operation) OpOutput {
	sp, err := GetMysqlOpQuery(operation)
	if err != nil {
		return OpOutput{
			Result: nil,
			Err:    err,
		}
	}

	rows, err := c.c.Query(sp)
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

// Ping returns an error if a connection cannot be established with the DBMS, else it returns nil.
func (c *MysqlConn) Ping() error {
	return c.c.Ping()
}

// GetMysqlOpQuery constructs a CALL query from the specified operation. Keys of operation.Inputs must be integers, they
// are converted from string to int and then used to sort the parameters in the stored procedure call. If keys are not
// specified as integers, an error is returned.
func GetMysqlOpQuery(operation Operation) (string, error) {
	inputs, err := getMysqlInputs(operation.Inputs)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("CALL %s(%s)", operation.Name, inputs), nil
}

func getMysqlInputs(inputs map[string]string) (string, error) {
	if len(inputs) == 0 {
		return "", nil
	}

	sortedParams := make([]string, len(inputs))
	// Store the values in slice in sorted order
	for k, v := range inputs {
		numKey, err := strconv.Atoi(k)
		if err != nil {
			if err == strconv.ErrSyntax {
				return "", fmt.Errorf("key of input '%s' should be an int: %s", k, err)
			}
			return "", nil
		}
		sortedParams[numKey] = v
	}

	var result string
	for _, param := range sortedParams {
		result = fmt.Sprintf("%s, '%s'", result, param)
	}
	result = result[2:]

	return result, nil
}
