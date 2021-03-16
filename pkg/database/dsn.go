package database

import "strings"

type Dsn string

func (s Dsn) GetDriver() string {
	return strings.Split(string(s), ":")[0]
}

// NewDsn initialized a new Dsn string.
func NewDsn(driver, username, password, host, port, dbName string) Dsn {
	return Dsn(driver + "://" + username + ":" + password + "@" + host + ":" + port + "/" + dbName)
}

// String returns a string from a Dsn.
func (s Dsn) String() string {
	return string(s)
}
