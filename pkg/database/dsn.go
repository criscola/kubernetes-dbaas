package database

import "strings"

type Dsn string

func (s Dsn) GetDriver() string {
	return strings.Split(string(s), ":")[0]
}

func (s Dsn) String() string {
	return string(s)
}

// WithTable set the table associated with the connection.
func (s Dsn) WithTable(table string) Dsn {
	// TODO: More checks
	return Dsn(s.String() + "/" + table)
}
