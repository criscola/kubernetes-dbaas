package database_test

import (
	"os"

	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(FormatTestDesc(Integration, "Postgres Direct CreateDb"), func() {
	// Setting up connection to DBMS
	dsn, err := database.Dsn(os.Getenv("POSTGRES_DSN")).GenPostgres()
	Expect(err).ToNot(HaveOccurred())

	conn, err := database.NewPsqlDirectConn(dsn)
	Expect(err).ToNot(HaveOccurred())

	Context("when Operation is defined correctly", func() {
		// Prepare test data
		createOperation := database.Operation{
			SqlTemplate: "CREATE DATABASE \"{{ .dbName }}\"",
			Secrets:     make(map[string]string),
			DSN:         os.Getenv("POSTGRES_DSN"),
		}

		// Execute tested operation
		result := conn.CreateDb(createOperation)

		It("should not return an error", func() {
			Expect(result.Err).ToNot(HaveOccurred())
		})

		It("should return a non-nil stored procedure Result", func() {
			Expect(result.Result).ToNot(BeNil())
		})

		It("should return the created db credentials", func() {
			Expect(result.Result).Should(
				HaveKey("username"),
				HaveKey("password"),
				HaveKey("dbName"),
				HaveKey("fqdn"),
				HaveKey("port"),
				HaveKey("lastRotation"),
			)
		})
	})
	Context("when Operation is defined wrongly", func() {
		// Prepare test data
		createOperation := database.Operation{
			SqlTemplate: "NOT_A_SQL_TEMPLATE {{ .missingno }}",
			Secrets:     make(map[string]string),
			DSN:         os.Getenv("POSTGRES_DSN"),
		}

		// Execute tested operation
		result := conn.CreateDb(createOperation)

		It("should return an error", func() {
			Expect(result.Err).To(HaveOccurred())
		})
	})
})
