package config_test

import (
	"bytes"
	"github.com/bedag/kubernetes-dbaas/pkg/config"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/bedag/kubernetes-dbaas/pkg/test"
	"github.com/ghodss/yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"
)

const shouldGenerateError = "test should generate error condition"

var testdataFolderPath string

// Read config (test file)
// Compare with literal struct
// Test validation by providing various wrong inputs
func TestMain(m *testing.M) {
	temp, err := test.GetTestdataFolderPath()
	if err != nil {
		log.Fatalf("could not get testdata folder path: %s", err)
	}
	testdataFolderPath = path.Join(temp, "config")
	code := m.Run()

	os.Exit(code)
}

// TestReadOperatorConfig tests whether the 'given.yaml' test file corresponds to the 'expected.yaml' test file.
func TestReadOperatorConfig(t *testing.T) {
	// Setup
	viper.Reset()
	v := viper.GetViper()
	v.SetConfigType("yaml")

	// Read 'given' test file
	givenPath := path.Join(testdataFolderPath, "given.yaml")
	testInputFile, err := ioutil.ReadFile(givenPath)
	err = v.ReadConfig(bytes.NewReader(testInputFile))
	if err != nil {
		t.Fatalf("error reading config with viper: %s", err)
	}

	// Read 'expected' test file
	expected := database.DbmsConfig{}
	expectedPath := path.Join(testdataFolderPath, "expected.yaml")
	afterReadingConfig, err := ioutil.ReadFile(expectedPath)
	err = yaml.Unmarshal(afterReadingConfig, &expected)
	if err != nil {
		t.Errorf("error unmarshaling: %s", err)
	}

	// Execute tested behavior
	err = config.ReadOperatorConfig(v)
	if err != nil {
		t.Fatalf("error reading operator config: %s", err)
	}

	// Verify test
	got := config.GetDbmsConfig()
	if !cmp.Equal(got, expected) {
		t.Fatalf("Got %q expected %q", got, expected)
	}
}

// TestReadOperatorConfigMissingEndpointDsn tests if the config validation works by supplying
// 'given_wrong_missing_endpointdsn.yaml' to the config package and checking if it returns an error as it should.
func TestReadOperatorConfigMissingEndpointDsn(t *testing.T) {
	got := testWrongInputConfig(t, "given_wrong_missing_endpointdsn.yaml")
	if got == nil {
		t.Fatal(shouldGenerateError)
	}
}

// TestReadOperatorConfigMissingEndpointDsn tests if the config validation works by supplying
// 'given_wrong_missing_endpointname.yaml' to the config package and checking if it returns an error as it should.
func TestReadOperatorConfigMissingEndpointName(t *testing.T) {
	got := testWrongInputConfig(t, "given_wrong_missing_endpointname.yaml")
	if got == nil {
		t.Fatal(shouldGenerateError)
	}
}

// TestReadOperatorConfigMissingEndpointDsn tests if the config validation works by supplying
// 'given_wrong_missing_createop.yaml' to the config package and checking if it returns an error as it should.
func TestReadOperatorConfigMissingCreateOperation(t *testing.T) {
	got := testWrongInputConfig(t, "given_wrong_missing_createop.yaml")
	if got == nil {
		t.Fatal(shouldGenerateError)
	}
}

// TestReadOperatorConfigMissingEndpointDsn tests if the config validation works by supplying
// 'given_wrong_missing_deleteop.yaml' to the config package and checking if it returns an error as it should.
func TestReadOperatorConfigMissingDeleteOperation(t *testing.T) {
	got := testWrongInputConfig(t, "given_wrong_missing_deleteop.yaml")
	if got == nil {
		t.Fatal(shouldGenerateError)
	}
}

// testWrongInputConfig is a helper function which resets the Viper configuration and reads a 'given_wrong_[..].yaml'
// test file.
func testWrongInputConfig(t *testing.T, inputPath string) error {
	// Setup
	viper.Reset()
	v := viper.GetViper()
	v.SetConfigType(inputPath[:strings.LastIndex(inputPath, ".")])

	// Read 'given' test file
	givenPath := path.Join(testdataFolderPath, inputPath)
	testInputFile, err := ioutil.ReadFile(givenPath)
	testInput := bytes.NewReader(testInputFile)

	// Execute tested behaviour
	err = v.ReadConfig(testInput)
	if err != nil {
		t.Fatalf("error reading config with viper: %s", err)
	}

	// Verify test
	got := config.ReadOperatorConfig(v)
	if got != nil {
		t.Log(got)
	}
	return got
}
