package database_test

import (
	"fmt"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(FormatTestDesc(Integration, "Mariadb CreateDb"), func() {
	var opResultAssertion database.OpOutput
	var result database.OpOutput
	var createOperation database.Operation

	// Setting up connection to DBMS
	dsn, err := database.Dsn("mariadb://root:Password&1@localhost:3306/mysql").GenMysql()
	Expect(err).ToNot(HaveOccurred())

	conn, err := database.NewMysqlConn(dsn)
	Expect(err).ToNot(HaveOccurred())

	BeforeEach(func() {
		// Prepare test data
		createOperation = database.Operation{
			Name: MysqlCreateOpName,
			Inputs: map[string]string{
				"0": "my-database-test",
			},
		}
	})
	JustBeforeEach(func() {
		// Prepare assertion data
		opResultAssertion = database.OpOutput{
			Result: map[string]string{
				"username": "testuser",
				"password": "testpassword",
				"dbName":   "my-database-test",
				"fqdn":     "localhost",
				"port":     "3306",
			},
			Err: nil,
		}
		// Execute tested behavior
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
			By("supplying strings instead of integers as keys")
			// Prepare test data
			createOperation = database.Operation{
				Name: MysqlCreateOpName,
				Inputs: map[string]string{
					"k8sName": "myTestDb",
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
					"0": "myTestDb",
				},
			}
		})
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
			Name: MysqlCreateOpName,
			Inputs: map[string]string{
				"4": "param4",
				"0": "param0",
				"1": "param1",
				"3": "param3",
				"2": "param2",
			},
		}
		outputAssert := fmt.Sprintf("CALL %s('param0', 'param1', 'param2', 'param3', 'param4')", MysqlCreateOpName)

		// Execute tested operation
		val, _ := database.GetMysqlOpQuery(createOperation)
		It("should match the expected output", func() {
			Expect(val).To(Equal(outputAssert))
		})
	})
})
