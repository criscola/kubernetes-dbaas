package database_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const spNameEav = "sp_create_rowset_EAV"

var _ = Describe(FormatTestDesc(Integration, "SQLServer CreateDb"), func() {
	// Setting up connection to DBMS
	dsn := "sqlserver://sa:Password&1@localhost:1433"
	conn, err := database.NewMssqlConn(database.Dsn(dsn))
	Expect(err).ToNot(HaveOccurred())

	Context("when Operation is defined correctly", func() {
		// Prepare test data
		createOperation := database.Operation{
			Name: spNameEav,
			Inputs: map[string]string{
				"k8sName": "integrTestDb",
			},
		}

		// Prepare assertion data
		opResultAssertion := database.OpOutput{
			Result: map[string]string{
				"username": "testuser",
				"password": "testpassword",
				"dbName": "integrTestDb",
				"fqdn": "localhost",
				"port": "1433",
			},
			Err: nil,
		}

		// Execute tested operation
		result := conn.CreateDb(createOperation)

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
})