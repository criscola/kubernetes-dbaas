package database_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	// Test data
	operationNameBasicCorrect = "sp_test"
)

var _ = Describe(FormatTestDesc(Unit, "RenderOperation"), func() {
	var templateOpInputs map[string]string
	var templatedOp database.Operation
	var opValues database.OpValues
	var renderedOperation database.Operation
	var operationAssertion database.Operation
	var err error

	BeforeEach(func() {
		// Prepare test data
		// Prepare correct template for templating operation to test
		templateOpInputs = map[string]string{
			"name":       "{{ .Metadata.name }}",
			"namespace":  "{{ .Metadata.namespace }}",
			"department": "{{ .Parameters.department }}",
			"assignee":   "{{ .Parameters.assignee }}",
		}
		templatedOp = database.Operation{
			Name:   operationNameBasicCorrect,
			Inputs: templateOpInputs,
		}
		// Prepare correct operation values for templated operation to test
		opValues = database.OpValues{
			Metadata: map[string]interface{}{
				"name":      "TestDb",
				"namespace": "TestDbNamespace",
			},
			Parameters: map[string]string{
				"department": "development",
				"assignee":   "John Doe",
			},
		}
	})
	JustBeforeEach(func() {
		// Prepare comparison data for assertions
		operationAssertion = database.Operation{
			Name: operationNameBasicCorrect,
			Inputs: map[string]string{
				"name":       "TestDb",
				"namespace":  "TestDbNamespace",
				"department": "development",
				"assignee":   "John Doe",
			},
			Secrets: make(map[string]string),
			DSN:     "",
		}
		// Execute tested behavior
		renderedOperation, err = templatedOp.RenderOperation(opValues)
	})
	Context("when Operation and OpValues are defined correctly", func() {
		It("should not produce an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("should render an Operation correctly", func() {
			Expect(renderedOperation).To(Equal(operationAssertion))
		})
	})
	Context("when an input to be rendered is specified but its value was not supplied", func() {
		BeforeEach(func() {
			By("not supplying the required value for 'assignee'")
			opValues = database.OpValues{
				Metadata: map[string]interface{}{
					"name":      "TestDb",
					"namespace": "TestDbNamespace",
				},
				Parameters: map[string]string{
					"department": "development",
				},
			}
		})
		It("should generate an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("map has no entry for key"))
		})
	})
})

var _ = Describe(FormatTestDesc(Unit, "RenderSecretFormat"), func() {
	var templateOpInputs map[string]string
	var createOpOutput database.OpOutput
	var renderedSecretFormat database.SecretFormat
	var secretFormatAssertion database.SecretFormat
	var err error

	BeforeEach(func() {
		// Prepare test data
		// Prepare correct inputs for templated operation to test
		templateOpInputs = map[string]string{
			"username": "{{ .Result.username }}",
			"password": "{{ .Result.password }}",
			"port":     "{{ .Result.port }}",
			"dbName":   "{{ .Result.dbName }}",
			"server":   "{{ .Result.fqdn }}",
			"dsn":      "sqlserver://{{ .Result.username }}:{{ .Result.password }}@{{ .Result.fqdn }}/{{ .Result.dbName }}",
		}
		createOpOutput = database.OpOutput{
			Result: map[string]string{
				"username": "sa",
				"password": "Password&1",
				"port":     "1433",
				"dbName":   "testDb",
				"fqdn":     "localhost",
			},
		}
	})
	JustBeforeEach(func() {
		secretFormatAssertion = map[string]string{
			"username": "sa",
			"password": "Password&1",
			"port":     "1433",
			"dbName":   "testDb",
			"server":   "localhost",
			"dsn":      "sqlserver://sa:Password&1@localhost/testDb",
		}
		// Execute tested behavior
		renderedSecretFormat, err = database.SecretFormat(templateOpInputs).RenderSecretFormat(createOpOutput)
	})

	Context("when SecretFormat and OpOutput are defined correctly", func() {
		It("does not produce an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("renders a SecretFormat correctly", func() {
			Expect(renderedSecretFormat).To(Equal(secretFormatAssertion))
		})
	})
	Context("when the return param was specified in SecretFormat but it is not returned by dbms", func() {
		By("not supplying 'fqdn'")
		BeforeEach(func() {
			createOpOutput = database.OpOutput{
				Result: map[string]string{
					"username": "sa",
					"password": "Password&1",
					"port":     "1433",
					"dbName":   "testDb",
				},
			}
		})
		It("should return an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("map has no entry for key"))
		})
	})
})
