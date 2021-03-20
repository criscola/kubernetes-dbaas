// This package attempts at opening a   nd retaining a pool of distinct DBMS connections
package pool

import (
	"fmt"
	dbaasv1 "github.com/bedag/kubernetes-dbaas/apis/databaseclass/v1"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
)

var pool dbmsPool

// key is the endpoint name
type dbmsPool map[string]dbmsPoolEntry

type dbmsPoolEntry struct {
	dbmsConn *database.DbmsConn
	dsn      string
	driver   string
}

func Register(dbms database.Dbms, dbClass dbaasv1.DatabaseClass) error {
	// Get driver from DatabaseClass
	driver := dbClass.Spec.Driver
	for _, endpoint := range dbms.Endpoints {
		conn, err := database.New(endpoint.Dsn, dbClass.Spec.Operations)
		if err != nil {
			return err
		}
		// Add entry to pool
		pool[endpoint.Name] = dbmsPoolEntry{conn, endpoint.Dsn.String(), driver}
	}

	return nil
}

func GetConnByEndpointName(endpointName string) (*database.DbmsConn, error) {
	if conn, ok := pool[endpointName]; ok {
		// Extra check in case the connection has gone down, probably unnecessary because database/sql reopens
		// db connections when necessary
		if err := conn.dbmsConn.Ping(); err != nil {
			return nil, err
		}
		return conn.dbmsConn, nil
	}
	return nil, fmt.Errorf("entry '%s' not found in dbms pool", endpointName)
}

/*

// Pool of DBMS connection pools where the key is the driver and the value is a slice of dbmsPoolEntry
type dbmsPool map[string][]dbmsPoolEntry

type dbmsPoolEntry struct {
	dbmsConn   *database.DbmsConn
	dbmsConfig database.Endpoint
}

// Register registers a new database.Dbms in the pool
func Register(dbms database.Dbms, dbClass dbaasv1.DatabaseClass) error {
	// Get driver from DatabaseClass
	driver := dbClass.Spec.Driver
	for _, endpoint := range dbms.Endpoints {
		conn, err := database.New(endpoint.Dsn, dbClass.Spec.Operations)
		if err != nil {
			return err
		}
		// Add entry to pool
		pool[driver] = append(pool[driver], dbmsPoolEntry{conn, endpoint})
	}

	return nil
}

// GetConnByDriverAndEndpointName tries to retrieve a database.DbmsConn from the pool of connections.
func GetConnByDriverAndEndpointName(driver, endpointName string) (*database.DbmsConn, error) {
	for _, v := range pool[driver] {
		if v.dbmsConfig.Name == endpointName {
			// Extra check in case the connection has gone down, probably unnecessary because database/sql reopens
			// db connections when necessary
			if err := v.dbmsConn.Ping(); err != nil {
				return nil, err
			}
			return v.dbmsConn, nil
		}
	}
	return nil, fmt.Errorf("entry '%s' with driver '%s' not found in dbms pool", endpointName, driver)
}

func GetConnByEndpointName(endpointName string) (*database.DbmsConn, error) {
	for _, v := range pool[driver] {
		if v.dbmsConfig.Name == endpointName {
			// Extra check in case the connection has gone down, probably unnecessary because database/sql reopens
			// db connections when necessary
			if err := v.dbmsConn.Ping(); err != nil {
				return nil, err
			}
			return v.dbmsConn, nil
		}
	}
	return nil, fmt.Errorf("entry '%s' with driver '%s' not found in dbms pool", endpointName, driver)
}

// SizeOf returns the number of connections in the current pool
func SizeOf(driver string) int {
	return len(pool[driver])
}

// String returns the pool formatted as a string.
func String() string {
	return fmt.Sprint(pool)
}

*/

func init() {
	pool = make(map[string]dbmsPoolEntry)
}
