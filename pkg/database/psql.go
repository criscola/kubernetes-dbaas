package database

import (
	"database/sql"
	"log"
)

type PsqlConn struct {
	c *sql.DB
}

func NewPsqlConn(dsn string) (*PsqlConn, error) {
	log.Fatal("psql driver not yet implemented")
	/*
		dbConn, err := sql.Open("psql", dsn)
		if err != nil {
			// TODO: Handle error
		}
		conn := PsqlConn{dbConn}
	*/
	return nil, nil
}

func (c *PsqlConn) CreateDb(name string) OpOutput {
	log.Fatal("psql driver not yet implemented")
	/*
		if err != nil {
			// TODO: Handle error
		}
	*/
	return OpOutput{}
}

func (c *PsqlConn) DeleteDb(name string) OpOutput {
	log.Fatal("psql driver not yet implemented")
	// TODO: Implement
	return OpOutput{}
}

func (c *PsqlConn) Ping() error {
	return c.c.Ping()
}
