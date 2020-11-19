package database

import (
	"database/sql"
	"fmt"
)

type PsqlConn struct {
	c *sql.DB
	// TODO add logging
	// Add MSSQL specific state here
}

func NewPsqlConn() (*MssqlConn, error) {
	dbConn, err := sql.Open("psql", "...")
	if err != nil {
		// TODO: Handle error
	}
	conn := MssqlConn{dbConn}
	return &conn, nil
}

func (c *PsqlConn) CreateDb(params ...string) ([]string, error) {
	rows, err := c.c.Query("CALL ...")

	if err != nil {
		// TODO: Handle error
	}
	fmt.Print(rows)
	return nil, nil
}

func (c *PsqlConn) DeleteDb() error {
	// TODO: Implement
	return nil
}
