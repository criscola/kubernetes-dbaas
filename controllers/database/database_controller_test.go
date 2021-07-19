package controllers_test

import (
	"context"
	"fmt"
	databasev1 "github.com/bedag/kubernetes-dbaas/apis/database/v1"
	. "github.com/bedag/kubernetes-dbaas/controllers/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	"github.com/bedag/kubernetes-dbaas/pkg/typeutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	k8sError "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"path"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	DbMariadbFilename   = "db-mariadb.yaml"
	DbPostgresFilename  = "db-postgres.yaml"
	DbSqlserverFilename = "db-sqlserver.yaml"
)

var _ = Describe(FormatTestDesc(E2e, "Database controller"), func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	// Prepare Database resources
	postgresDatabaseRes, err := getDbFromTestdata(DbPostgresFilename)
	Expect(err).NotTo(HaveOccurred())
	sqlserverDatabaseRes, err := getDbFromTestdata(DbSqlserverFilename)
	Expect(err).NotTo(HaveOccurred())
	mariadbDatabaseRes, err := getDbFromTestdata(DbMariadbFilename)
	Expect(err).NotTo(HaveOccurred())

	Context("when testing the happy path", func() {
		Context("when reconciling a PostgreSQL Database resource", func() {
			It("should handle its lifecycle correctly", func() {
				testDatabaseLifecycleHappyPath(postgresDatabaseRes, duration, timeout, interval)
			})
		})
		Context("when reconciling a MariaDB Database resource", func() {
			It("should handle its lifecycle correctly", func() {
				testDatabaseLifecycleHappyPath(mariadbDatabaseRes, duration, timeout, interval)
			})
		})
		Context("when reconciling a SQLServer Database resource", func() {
			It("should handle its lifecycle correctly", func() {
				testDatabaseLifecycleHappyPath(sqlserverDatabaseRes, duration, timeout, interval)
			})
		})
	})
	// It assertion on Create


	// Eventually to retry a number of times until either the function's output matches the Should() assertion,
	// or the number of attempts * interval period exceed the provided timeout value

	// Now that we've create a Database in our test cluster, we must test that the DB actually created it
	// Using Consistently we can ensure that a status field remains set to a certain value for a certain amount of time

	// Also test as much behaviour as possible, e.g. Secret recreation
})

func testDatabaseLifecycleHappyPath(v databasev1.Database, duration, timeout, interval interface{}) {
	By("creating the API resource successfully with condition Ready set to true", func() {
		//Expect(err).NotTo(HaveOccurred())
		err := k8sClient.Create(context.Background(), &v)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() error {
			return checkDbReady(&v)
		}, timeout, interval).Should(BeNil())
		// We don't just check the Ready condition would be eventually True, we also check that it
		// stays that way for a certain period of time as an additional assurance
		Consistently(func() error {
			return checkDbReady(&v)
		}, duration, interval).Should(BeNil())
	})
	By("creating the relative Secret resource successfully", func() {
		Eventually(func() error {
			secret := v1.Secret{}
			err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: v.Namespace,
				Name: FormatSecretName(&v)}, &secret)
			return err
		}, timeout, interval).Should(BeNil())
	})
	By("rotating the credentials", func() {
		// TODO: Create Rotate sample stored procedure and then implement credential rotation
		// TODO: Get secret data, apply rotation, get secret data again and compare it with the older data
		// Expect password to have changed. Expect annotation to be removed.
	})
	By("deleting the API resource successfully", func() {
		err := k8sClient.Delete(context.Background(), &v)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() bool {
			return k8sError.IsNotFound(checkDbReady(&v))
		}, timeout, interval).Should(BeTrue())
	})
}

func getDbFromTestdata(filename string) (databasev1.Database, error) {
	return readDbYaml(filename)
}

func readDbYaml(filename string) (databasev1.Database, error) {
	dbFilepath := path.Join(testdataFilepath, filename)
	db := databasev1.Database{}
	dat, err := ioutil.ReadFile(dbFilepath)
	if err != nil {
		return databasev1.Database{}, err
	}
	err = yaml.Unmarshal(dat, &db)
	if err != nil {
		return databasev1.Database{}, err
	}
	return db, nil
}

func checkDbReady(db *databasev1.Database) error {
	// Get a fresh Database resource from the API server
	freshDb := databasev1.Database{}
	err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(db), &freshDb)
	if err != nil {
		return err
	}
	ready := meta.FindStatusCondition(freshDb.Status.Conditions, typeutil.TypeReady)
	if ready == nil {
		return fmt.Errorf("ready condition is nil")
	}
	if ready.Status != metav1.ConditionTrue {
		return fmt.Errorf("database is not ready: %s: %s", ready.Reason, ready.Message)
	}
	return nil
}