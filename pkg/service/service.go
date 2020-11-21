package service

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
)

type DbService interface {
	CreateDb(name, stage string) ([]string, error)
	DeleteDb() error
}

// Add service specific state here
type DbmsConn struct {
	DbService
}

// TODO: Add "stage"
func Open(driver string) (*DbmsConn, error) {
	var dbmsConn *DbmsConn

	if driver == "mssql" {
		mssqlConn, err := database.NewMssqlConn()
		if err != nil {
			return nil, err
		}
		dbmsConn = &DbmsConn{mssqlConn}
	} else if driver == "psql" {
		psqlConn, err := database.NewPsqlConn()
		if err != nil {
			return nil, err
		}
		dbmsConn = &DbmsConn{psqlConn}
	}
	return dbmsConn, nil
}

func GetSupportedDbms() []string {
	return []string{"mssql", "psql"}
}
