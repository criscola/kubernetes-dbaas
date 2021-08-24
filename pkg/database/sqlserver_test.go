package database_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

const SqlserverCreateOpNameUnescapedBug = "sp_create_rowset_eav_unescaped_bug"

var _ = Describe(FormatTestDesc(Integration, "Sqlserver CreateDb"), func() {
	var opResultAssertion database.OpOutput
	var result database.OpOutput
	var createOperation database.Operation

	// Setting up connection to DBMS
	dsn, err := database.Dsn(os.Getenv("SQLSERVER_DSN")).GenSqlserver()
	Expect(err).ToNot(HaveOccurred())

	conn, err := database.NewSqlserverConn(dsn)
	Expect(err).ToNot(HaveOccurred())

	BeforeEach(func() {
		// Prepare test data
		createOperation = database.Operation{
			Name: SqlserverCreateOpName,
			Inputs: map[string]string{
				"k8sName": "db-name-with-dashes",
			},
		}
	})
	JustBeforeEach(func() {
		// Prepare assertion data
		opResultAssertion = database.OpOutput{
			Result: map[string]string{
				"username":     "testuser",
				"password":     "testpassword",
				"dbName":       "db-name-with-dashes",
				"fqdn":         "localhost",
				"port":         "1433",
				"lastRotation": "",
			},
			Err: nil,
		}

		result = conn.CreateDb(createOperation)
	})
	Context("when Operation is defined correctly", func() {
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
	Context("when the inputs of an Operation are defined wrongly", func() {
		BeforeEach(func() {
			By("supplying a wrong key name")
			// Prepare test data
			createOperation = database.Operation{
				Name: MysqlCreateOpName,
				Inputs: map[string]string{
					"wrongKey": "db-name-with-dashes",
				},
			}
		})
		It("should return an error", func() {
			Expect(result.Err).To(HaveOccurred())
		})
	})
	Context("when the operation name is defined wrongly", func() {
		BeforeEach(func() {
			// Prepare test data
			createOperation = database.Operation{
				Name: "fake_sp_name",
				Inputs: map[string]string{
					"k8sName": "myTestDb",
				},
			}
		})
		It("should return an error", func() {
			Expect(result.Err).To(HaveOccurred())
		})
	})
	// TODO: Investigate: bug in the driver (sql.Named of pgx)? The connection should return an error, because the
	// create statement is not executed and calling the procedure manually does produce an error.
	PContext("when the unescaped stored procedure is used in the Create call", func() {
		createOperation = database.Operation{
			Name: SqlserverCreateOpNameUnescapedBug,
			Inputs: map[string]string{
				"k8sName": "db-name-with-dashes",
			},
		}
		It("should return a syntax error due to unescaped characters", func() {
			Expect(result.Err).To(HaveOccurred())
		})
	})
})
