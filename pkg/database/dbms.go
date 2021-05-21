package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"text/template"
)

const (
	Sqlserver               = "sqlserver"
	Postgres                = "postgres"
	Mysql                   = "mysql"
	Mariadb                 = "mariadb"
	CreateMapKey            = "create"
	DeleteMapKey            = "delete"
	RotateMapKey            = "rotate"
	OperationsConfigKey     = "operations"
	ErrorOnMissingKeyOption = "missingkey=error"
	DbmsConfigKey           = "dbms"
)

// Driver represents a struct responsible for executing CreateDb and DeleteDb operations on a system it supports. Drivers
// should provide a way to check their current status (i.e. whether it can accept CreateDb and DeleteDb operations at the
// moment of a Ping call
type Driver interface {
	CreateDb(operation Operation) OpOutput
	DeleteDb(operation Operation) OpOutput
	Rotate(operation Operation) OpOutput
	Ping() error
}

// DbmsConn represents the DBMS connection. See Driver.
type DbmsConn struct {
	Driver
}

// +kubebuilder:object:generate=true
// Operation represents an operation performed on a DBMS identified by name and containing a map of inputs and a map
// of outputs.
type Operation struct {
	Name    string            `json:"name,omitempty"`
	Inputs  map[string]string `json:"inputs,omitempty"`
}

// OpOutput represents the return values of an operation. If the operation generates an error, it must be set in the Err
// field. If Err is nil, the operation is assumed to be successful.
type OpOutput struct {
	Result map[string]string
	Err    error
}

// OpValues represent the input values of an operation.
type OpValues struct {
	Metadata   map[string]interface{}
	Parameters map[string]string
}

// +kubebuilder:object:generate=true
// Dbms is the instance associated with a Dbms resource. It contains the Driver responsible for the Operations executed on
// Endpoints.
type Dbms struct {
	DatabaseClassName string     `json:"databaseClassName"`
	Endpoints         []Endpoint `json:"endpoints"`
}

// +kubebuilder:object:generate=true
// DbmsList is a slice containing Dbms structs.
type DbmsList []Dbms

// +kubebuilder:object:generate=true
// +kubebuilder:kubebuilder:validation:MinItems=1
// Endpoint represent the configuration of a DBMS endpoint identified by a name.
type Endpoint struct {
	Name string `json:"name"`
	Dsn  Dsn    `json:"dsn"`
}

// +kubebuilder:object:generate=true
type SecretFormat map[string]string

// New initializes a Dbms instance based on a map of Operation. It expects a dsn like that:
// driver://username:password@host/instance?param1=value&param2=value
//
// See the individual Driver implementations.
func New(driver string, dsn Dsn) (*DbmsConn, error) {
	var dbmsConn *DbmsConn

	switch driver {
	case Sqlserver:
		parsedDsn, err := dsn.GenSqlserver()
		if err != nil {
			return nil, formatDsnParseError(err)
		}
		sqlserverConn, err := NewSqlserverConn(parsedDsn)
		if err != nil {
			return nil, err
		}
		dbmsConn = &DbmsConn{sqlserverConn}
	case Postgres:
		parsedDsn, err := dsn.GenPostgres()
		if err != nil {
			return nil, formatDsnParseError(err)
		}
		psqlConn, err := NewPsqlConn(parsedDsn)
		if err != nil {
			return nil, err
		}
		dbmsConn = &DbmsConn{psqlConn}
	case Mysql, Mariadb:
		parsedDsn, err := dsn.GenMysql()
		if err != nil {
			return nil, formatDsnParseError(err)
		}
		mysqlConn, err := NewMysqlConn(parsedDsn)
		if err != nil {
			return nil, err
		}
		dbmsConn = &DbmsConn{mysqlConn}
	default:
		return nil, fmt.Errorf("invalid dsn '%s': driver not found", dsn)
	}

	if err := dbmsConn.Ping(); err != nil {
		return nil, err
	}

	return dbmsConn, nil
}

