package database_test

import (
	"fmt"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const mysqlCreateOpName = "sp_create_db_rowset_eav"

var _ = Describe(FormatTestDesc(Integration, "Mariadb CreateDb"), func() {
	// Setting up connection to DBMS
	dsn, err := database.Dsn("mariadb://root:Password&1@localhost:3306/mysql").GenMysql()
	Expect(err).ToNot(HaveOccurred())

	conn, err := database.NewMysqlConn(dsn)
	Expect(err).ToNot(HaveOccurred())

	// Prepare assertion data
	opResultAssertion := database.OpOutput{
		Result: map[string]string{
			"username": "testuser",
			"password": "testpassword",
			"dbName": "myTestDb",
			"fqdn": "localhost",
			"port": "3306",
		},
		Err: nil,
	}

	Context("when Operation is defined correctly", func() {
		// Prepare test data
		createOperation := database.Operation{
			Name: mysqlCreateOpName,
			Inputs: map[string]string{
				"0": "myTestDb",
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

	Context("when an Operation is defined wrongly", func() {
		By("supplying strings instead of ints as key of the inputs of the create procedure")

		// Prepare test data
		createOperation := database.Operation{
			Name: mysqlCreateOpName,
			Inputs: map[string]string{
				"k8sName": "myTestDb",
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

var _ = Describe(FormatTestDesc(Unit, "GetMysqlOpQuery"), func() {
	Context("when Operation is defined correctly", func() {
		By("having 5 inputs")
		// Prepare test data
		createOperation := database.Operation{
			Name: mysqlCreateOpName,
			Inputs: map[string]string{
				"4": "param4",
				"0": "param0",
				"1": "param1",
				"3": "param3",
				"2": "param2",
			},
		}

		outputAssert := fmt.Sprintf("CALL %s('param0', 'param1', 'param2', 'param3', 'param4')", mysqlCreateOpName)

		// Execute tested operation
		val, _ := database.GetMysqlOpQuery(createOperation)

		It("should match the expected output", func() {
			Expect(val).To(Equal(outputAssert))
		})
	})
})