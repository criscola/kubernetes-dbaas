package database_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/bedag/kubernetes-dbaas/pkg/database/sqlserver"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

const spNameEav = "sp_create_rowset_EAV"

var _ = Describe(FormatTestDesc(Integration, "NewRateLimitedDbmsConn", Slow), func() {
	// Create a connection
	// Setting up connection to DBMS
	dsn := "sqlserver://sa:Password&1@localhost:1433"
	conn, err := sqlserver.NewMssqlConn(database.Dsn(dsn))
	Expect(err).ToNot(HaveOccurred())

	rateLimitedConn, err := database.NewRateLimitedDbmsConn(conn, 1)
	Expect(err).ToNot(HaveOccurred())

	createOperation := database.Operation{
		Name: spNameEav,
		Inputs: map[string]string{
			"k8sName": "rateLimiterTest",
		},
	}

	Context("when CreateDb is called 10 times in a row", func() {
		Context("when RPS is 1", func() {
			var wg sync.WaitGroup

			var callTimes []time.Time
			beforeAll := time.Now()
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer func() { callTimes = append(callTimes, time.Now()) }()
					defer wg.Done()

					rateLimitedConn.CreateDb(createOperation)
				}()
			}
			wg.Wait()
			elapsedSeconds := time.Now().Sub(beforeAll).Seconds()
			It("should not take less than 9 seconds", func() {
				Expect(elapsedSeconds).NotTo(BeNumerically("<", 9))
			})
			It("should execute each operation with a pause of 1 second in-between", func() {
				for i := 0; i < len(callTimes) - 1; i++ {
					diff := int(callTimes[1].Sub(callTimes[0]).Seconds())  // truncate float
					Expect(diff).To(BeNumerically("==", 1))
				}
			})
		})
	})
})