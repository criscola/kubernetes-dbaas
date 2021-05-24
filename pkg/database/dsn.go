package database

import "github.com/xo/dburl"

type Dsn string

func (dsn Dsn) GenMysql() (string, error) {
	u, err := dburl.Parse(dsn.String())
	if err != nil {
		return "", err
	}
	parsedDsn, err := dburl.GenMySQL(u)
	if err != nil {
		return "", err
	}
	return parsedDsn, nil
}

func (dsn Dsn) GenSqlserver() (string, error) {
	u, err := dburl.Parse(dsn.String())
	if err != nil {
		return "", err
	}
	parsedDsn, err := dburl.GenSQLServer(u)
	if err != nil {
		return "", err
	}
	return parsedDsn, nil
}

func (dsn Dsn) GenPostgres() (string, error) {
	u, err := dburl.Parse(dsn.String())
	if err != nil {
		return "", err
	}
	parsedDsn, err := dburl.GenPostgres(u)
	if err != nil {
		return "", err
	}
	return parsedDsn, nil
}

// String returns a string from a Dsn.
func (dsn Dsn) String() string {
	return string(dsn)
}