// RenderOperation renders "actions" specified through the use of the Go text/template format. It renders Input of
// the receiver. Data to be inserted is taken directly from values. See OpValues. If the rendering is successful, the
// method returns ah na rendered Operation, if an error is generated, it is returned along with an empty Operation struct.
// Keys which are specified but not found generate an error (i.e. no unreferenced keys are allowed).
func (op Operation) RenderOperation(values OpValues) (Operation, error) {
	// Transform map[string]string to a single json string
	operationTemplate, err := json.Marshal(op.Inputs)
	if err != nil {
		return Operation{}, err
	}
	renderedInputsString, err := RenderGoTemplate(string(operationTemplate), values, ErrorOnMissingKeyOption)
	if err != nil {
		return Operation{}, err
	}
	var renderedInputsMap map[string]string
	err = json.Unmarshal([]byte(renderedInputsString), &renderedInputsMap)
	if err != nil {
		return Operation{}, err
	}
	renderedOp := Operation{
		Name:    op.Name,
		Inputs:  renderedInputsMap,
	}

	return renderedOp, nil
}

// GetDatabaseClassNameByEndpointName performs a linear search on the receiver in search of endpointName. If an entry
// is found, it is returned.
func (c DbmsList) GetDatabaseClassNameByEndpointName(endpointName string) string {
	for _, dbms := range c {
		if contains(dbms.Endpoints, endpointName) {
			return dbms.DatabaseClassName
		}
	}
	return ""
}

// IsNamePresent return true if an endpoint name is not empty, else it returns false.
func (e Endpoint) IsNamePresent() bool {
	return e.Name != ""
}

// IsDsnPresent return true if an endpoint dsn is not empty, else it returns false.
func (e Endpoint) IsDsnPresent() bool {
	return e.Dsn != ""
}

// RenderSecretFormat renders a secret from OpOutput.Result create operation.
func (s SecretFormat) RenderSecretFormat(createOpOutput OpOutput) (SecretFormat, error) {
	// Transform map[string]string to a single json string
	stringInputs, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	renderedInputsString, err := RenderGoTemplate(string(stringInputs), createOpOutput, ErrorOnMissingKeyOption)
	if err != nil {
		return nil, err
	}
	var renderedInputsMap map[string]string
	err = json.Unmarshal([]byte(renderedInputsString), &renderedInputsMap)
	if err != nil {
		return nil, err
	}
	renderedSecretFormat := SecretFormat(renderedInputsMap)

	return renderedSecretFormat, nil
}

// From merges the receiver with toMerge. s values are overwritten by toMerge values when they have same key.
// New k-v pairs are added to the result. Empty strings in toMerge are ignored (not overwritten nor added to the result).
func (s SecretFormat) From(toMerge map[string]string) SecretFormat {
	for k, v := range toMerge {
		if v != "" {
			s[k] = v
		}
	}
	return s
}

// RenderGoTemplate takes the text to be parsed as a Go template and values to be rendered. For options see template.Option.
func RenderGoTemplate(text string, values interface{}, options ...string) (string, error) {
	// Setup the template to be rendered based on the inputs
	tmpl, err := template.New("gotmpl").Parse(text)
	if err != nil {
		return "", err
	}
	tmpl.Option(options...)
	// Create a new buffer for the rendering result
	renderedInputsBuf := bytes.NewBufferString("")
	// Render each templated value by taking the values from values
	err = tmpl.Execute(renderedInputsBuf, values)
	if err != nil {
		return "", err
	}
	return renderedInputsBuf.String(), nil
}

// contains is a very small utility function which returns true if s has been found in list.
func contains(list []Endpoint, s string) bool {
	for _, v := range list {
		if v.Name == s {
			return true
		}
	}
	return false
}

// getQueryInputs returns a slice of sql.Named(k, v) from values where k is the key and v is the value.
func getQueryInputs(values map[string]string) []interface{} {
	var inputParams []interface{}
	for k, v := range values {
		inputParams = append(inputParams, sql.Named(k, v))
	}
	return inputParams
}

func formatDsnParseError(err error) error {
	return fmt.Errorf("error parsing dsn: %s", err)
}