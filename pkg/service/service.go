package service

import "github.com/bedag/kubernetes-dbaas/pkg/database"

type DbService interface {
	// TODO pass "enum"
	CreateDb(params ...string) ([]string, error)
	DeleteDb() error
}

// Add service specific state here
type DbmsConn struct {
	DbService
}

func Open(driver string) (*DbmsConn, error) {
	var dbmsConn *DbmsConn

	if driver == "mssql" {
		mssqlConn, err := database.NewMssqlConn()
		if err != nil {
			return nil, err
		}
		dbmsConn = &DbmsConn{mssqlConn}
	}
	if driver == "psql" {
		psqlConn, err := database.NewPsqlConn()
		if err != nil {
			return nil, err
		}
		dbmsConn = &DbmsConn{psqlConn}
	}
	return dbmsConn, nil
}
