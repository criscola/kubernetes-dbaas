package database

import (
	"database/sql"
	"fmt"
)

type MssqlConn struct {
	c *sql.DB
}

func NewMssqlConn() (*MssqlConn, error) {
	dbConn, err := sql.Open("mssql", "...")
	if err != nil {
		// TODO: Handle error
	}
	conn := MssqlConn{dbConn}
	return &conn, nil
}

func (c *MssqlConn) CreateDb(name, stage string) ([]string, error) {
	rows, err := c.c.Query("CALL ...")

	if err != nil {
		// TODO: Handle error
	}
	fmt.Print(rows)
	return nil, nil
}

func (c *MssqlConn) DeleteDb() error {
	// TODO: Implement
	return nil
}
