package config

import (
	"fmt"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	. "github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

// OperatorConfig represents the key-value configuration of the operator
type OperatorConfig map[string]database.DbmsConfig

const (
	DbmsMapKey = "dbms"
)

var c OperatorConfig

// ReadOperatorConfig unmarshalls the operator configuration from a viper.Viper struct into a private struct.
//
// See GetDbmsConfig.
func ReadOperatorConfig(v *Viper) error {
	if v.GetStringSlice(DbmsMapKey) == nil {
		return fmt.Errorf("dbms configuration is not present in %s", v.ConfigFileUsed())
	}
	// Disallow unexpected attributes
	decoderConfig := func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnused = true
	}
	// Unmarshal config into struct
	if err := v.Unmarshal(&c, decoderConfig); err != nil {
		fmt.Println("error decoding into config struct")
		return err
	}
	err := validateDbmsConfig(c)
	if err := validateDbmsConfig(c); err != nil {
		err = errors.Unwrap(err)
	}
	return err
}

// GetDbmsConfig returns the private struct containing the DBMS configuration
func GetDbmsConfig() database.DbmsConfig {
	return c[DbmsMapKey]
}

// GetRootProjectPath is a utility function which gets the root project path based on the location of the main.go file.
// Do not move main.go in a place other than the root project path.
func GetRootProjectPath() (string, error) {
	var p string
	// Walk up the FS tree until a path containing main.go is found
	currentPath, _ := os.Getwd()
	mainFile := "main.go"

	for {
		p = path.Join(currentPath, mainFile)
		// Look if the current path contains main.go
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			// path of main.go found
			break
		}

		// Walk up the FS tree
		i := strings.LastIndex(currentPath, "/")
		currentPath = currentPath[:i]

		if currentPath == "" {
			return "", fmt.Errorf("path not found")
		}
	}
	return currentPath, nil
}

// validateDbmsConfig validates the configuration of a DBMS configuration. It returns wrapped errors with the issues
// it was able to detect.
func validateDbmsConfig(config OperatorConfig) error {
	var err error

	// Check that there is at least 1 dbms configured
	if len(config[DbmsMapKey]) < 1 {
		return wrapError(err, "at least 1 dbms entry must be specified")
	}

	// Check that there is at least 1 endpoint specified for each dbms type
	for _, driver := range config[DbmsMapKey] {
		if len(driver.Endpoints) < 1 {
			err = wrapError(err, fmt.Sprintf("at least 1 endpoint must be specified for each '%s' entry",
				driver.Driver))
		}
		// Check that for each endpoint, Name and Dsn have been specified
		for _, endpoint := range driver.Endpoints {
			if !endpoint.IsNamePresent() {
				err = wrapError(err, fmt.Sprintf(" missing in '%s'", driver.Driver))
			}
			if !endpoint.IsDsnPresent() {
				err = wrapError(err, fmt.Sprintf("endpoint dsn missing in '%s'", driver.Driver))
			}
		}
		if _, ok := driver.Operations[database.CreateMapKey]; !ok {
			err = wrapError(err, fmt.Sprintf("create operation missing in '%s'", driver.Driver))
		}

		if _, ok := driver.Operations[database.DeleteMapKey]; !ok {
			err = wrapError(err, fmt.Sprintf("delete operation missing in '%s'", driver.Driver))
		}
	}

	return err
}

func wrapErrorVerbose(err error, msg string) error {
	if err == nil {
		return errors.New(msg)
	}
	return errors.Wrap(err, msg)
}

func wrapError(err error, msg string) error {
	if err == nil {
		return fmt.Errorf(msg)
	}
	return fmt.Errorf("%w, %s", err, msg)
}

func init() {
	c = OperatorConfig{}
}
