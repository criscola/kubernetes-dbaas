package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xo/dburl"
)

// PsqlDirectConn represents a connection to a Postgres DBMS.
type PsqlDirectConn struct {
	c *pgxpool.Pool
}

type CreationQuery struct {
	Query      string
	Parameters map[string]string
}

// NewPsqlDirectConn opens a new PostgreSQL connection from a given dsn.
func NewPsqlDirectConn(dsn string) (*PsqlDirectConn, error) {
	dbConn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	conn := PsqlDirectConn{dbConn}

	return &conn, nil
}

func normalizeRune(r rune) rune {
	if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-._", r) {
		return r
	}
	return '-'
}

func (c *PsqlDirectConn) getDatabaseParameters(operation Operation) map[string]string {

	dbParams := make(map[string]string)

	namespace := operation.Inputs["namespace"]
	ressourceName := operation.Inputs["name"]

	baseName := fmt.Sprintf("%s-%s", namespace, ressourceName)
	length := len(baseName)

	if length > 30 {
		length = 30
	}

	baseName = baseName[:length]
	password := c.genPassword()

	DSN, _ := dburl.Parse(operation.DSN)

	host := DSN.Host
	port := DSN.Port()

	dbParams["dbName"] = "d_" + baseName
	dbParams["fqdn"] = host
	dbParams["password"] = password
	dbParams["port"] = port
	dbParams["username"] = "u_" + baseName

	return dbParams
}

func (c *PsqlDirectConn) CreateDb(createOp Operation) OpOutput {
	newDbParameters := c.getDatabaseParameters(createOp)

	for k, v := range newDbParameters {
		createOp.Secrets[k] = v
	}

	creationQuery, err := c.getPsqlOpQuery(createOp)

	if err != nil {
		return OpOutput{nil, err}
	}

	creationError := queryMutliLine(creationQuery, c)

	result := newDbParameters
	result["dsn"] = createOp.DSN
	result["lastRotation"] = time.Now().Format(time.RFC3339)

	return OpOutput{result, creationError}
}

func (c *PsqlDirectConn) DeleteDb(deleteOp Operation) OpOutput {
	deleteQuery, err := c.getPsqlOpQuery(deleteOp)

	expectedDbParam := c.getDatabaseParameters(deleteOp)
	expectedDbName := expectedDbParam["dbName"]

	if expectedDbName != deleteOp.Secrets["dbName"] {
		return OpOutput{
			nil,
			fmt.Errorf("anomaly detected: expected data-basename to be %s, got %s, secrets migth have been tempered with. deletion aborted", expectedDbName, deleteOp.Secrets["dbName"]),
		}
	}

	if err != nil {
		return OpOutput{nil, err}
	}

	deletionError := queryMutliLine(deleteQuery, c)

	return OpOutput{nil, deletionError}
}

// Rotate attempts to rotate the credentials of a connection.
func (c *PsqlDirectConn) Rotate(rotate Operation) OpOutput {
	dbParams := c.getDatabaseParameters(rotate)

	for k, v := range dbParams {
		rotate.Secrets[k] = v
	}

	rotateQuery, err := c.getPsqlOpQuery(rotate)

	if err != nil {
		return OpOutput{nil, err}
	}

	rotationError := queryMutliLine(rotateQuery, c)

	rotate.Secrets["password"] = dbParams["password"]
	rotate.Secrets["lastRotation"] = time.Now().Format(time.RFC3339)

	return OpOutput{rotate.Secrets, rotationError}
}

// Ping returns an error if a connection cannot be established with the DBMS, else it returns nil.
func (c *PsqlDirectConn) Ping() error {
	return c.c.Ping(context.Background())
}

func (c *PsqlDirectConn) getPsqlOpQuery(operation Operation) (string, error) {
	template := operation.SqlTemplate

	psqlQuery, err := RenderGoTemplate(template, operation.Secrets, ErrorOnMissingKeyOption)

	if err != nil {
		return "", err
	}

	return psqlQuery, nil
}

func (c *PsqlDirectConn) genPassword() string {
	return uuid.New().String()
}

func queryMutliLine(queries string, psql *PsqlDirectConn) error {
	subQuery := sqlCommandSplitter(queries)

	batch := &pgx.Batch{}

	for _, query := range subQuery {
		batch.Queue(query)
	}

	return psql.c.SendBatch(context.Background(), batch).Close()
}

func sqlCommandSplitter(s string) []string {
	return strings.FieldsFunc(
		s,
		func(r rune) bool {
			return r == ';' || r == '\n'
		},
	)
}
