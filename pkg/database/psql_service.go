package database

import "database/sql"

type PsqlConn struct {
	conn *sql.DB
	// TODO add logging
	// Add MSSQL specific state here
}

func (*PsqlConn) CreateDb(params... string) ([]string, error) {
	// TODO: Implement
	return nil, nil
}

func (*PsqlConn) DeleteDb() error {
	// TODO: Implement
	return nil
}