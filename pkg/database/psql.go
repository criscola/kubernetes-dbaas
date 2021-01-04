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

func (c *PsqlConn) CreateDb(name string) QueryOutput {
	log.Fatal("psql driver not yet implemented")
	/*
		if err != nil {
			// TODO: Handle error
		}
	*/
	return QueryOutput{}
}

func (c *PsqlConn) DeleteDb(name string) QueryOutput {
	log.Fatal("psql driver not yet implemented")
	// TODO: Implement
	return QueryOutput{}
}

func (c *PsqlConn) Ping() error {
	return c.c.Ping()
}
