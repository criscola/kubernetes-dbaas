// This package opens and retains a pool of distinct DBMS connections.
package pool

import (
	"fmt"
	databaseclassv1 "github.com/bedag/kubernetes-dbaas/apis/databaseclass/v1"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
)

var pool dbmsPool

// key is the endpoint name.
type dbmsPool map[string]dbmsPoolEntry

type dbmsPoolEntry struct {
	dbmsConn *database.DbmsConn
	dsn      string
	driver   string
}

// Register registers a new database.Dbms in the pool.
func Register(dbms database.Dbms, dbClass databaseclassv1.DatabaseClass) error {
	// Get driver from DatabaseClass.
	driver := dbClass.Spec.Driver
	for _, endpoint := range dbms.Endpoints {
		if _, exists := pool[endpoint.Name]; exists {
			return fmt.Errorf("%s is already present in the pool. Endpoint names must be unique within the list "+
				"of endpoints", endpoint.Name)
		}

		conn, err := database.NewDbmsConn(driver, endpoint.Dsn)
		if err != nil {
			return fmt.Errorf("problem opening connection to endpoint: %s", err)
		}
		// Add entry to pool
		pool[endpoint.Name] = dbmsPoolEntry{conn, endpoint.Dsn.String(), driver}
	}

	return nil
}

// GetConnByEndpointName tries to get a connection by endpoint name. It returns an error if a connection related to
// endpointName is not found in the pool.
func GetConnByEndpointName(endpointName string) (*database.DbmsConn, error) {
	if conn, ok := pool[endpointName]; ok {
		return conn.dbmsConn, nil
	}
	return nil, fmt.Errorf("entry '%s' not found in dbms pool", endpointName)
}

func init() {
	pool = make(map[string]dbmsPoolEntry)
}
