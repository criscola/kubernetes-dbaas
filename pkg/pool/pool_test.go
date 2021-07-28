package pool_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/bedag/kubernetes-dbaas/pkg/pool"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(FormatTestDesc(Integration, "RegisterDbms"), func() {
	var dbmsPool pool.DbmsPool
	var dbmsList []database.Dbms
	var err error
	BeforeEach(func() {
		dbmsList = []database.Dbms{{
			DatabaseClassName: "testDbc",
			Endpoints: []database.Endpoint{
				{
					Name: "postgres",
					Dsn:  "postgres://postgres:Password&1@localhost:5432",
				},
			},
		}}
	})
	JustBeforeEach(func() {
		dbmsPool = pool.NewDbmsPool(0)
		for _, dbms := range dbmsList {
			err = dbmsPool.RegisterDbms(dbms, database.Postgres)
		}
	})
	Context("when registering pool entries", func() {
		It("should not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})
		It("should register the postgres entry correctly", func() {
			driver := dbmsPool.Get("postgres")
			Expect(driver).ToNot(BeNil())
			Expect(driver.Ping()).To(Succeed())
		})
	})
	Context("when trying to register a duplicate pool entry", func() {
		BeforeEach(func() {
			dbmsList = []database.Dbms{{
				DatabaseClassName: "testDbc",
				Endpoints: []database.Endpoint{
					{
						Name: "duplicate",
						Dsn:  "postgres://postgres:Password&1@localhost:5432",
					},
					{
						Name: "duplicate",
						Dsn:  "postgres://postgres:Password&1@localhost:5432",
					},
				},
			}}
		})
		It("should return an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Endpoint names must be unique within the list of endpoints"))
		})
	})
	Context("when trying to register a malformed DSN", func() {
		BeforeEach(func() {
			dbmsList = []database.Dbms{{
				DatabaseClassName: "testDbc",
				Endpoints: []database.Endpoint{
					{
						Name: "postgres1",
						Dsn:  "malformed123",
					},
				},
			}}
		})
		It("should return an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})
	Context("when trying to register an unexisting DSN", func() {
		BeforeEach(func() {
			dbmsList = []database.Dbms{{
				DatabaseClassName: "testDbc",
				Endpoints: []database.Endpoint{
					{
						Name: "postgres1",
						Dsn:  "postgres://postgres:fakepassword@fakehostname:12345",
					},
				},
			}}
		})
		It("should return an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})
})
