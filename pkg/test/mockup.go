package test

import (
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"strconv"
)

const (
	// Stored procedure values (should match the stored procedure signature on the DBMS side)
	opSpCreateName = "sp_create"
	opSpDeleteName = "sp_delete"
	opInputK8sName = "k8sname"
	opOutputUser   = "username"
	opOutputPass   = "password"
	opOutputDbName = "dbName"

	endptName = "myendpoint"

	MockK8sName      = "281cb030-9d56-48fa-b433-57237c9e1188"
	MockOutputUser   = "testuser"
	MockOutputPass   = "testpass"
	MockOutputDbName = "mydbtest281cb030" // DbNamePrefix + the first 8 chars of MockK8sName

	// Testing stored procedure-specific values
	DbNamePrefix   = "mydbtest"
	LeftTrimLength = "8"

	// TODO: Check that stored procedure is idempotent, i.e. it returns the correct username and password depending on
	// the k8sName, and not always the same
	CreateSp = `create proc sp_create (@` + opInputK8sName + ` varchar(max), @` + opOutputUser +
		` varchar(max) OUTPUT, @` + opOutputPass + ` varchar(max) OUTPUT, @` + opOutputDbName + ` varchar(max) OUTPUT)
					as
					declare @sql varchar(max)
					declare @tempDbName varchar(max)	
					set @tempDbName=CONCAT('` + DbNamePrefix + `', (SELECT LEFT(CONVERT(varchar(255), @k8sName),` + LeftTrimLength + `)))
					IF COUNT((DB_ID(@tempDbName))) = 0 
						BEGIN
							set @sql = CONCAT('create database ',@tempDbName)
							exec (@sql)
						END
					select @username='` + MockOutputUser + `',@password='` + MockOutputPass + `',@dbName='` + MockOutputDbName + `'`

	DeleteSp = `create proc sp_delete (@` + opInputK8sName + ` varchar(max))
					as
					declare @sql varchar(max)
					declare @dbname varchar(255)
					set @dbname=CONCAT('` + DbNamePrefix + `', (SELECT LEFT(CONVERT(varchar(255), @k8sName),8)))
					IF COUNT((DB_ID(@dbname))) > 0 
						BEGIN
							set @sql = CONCAT('drop database ',@dbname)
							exec (@sql)
						END`
)

// GetMockOps returns a mock of map[string]database.Operation
func GetMockOps() map[string]database.Operation {
	// Mock map[string]Operation
	ops := map[string]database.Operation{}
	inputs := map[string]string{}
	outputs := map[string]string{}
	inputs[database.K8sMapKey] = opInputK8sName
	outputs[database.UserMapKey] = opOutputUser
	outputs[database.PassMapKey] = opOutputPass
	outputs[database.DbNameMapKey] = opOutputDbName
	ops[database.CreateMapKey] = database.Operation{
		Name:    opSpCreateName,
		Inputs:  inputs,
		Outputs: outputs,
	}
	ops[database.DeleteMapKey] = database.Operation{
		Name:   opSpDeleteName,
		Inputs: inputs,
	}
	return ops
}

// GetMockEndpoints returns a mock of []database.Endpoint given a slice of database.Dsn
func GetMockEndpoints(dsnSlice []database.Dsn) []database.Endpoint {
	var endpts []database.Endpoint

	for i, dsn := range dsnSlice {
		endpts = append(endpts, database.Endpoint{
			Name: endptName + strconv.Itoa(i),
			Dsn:  dsn,
		})
	}
	return endpts
}
