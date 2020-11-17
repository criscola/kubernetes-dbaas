package database

import "database/sql"

type MssqlService struct {
	conn *sql.DB
	// TODO add logging
	// Add MSSQL specific state here
}

func (*MssqlService) CreateDb(params... string) ([]string, error) {
	// TODO: Implement
	return nil, nil
}

func (*MssqlService) DeleteDb() error {
	// TODO: Implement
	return nil
}