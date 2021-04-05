// Package config contains helper methods used to retrieve configuration about the operator. It is a wrapper around the
// Viper library.
package config

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/spf13/viper"
)

var (
	dbms database.DbmsList
)

// Get returns the dbms endpoint list of the operator as specified in the operator's configuration file.
func GetDbmsList() (database.DbmsList, error) {
	if dbms == nil {
		if err := viper.UnmarshalKey(database.DbmsConfigKey, &dbms); err != nil {
			return nil, err
		}
	}
	return dbms, nil
}
