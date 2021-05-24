package database_test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	. "github.com/bedag/kubernetes-dbaas/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

const (
	// Test data
	operationNameBasicCorrect = "sp_test"
)

func TestDatabase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database package suite")
}

var _ = Describe(FormatTestDesc(Unit, "RenderOperation"), func() {
	var mapTemplateBasicCorrect map[string]string

	Context("when Operation and OpValues are defined correctly", func() {
		// Prepare test data
		// Prepare correct template for templating operation to test
		mapTemplateBasicCorrect = map[string]string{
			"name":       "{{ .Metadata.name }}",
			"namespace":  "{{ .Metadata.namespace }}",
			"department": "{{ .Parameters.department }}",
			"assignee":   "{{ .Parameters.assignee }}",
		}

		operationStructTemplatedBasicCorrect := database.Operation{
			Name:   operationNameBasicCorrect,
			Inputs: mapTemplateBasicCorrect,
		}

		// Prepare correct operation values for templated operation to test
		opValuesBasicCorrect := database.OpValues{
			Metadata: map[string]interface{}{
				"name":      "TestDb",
				"namespace": "TestDbNamespace",
			},
			Parameters: map[string]string{
				"department": "development",
				"assignee":   "John Doe",
			},
		}

		// Prepare comparison data for assertions
		operationAssertion := database.Operation{
			Name: operationNameBasicCorrect,
			Inputs: map[string]string{
				"name":       "TestDb",
				"namespace":  "TestDbNamespace",
				"department": "development",
				"assignee":   "John Doe",
			},
		}

		// Execute tested behavior
		renderedOperation, err := operationStructTemplatedBasicCorrect.RenderOperation(opValuesBasicCorrect)

		It("does not produce an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("renders an Operation correctly", func() {
			Expect(renderedOperation).To(Equal(operationAssertion))
		})
	})
})

var _ = Describe(FormatTestDesc(Unit, "RenderSecretFormat"), func() {
	var mapInputsBasicCorrect map[string]string

	Context("when SecretFormat and OpOutput are defined correctly", func() {
		// Prepare test data
		// Prepare correct inputs for templated operation to test
		mapInputsBasicCorrect = map[string]string{
			"username": "{{ .Result.username }}",
			"password": "{{ .Result.password }}",
			"port":     "{{ .Result.port }}",
			"dbName":   "{{ .Result.dbName }}",
			"server":   "{{ .Result.fqdn }}",
			"dsn":      "sqlserver://{{ .Result.username }}:{{ .Result.password }}@{{ .Result.fqdn }}/{{ .Result.dbName }}",
		}

		createOpOutput := database.OpOutput{
			Result: map[string]string{
				"username": "sa",
				"password": "Password&1",
				"port":     "1433",
				"dbName":   "testDb",
				"fqdn":     "localhost",
			},
		}

		secretFormatAssertion := database.SecretFormat(map[string]string{
			"username": "sa",
			"password": "Password&1",
			"port":     "1433",
			"dbName":   "testDb",
			"server":   "localhost",
			"dsn":      "sqlserver://sa:Password&1@localhost/testDb",
		})

		// Execute tested behavior
		renderedSecretFormat, err := database.SecretFormat(mapInputsBasicCorrect).RenderSecretFormat(createOpOutput)

		It("does not produce an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("renders a SecretFormat correctly", func() {
			Expect(renderedSecretFormat).To(Equal(secretFormatAssertion))
		})
	})
})
