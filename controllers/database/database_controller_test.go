package controllers_test

import (
	"context"

	databaseclassv1 "github.com/bedag/kubernetes-dbaas/apis/databaseclass/v1"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe(FormatTestDesc(E2e, "Database controller"), func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		DbName      = "database-sample"
		DbNamespace = "default"
		DbcName = "databaseclass-sample-postgres"
		DbcDriver = database.Postgres

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)
	var (
		ctx = context.Background()
	)

	Context("when creating a new Database resource", func() {
		// Initialize DatabaseClass
		By("creating a DatabaseClass")
		dbc := &databaseclassv1.DatabaseClass{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "databaseclass.dbaas.bedag.ch/v1",
				Kind: "DatabaseClass",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: DbcName,
			},
			Spec: databaseclassv1.DatabaseClassSpec{
				Driver:       DbcDriver,
				Operations: map[string]database.Operation{
					database.CreateMapKey: {
						Name: PostgresCreateOpName,
						Inputs: map[string]string{
							"k8sName": "{{ .Metadata.uid }}",
						},
					},
				},
				SecretFormat: database.SecretFormat{
					"username": "{{ .Result.username }}",
					"password": "{{ .Result.password }}",
					"port": "{{ .Result.port }}",
					"dbName": "{{ .Result.dbName }}",
					"server": "{{ .Result.fqdn }}",
					"dsn": "postgres://{{ .Result.username }}:{{ .Result.password }}@{{ .Result.fqdn }}/{{ .Result.dbName }}",
				},
			},
		}

		It("should create a DatabaseClass", func() {
			Expect(k8sClient.Create(ctx, dbc)).Should(Succeed())
		})
		// It assertion on Create

		// Eventually to retry a number of times until either the function's output matches the Should() assertion,
		// or the number of attempts * interval period exceed the provided timeout value

		// Now that we've create a Database in our test cluster, we must test that the DB actually created it
		// Using Consistently we can ensure that a status field remains set to a certain value for a certain amount of time

		// Also test as much behaviour as possible, e.g. Secret recreation
	})
})