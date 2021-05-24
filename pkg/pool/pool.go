// This package opens and retains a pool of distinct DBMS connections.
package pool

import (
	"fmt"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
)

// Pool specifies the generic interface for a Pool of DBMS connections.
type Pool interface {
	Get(name string) Entry
	Register(name string, driver string, dsn database.Dsn) error
}

// Entry specifies the generic interface for an entry of DbmsPool.
type Entry interface {
	database.Driver
}

// DbmsPool is a map of pool entries identified by a unique name.
type DbmsPool struct {
	entries map[string]Entry
	rps     int
}

// Get retrieves an Entry from pool.
func (pool DbmsPool) Get(name string) Entry {
	return pool.entries[name]
}

// DbmsEntry represents a standard Dbms connection.
type DbmsEntry struct {
	Entry
	driver string
	dsn    database.Dsn
}

// NewDbmsPool initializes a DbmsPool struct with the given rps. See also database.RateLimitedDbmsConn.
func NewDbmsPool(rps int) DbmsPool {
	return DbmsPool{
		entries: make(map[string]Entry),
		rps:     rps,
	}
}

// RegisterDbms is a utility function around Register. It iterates over database.Dbms.Endpoints and registers a connection for
// each endpoint.
func (pool DbmsPool) RegisterDbms(dbms database.Dbms, driver string) error {
	for _, endpoint := range dbms.Endpoints {
		if err := pool.Register(endpoint.Name, driver, endpoint.Dsn); err != nil {
			return err
		}
	}
	return nil
}

// Register registers a new database.Dbms in the pool.
func (pool DbmsPool) Register(name string, driver string, dsn database.Dsn) error {
	conn, err := database.New(driver, dsn)
	if err != nil {
		return fmt.Errorf("problem opening connection to endpoint with driver: '%s': %s", driver, err)
	}
	rateLimitedConn, err := database.NewRateLimitedDbmsConn(conn, pool.rps)
	if err != nil {
		return err
	}
	if _, exists := pool.entries[name]; exists {
		return fmt.Errorf("%s is already present in the pool. Endpoint names must be unique within the list "+
			"of endpoints", name)
	}
	pool.entries[name] = DbmsEntry{rateLimitedConn, driver, dsn}
	return err
}
