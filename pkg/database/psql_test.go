package database_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe(FormatTestDesc(Integration, "Postgres CreateDb"), func() {
	// Setting up connection to DBMS
	dsn, err := database.Dsn(os.Getenv("POSTGRES_DSN")).GenPostgres()
	Expect(err).ToNot(HaveOccurred())

	conn, err := database.NewPsqlConn(dsn)
	Expect(err).ToNot(HaveOccurred())

	// Prepare assertion data
	opResultAssertion := database.OpOutput{
		Result: map[string]string{
			"username":     "testuser",
			"password":     "testpassword",
			"dbName":       "my-test-db",
			"fqdn":         "localhost",
			"port":         "5432",
			"lastRotation": "",
		},
		Err: nil,
	}

	Context("when Operation is defined correctly", func() {
		// Prepare test data
		createOperation := database.Operation{
			Name: PostgresCreateOpName,
			Inputs: map[string]string{
				"k8sName": "my-test-db",
			},
		}

		// Execute tested operation
		var result database.OpOutput
		result = conn.CreateDb(createOperation)

		It("should not return an error", func() {
			Expect(result.Err).ToNot(HaveOccurred())
		})

		It("should return a non-nil stored procedure Result", func() {
			Expect(result.Result).ToNot(BeNil())
		})

		It("should return a rowset as specified in the stored procedure", func() {
			Expect(result).To(Equal(opResultAssertion))
		})
	})
	Context("when Operation is defined wrongly", func() {
		// Prepare test data
		createOperation := database.Operation{
			Name: "fake_sp_name",
			Inputs: map[string]string{
				"k8sName": "my-test-db",
			},
		}

		// Execute tested operation
		var result database.OpOutput
		result = conn.CreateDb(createOperation)

		It("should return an error", func() {
			Expect(result.Err).To(HaveOccurred())
		})
	})
})
