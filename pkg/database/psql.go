package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PsqlConn represents a connection to a SQL Server DBMS.
type PsqlConn struct {
	c *pgxpool.Pool
}

// NewPsqlConn opens a new PostgreSQL connection from a given dsn.
func NewPsqlConn(dsn string) (*PsqlConn, error) {
	dbConn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	conn := PsqlConn{dbConn}

	return &conn, nil
}

func normalizeRune(r rune) rune {
	if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-.", r) {
		return r
	}
	return '-'
}

func getDatabaseParameters(operation Operation) map[string]string {
	username := uuid.New().String()
	password := uuid.New().String()
	databaseCredentials := make(map[string]string)

	if databaseName, ok := operation.Inputs["databaseName"]; ok {
		databaseCredentials["PGDATABASE_CREATE"] = strings.Map(normalizeRune, fmt.Sprintf("%v-%v", operation.Inputs["namespace"], databaseName))
	} else {
		databaseCredentials["PGDATABASE_CREATE"] = strings.Map(normalizeRune, "d_"+username)
	}
	databaseCredentials["PGDATABASE"] = databaseCredentials["PGDATABASE_CREATE"] + "@" + operation.DSN.Hostname()
	databaseCredentials["PGHOST"] = operation.DSN.Hostname()
	databaseCredentials["PGPASSWORD"] = password
	databaseCredentials["PGPORT"] = operation.DSN.Port()
	databaseCredentials["PGUSER_CREATE"] = "u_" + username
	databaseCredentials["PGUSER"] = "u_" + username + "@" + operation.DSN.Hostname()
	fmt.Printf("%v\n", databaseCredentials)
	return databaseCredentials
}

// CreateDb attempts to create a new database as specified in the operation parameter. It returns an OpOutput with the
// result of the call.
func (c *PsqlConn) CreateDb(operation Operation) OpOutput {
	val := getPsqlOpQuery(operation)

	fmt.Printf("%v\n", val)
	rows, err := c.c.Query(context.Background(), val)
	if err != nil {
		return OpOutput{Result: nil, Err: err}
	}

	var key string
	var value string
	result := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(&key, &value)
		if err != nil {
			return OpOutput{nil, err}
		}
		result[key] = value
	}

	return OpOutput{result, nil}
}

// DeleteDb attempts to delete a database instance as specified in the operation parameter. It returns an OpOutput with the
// result of the call if present.
func (c *PsqlConn) DeleteDb(operation Operation) OpOutput {
	_, err := c.c.Exec(context.Background(), getPsqlVoidOpQuery(operation))
	if err != nil {
		return OpOutput{nil, err}
	}

	return OpOutput{}
}

// Rotate attempts to rotate the credentials of a connection.
func (c *PsqlConn) Rotate(operation Operation) OpOutput {
	val := getPsqlOpQuery(operation)
	rows, err := c.c.Query(context.Background(), val)
	if err != nil {
		return OpOutput{nil, err}
	}

	var key string
	var value string
	result := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(&key, &value)
		if err != nil {
			return OpOutput{nil, err}
		}
		result[key] = value
	}

	return OpOutput{result, nil}
}

// Ping returns an error if a connection cannot be established with the DBMS, else it returns nil.
func (c *PsqlConn) Ping() error {
	return c.c.Ping(context.Background())
}

func getPsqlOpQuery(operation Operation) string {
	template := strings.ReplaceAll(operation.Name, "{ { ", "{{")
	templateData := make(map[string]string)
	for k, v := range operation.Inputs {
		templateData[k] = v
	}
	for k, v := range getDatabaseParameters(operation) {
		templateData[k] = v
	}
	renderedCreateOperation, _ := RenderGoTemplate(template, templateData)
	return renderedCreateOperation
}

func getPsqlVoidOpQuery(operation Operation) string {
	return fmt.Sprintf("select %s(%s)", operation.Name, getPsqlInputs(operation.Inputs))
}

func getPsqlInputs(values map[string]string) string {
	if len(values) == 0 {
		return ""
	}
	var result string
	for k, v := range values {
		result = fmt.Sprintf("%s := '%s', %s", k, v, result) // params specified on reverse order on purpose as a sanity check for postgres
	}

	result = result[:len(result)-2]
	return result
}
