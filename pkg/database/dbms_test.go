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
	operationNameBasicCorrect      = "sp_test"
)

func TestDatabase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books Suite")
}

var _ = Describe(FormatTestDesc(Unit, "RenderOperation"), func() {
	mapInputsBasicCorrect := make(map[string]string)
	mapOutputsBasicCorrect := make(map[string]string)

	Context("when Operation and OpValues are defined correctly", func() {
		// Prepare test data
		// Prepare correct inputs for templated operation to test
		mapInputsBasicCorrect = map[string]string{
			"name": "{{ .Metadata.name }}",
			"namespace": "{{ .Metadata.namespace }}",
			"department": "{{ .Parameters.department }}",
			"assignee": "{{ .Parameters.assignee }}",
		}

		// Prepare correct outputs for templated operation to test
		mapOutputsBasicCorrect = map[string]string{
			"password": "password",
			"username": "username",
			"dbName": "dbName",
			"fqdn": "fqdn",
			"port": "port",
		}

		operationStructTemplatedBasicCorrect := database.Operation{
			Name:    operationNameBasicCorrect,
			Inputs:  mapInputsBasicCorrect,
			Outputs: mapOutputsBasicCorrect,
		}

		// Prepare correct operation values for templated operation to test
		opValuesBasicCorrect := database.OpValues{
			Metadata: map[string]interface{}{
				"name": "TestDb",
				"namespace": "TestDbNamespace",
			},
			Parameters: map[string]string{
				"department": "development",
				"assignee":   "John Doe",
			},
		}

		// Prepare comparison data for assertions
		operationAssertionData := database.Operation{
			Name:    operationNameBasicCorrect,
			Inputs: map[string]string{
				"name": "TestDb",
				"namespace": "TestDbNamespace",
				"department": "development",
				"assignee":   "John Doe",
			},
			Outputs: mapOutputsBasicCorrect,
		}

		// Execute tested behavior
		renderedOperation, err := operationStructTemplatedBasicCorrect.RenderOperation(opValuesBasicCorrect)

		It("does not produce an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("renders an Operation correctly", func() {
			Expect(renderedOperation).To(Equal(operationAssertionData))
		})
	})
})
