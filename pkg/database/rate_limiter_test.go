package database_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"sync"
	"time"
)

var _ = Describe(FormatTestDesc(Integration, "NewRateLimitedDbmsConn", Slow), func() {
	// Setting up connection to DBMS
	dsn, err := database.Dsn(os.Getenv("SQLSERVER_DSN")).GenSqlserver()
	Expect(err).ToNot(HaveOccurred())

	conn, err := database.NewSqlserverConn(dsn)
	Expect(err).ToNot(HaveOccurred())

	rateLimitedConn, err := database.NewRateLimitedDbmsConn(conn, 1)
	Expect(err).ToNot(HaveOccurred())

	createOperation := database.Operation{
		Name: SqlserverCreateOpName,
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
				go func(callTimes []time.Time) {
					defer func() { callTimes = append(callTimes, time.Now()) }()
					defer wg.Done()

					rateLimitedConn.CreateDb(createOperation)
				}(callTimes)
			}
			wg.Wait()
			elapsedSeconds := time.Now().Sub(beforeAll).Seconds()
			It("should not take less than 9 seconds", func() {
				Expect(elapsedSeconds).NotTo(BeNumerically("<", 9))
			})
			It("should execute each operation with a pause of at least 1 second in-between", func() {
				for i := 0; i < len(callTimes)-1; i++ {
					diff := callTimes[1].Sub(callTimes[0]).Seconds()
					Expect(diff).To(BeNumerically(">=", 0.9)) // tolerance of 0.1s
				}
			})
		})
	})
})
