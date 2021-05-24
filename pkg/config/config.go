// Package config contains helper methods used to retrieve configuration about the operator. It is a wrapper around the
// Viper library.
package config

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
)

var (
	dbms database.DbmsList
)
