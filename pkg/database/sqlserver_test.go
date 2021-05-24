package database_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const SqlserverCreateOpNameUnescapedBug = "sp_create_rowset_eav_unescaped_bug"

var _ = Describe(FormatTestDesc(Integration, "Sqlserver CreateDb"), func() {
	// Setting up connection to DBMS
	dsn, err := database.Dsn("sqlserver://sa:Password&1@localhost:1433").GenSqlserver()
	Expect(err).ToNot(HaveOccurred())

	conn, err := database.NewSqlserverConn(dsn)
	Expect(err).ToNot(HaveOccurred())

	Context("when Operation is defined correctly", func() {
		// Prepare test data
		createOperation := database.Operation{
			Name: SqlserverCreateOpName,
			Inputs: map[string]string{
				"k8sName": "db-name-with-dashes",
			},
		}

		// Prepare assertion data
		opResultAssertion := database.OpOutput{
			Result: map[string]string{
				"username": "testuser",
				"password": "testpassword",
				"dbName":   "db-name-with-dashes",
				"fqdn":     "localhost",
				"port":     "1433",
			},
			Err: nil,
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
	// TODO: Investigate: bug in the driver (sql.Named of pgx)? The connection should return an error, because the
	// create statement is not executed and calling the procedure manually does produce an error.
	PContext("when the unescaped stored procedure is used in the Create call", func() {
		createOperationUnescaped := database.Operation{
			Name: SqlserverCreateOpNameUnescapedBug,
			Inputs: map[string]string{
				"k8sName": "db-name-with-dashes",
			},
		}
		// Execute tested operation
		var result database.OpOutput
		result = conn.CreateDb(createOperationUnescaped)

		It("should return a syntax error due to unescaped characters", func() {
			Expect(result.Err).To(HaveOccurred())
		})
	})
})
