// This package attempts at opening a   nd retaining a pool of distinct DBMS connections
// TODO: Versioning. Especially for this package (service layer)
package pool

import (
	"fmt"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
)

var pool dbmsPool

// Pool of DBMS connection pools where the key is the driver and the value is a slice of dbmsPoolEntry
type dbmsPool map[string][]dbmsPoolEntry

//
type dbmsPoolEntry struct {
	dbmsConn   *database.DbmsConn
	dbmsConfig database.Endpoint
}

// Register registers a new database.Dbms in the pool
func Register(dbms database.Dbms) error {
	// Register dbms endpoints
	driver := dbms.Driver
	for _, endpoint := range dbms.Endpoints {
		conn, err := database.New(endpoint.Dsn, dbms.Operations)
		if err != nil {
			return err
		}
		// Add entry to pool
		pool[driver] = append(pool[driver], dbmsPoolEntry{conn, endpoint})
	}

	return nil
}

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

func GetDsnByDriverAndEndpointName(driver, endpointName string) (database.Dsn, error) {
	for _, v := range pool[driver] {
		if v.dbmsConfig.Name == endpointName {
			return v.dbmsConfig.Dsn, nil
		}
	}
	return "", fmt.Errorf("entry '%s' with driver '%s' not found in dbms pool", endpointName, driver)
}

func SizeOf(driver string) int {
	return len(pool[driver])
}

func String() string {
	return fmt.Sprint(pool)
}

func init() {
	pool = make(map[string][]dbmsPoolEntry)
}
