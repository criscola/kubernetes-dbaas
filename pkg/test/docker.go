// Package test provides the facilities to perform integration testing using the package dockertest.
package test

import (
	"database/sql"
	"fmt"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test/dpool"
	"github.com/ory/dockertest/v3"
	"log"
)

const (
	imageRepo       = "mcr.microsoft.com/mssql/server"
	imageTag        = "2019-latest"
	acceptEulaParam = "ACCEPT_EULA=Y"
	saPassword      = "Password&1"
	saPasswordParam = "SA_PASSWORD=" + saPassword

	SqlServerDriver = "sqlserver"
)

// SetupSingleSqlServerContainer provisions a single sql server container
func SetupSingleSqlServerContainer() (database.Dsn, *sql.DB, *dockertest.Resource, error) {
	var dockerPool *dockertest.Pool
	var sqlserverDsn database.Dsn
	var resource *dockertest.Resource
	var db *sql.DB
	var err error

	dockerPool, err = GetDockerPool()
	if err != nil {
		return "", nil, nil, FormatDockerPoolError(err)
	}
	sqlserverDsn, db, resource, err = RunSqlServerContainer(dockerPool)
	if err != nil {
		return "", nil, nil, FormatSqlServerContainerSetupError(err)
	}

	// Preload create and delete stored procedure in db
	_, err = db.Query(CreateSp)
	if err != nil {
		return "", nil, nil, fmt.Errorf("could not load create stored procedure: %s", err)
	}
	_, err = db.Query(DeleteSp)
	if err != nil {
		return "", nil, nil, fmt.Errorf("could not load delete stored procedure: %s", err)
	}
	return sqlserverDsn, db, resource, err
}

// RunSqlServerContainer creates a new docker container running on dockerPool.
func RunSqlServerContainer(dockerPool *dockertest.Pool) (database.Dsn, *sql.DB, *dockertest.Resource, error) {
	var dsn database.Dsn
	var db *sql.DB
	var resource *dockertest.Resource
	var err error

	resource, err = dockerPool.Run(imageRepo, imageTag, []string{acceptEulaParam, saPasswordParam})
	if err != nil {
		return "", nil, nil, fmt.Errorf("could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = dockerPool.Retry(func() error {
		dsn = database.Dsn(SqlServerDriver + "://sa:" + saPassword + "@localhost:" + resource.GetPort("1433/tcp"))
		db, err = sql.Open(SqlServerDriver, dsn.String())
		if err != nil {
			log.Println("failed to connect")
			return err
		}

		return db.Ping()
	}); err != nil {
		return "", nil, nil, fmt.Errorf("could not connect to docker: %s", err)
	}
	return dsn, db, resource, nil
}

// TeardownContainer purges a container and linked volumes from docker.
func TeardownContainer(resource *dockertest.Resource) error {
	dockerPool, err := GetDockerPool()
	if err != nil {
		return FormatDockerPoolError(err)
	}

	if err := dockerPool.Purge(resource); err != nil {
		return fmt.Errorf("could not purge resource: %s", err)
	}
	return nil
}

func FormatDockerPoolError(err error) error {
	return fmt.Errorf("could not connect to docker: %s", err)
}

func FormatDockerContainerError(err error) error {
	return fmt.Errorf("could not start docker container: %s", err)
}

func FormatContainerTeardownError(err error) error {
	return fmt.Errorf("could not tear down sqlserver container: %s", err)
}

func FormatSqlServerContainerSetupError(err error) error {
	return fmt.Errorf("could not set up sqlserver container: %s", err)
}
