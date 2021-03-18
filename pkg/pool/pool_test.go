package pool_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/bedag/kubernetes-dbaas/pkg/test/dpool"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"testing"
)

const driver = SqlServerDriver

var dsn1 database.Dsn
var dsn2 database.Dsn
var dsnSlice []database.Dsn

// TestMain is executed before each test case.
func TestMain(m *testing.M) {
	var resource1 *dockertest.Resource
	var resource2 *dockertest.Resource

	// Setup
	dockerPool, err := GetDockerPool()
	if err != nil {
		log.Fatal(FormatDockerPoolError(err))
	}
	dsn1, _, resource1, err = RunSqlServerContainer(dockerPool)
	if err != nil {
		log.Fatal(FormatDockerContainerError(err))
	}
	dsn2, _, resource2, err = RunSqlServerContainer(dockerPool)
	if err != nil {
		log.Fatal(FormatDockerContainerError(err))
	}
	dsnSlice = []database.Dsn{dsn1, dsn2}

	// Run tests
	code := m.Run()

	// Cleanup
	if err = TeardownContainer(resource1); err != nil {
		log.Fatal(FormatContainerTeardownError(err))
	}

	if err = TeardownContainer(resource2); err != nil {
		log.Fatal(FormatContainerTeardownError(err))
	}

	os.Exit(code)
}

// TestRegister tests whether the pool is able to register a new instance.
// TODO: Update
func TestRegister(t *testing.T) {
	/*
		endpoints := GetMockEndpoints(dsnSlice)
		dbms := database.Dbms{
			Driver:     driver,
			Operations: GetMockOps(),
			Endpoints:  endpoints,
		}
		err := pool.Register(dbms)
		if err != nil {
			t.Fatalf("could not register DBMS in pool: %s", err)
		}
		t.Log(pool.String())
		// There should be an entry for each endpoint of database.Dbms
		if pool.SizeOf(driver) != len(endpoints) {
			t.Error("could not get correct number of pool entries")
			t.Fatalf("expected: %d got: %d", len(endpoints), pool.SizeOf(driver))
		}
	*/
}

// TestGetByDriverAndEndpointName tests whether an entry by driver and endpoint name can be retrieved from the pool.
/*
func TestGetByDriverAndEndpointName(t *testing.T) {
	endpoints := GetMockEndpoints(dsnSlice)
	for _, v := range endpoints {
		dbmsConn, err := pool.GetConnByDriverAndEndpointName(driver, v.Name)
		if err != nil {
			t.Fatalf("could not get DBMS (%s, %s) from pool: %s", driver, v.Name, err)
		}
		// TODO: How to check that a connections matches its Endpoint? Probably need to expose something in DbmsConn
		if dbmsConn == nil {
			t.Fatal("dbmsConn cannot be nil")
		}
		if err := dbmsConn.Ping(); err != nil {
			t.Fatalf("could not ping database (%s, %s): %s ", v.Name, v.Dsn, err)
		}
	}
}
*/
