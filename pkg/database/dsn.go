package database

import "strings"

type Dsn string

// Must be like driver://username:password@host/instance?param1=value&param2=value
func (s Dsn) GetDriver() string {
	return strings.Split(string(s), ":")[0]
}

func (s Dsn) String() string {
	return string(s)
}

func (s Dsn) WithTable(table string) Dsn {
	// TODO: More checks
	return Dsn(s.String() + "/" + table)
}
