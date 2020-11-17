package database

import (
	"database/sql"
	"fmt"
)

// Add service specific state here
type DbConn struct {
	DbService
}

// Inject dependency
// TODO: Type safety
func NewDbConn(driver string) (*DbConn, error) {
	if driver == "mssql" {
		return &DbConn{NewMssqlConn(driver)}, nil
	}
	if driver == "psql" {
		return &DbConn{NewPsqlConn(driver)}, nil
	}
	// TODO: More logging
	return nil, fmt.Errorf("error connecting to %s", driver)
}

func NewMssqlConn(driver string) *MssqlService {
	// TODO: Implement
	conn, _ := sql.Open(driver, "driver://username:password@host/")
	return &MssqlService{conn: conn}
}

func NewPsqlConn(driver string) *PsqlConn {
	// TODO: Implement
	conn, _ := sql.Open(driver, "driver://username:password@host/")
	return &PsqlConn{conn: conn}
}
