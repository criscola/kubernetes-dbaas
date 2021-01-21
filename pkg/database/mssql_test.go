package database_test

import (
	"database/sql"
	"fmt"
	. "github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
)

var sqlserverDsn Dsn

// TestMain is executed before each test case.
func TestMain(m *testing.M) {
	var resource *dockertest.Resource
	var err error

	// Set up
	sqlserverDsn, _, resource, err = SetupSingleSqlServerContainer()
	if err != nil {
		log.Fatal(FormatSqlServerContainerSetupError(err))
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if err = TeardownContainer(resource); err != nil {
		log.Fatal(FormatContainerTeardownError(err))
	}

	os.Exit(code)
}

// TODO: Update test to reflect new templating feature
func TestMssqlConn_CreateDb(t *testing.T) {
	/*
		ops := GetMockOps()

		conn, err := NewMssqlConn(sqlserverDsn, ops)
		if err != nil {
			t.Fatalf("could not initialize MssqlConn: %s", err)
		}
		out := conn.CreateDb(MockK8sName)
		if out.Err != nil {
			t.Fatalf("could not create database: %s", out.Err)
		}
		if out.Out[0] != MockOutputUser {
			t.Error("could not get username from stored procedure")
			t.Fatalf("expected: %s got: %s", MockOutputUser, out.Out[0])
		}
		if out.Out[1] != MockOutputPass {
			t.Error("could not get password from stored procedure")
			t.Fatalf("expected: %s got: %s", MockOutputPass, out.Out[1])
		}

		// Let's check if the db was really created...
		isDbPresent, err := dbmsContainsDb(conn)
		if err != nil {
			t.Fatalf("could not create database correctly: %s", err)
		}
		if !isDbPresent {
			t.Fatal("could not create database; database is not present in dbms")
		}*/
}

// TODO: Update test to reflect new templating feature
func TestMssqlConn_DeleteDb(t *testing.T) {
	/*
		t.Run("create db before deletion", TestMssqlConn_CreateDb)

		ops := GetMockOps()

		conn, err := NewMssqlConn(sqlserverDsn, ops)
		if err != nil {
			t.Fatalf("could not initialize MssqlConn: %s", err)
		}
		out := conn.DeleteDb(MockK8sName)
		if out.Err != nil {
			t.Fatalf("could not create database: %s", out.Err)
		}
		// Let's check if the db was really deleted...
		isDbPresent, err := dbmsContainsDb(conn)
		if err != nil {
			t.Fatalf("could not create database correctly: %s", err)
		}
		if isDbPresent {
			t.Fatal("could not delete database; database is still present in dbms")
		}*/
}

// Checks that the dbms contains the mock database, returns true if it contains it, false otherwise
func dbmsContainsDb(conn *MssqlConn) (bool, error) {
	c := getClient(conn)
	// TODO: Agree on how the database names are treated. For the moment, the dbname is DbNamePrefix + the first eight chars of the dbresource UID)
	left, _ := strconv.Atoi(LeftTrimLength)
	dbName := DbNamePrefix + MockK8sName[:left]
	row, err := c.Query("SELECT COUNT((DB_ID('" + dbName + "')))")
	if err != nil {
		return false, fmt.Errorf("could not create database: %s", err)
	}
	row.Next()
	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("could not scan database count: %s", err)
	}
	if count != 1 {
		return false, nil
	}
	return true, nil
}

func getClient(conn *MssqlConn) *sql.DB {
	v := reflect.ValueOf(conn).Elem().FieldByName("c")
	return GetUnexportedField(v).(*sql.DB)
}
