// Package typeutil defines various enums and messages for conditions, events and logging purposes. It should stay
// backward-compatible because external systems (e.g. alerting systems) might rely upon those values.
package typeutil

const (
	// Condition types
	TypeReady = "Ready"

	// UpperCamelCase reasons enumerable, generic format is <Subject>[Verb]<Outcome (e.g. "Ready", "InProgress", "Failed"...)>
	RsnCreate               = "DatabaseReady"
	RsnDbGetFail            = "DatabaseGetFailed"
	RsnDbUpdateFail         = "DatabaseUpdateFailed"
	RsnDbCreateInProg       = "DatabaseCreateInProgress"
	RsnDbmsConfigGetFail    = "DbmsConfigGetFailed"
	RsnDbcConfigGetFail     = "DatabaseClassConfigGetFailed"
	RsnDbcGetFail           = "DatabaseClassGetFailed"
	RsnOpNotSupported       = "OperationNotSupported"
	RsnDbDeleteInProg       = "DatabaseDeleteInProg"
	RsnOpRenderFail         = "OperationRenderFailed"
	RsnDbmsEndpointNotFound = "DbmsEndpointConnectFailed"
	RsnDbDeleteFail         = "DatabaseDeleteFailed"
	RsnDbCreateFail         = "DatabaseCreateFailed"
	RsnSecretCreateFail     = "SecretCreateFailed"
	RsnDbMetaParseFail      = "DatabaseMetaParseFailed"
	RsnDbSpecParseFail      = "DatabaseSpecParseFailed"
	RsnDbCreateSucc         = "DatabaseCreateSuccess"
	RsnDbOpQueueSucc        = "DatabaseQueueSuccess"
	RsnDbmsConnFail         = "DbmsConnectionFailed"
	RsnSecretCreateSucc     = "SecretCreateSuccess"
	RsnSecretUpdateSucc     = "SecretUpdateSuccess"
	RsnSecretUpdateFail     = "SecretUpdateFailed"
	RsnSecretGetFail        = "SecretGetFailed"
	RsnReadyCondUpdateFail  = "ReadyConditionUpdateFailed"

	// Human-readable messages
	MsgDbCreateSucc         = "database instance provisioned successfully on dbms endpoint"
	MsgDbDeleted            = "database resource not found. Ignoring since object must be deleted"
	MsgDbGetFail            = "database resource get failed"
	MsgDbUpdateFail         = "could not update database resource, retrying"
	MsgDbCreateInProg       = "database instance is being provisioned on dbms endpoint"
	MsgDbmsConfigGetFail    = "could not retrieve dbms list from operator config"
	MsgDbcConfigGetFail     = "could not retrieve databaseclass name from dbms config"
	MsgDbcGetFail           = "databaseclass resource get failed"
	MsgOpNotSupported       = "operation is not supported for databaseclass"
	MsgDbDeleteInProg       = "database instance is being deleted from dbms endpoint"
	MsgOpRenderFail         = "could not render operation values"
	MsgDbmsEndpointNotFound = "dbms connection not found in pool of connections"
	MsgDbDeleteFail         = "could not delete database instance from dbms endpoint"
	MsgDbCreateFail         = "could not create database instance on dbms endpoint"
	MsgSecretCreateFail     = "could not create secret resource for database resource"
	MsgDbMetaParseFail      = "could not parse metadata field of database resource during operation values creation"
	MsgDbSpecParseFail      = "could not parse spec field of database resource during operation values creation"
	MsgDbOpQueueSucc        = "database operation queued successfully"
	MsgDbmsConnFail         = "could not establish connection to dbms endpoint"
	MsgSecretCreateSucc     = "secret created successfully"
	MsgSecretUpdateSucc     = "secret updated successfully"
	MsgSecretUpdateFail     = "secret update failed"
	MsgSecretGetFail        = "secret get failed"
	MsgReadyCondUpdateFail  = "could not update ready condition of resource"

	// Event types
	Normal  = "Normal"
	Warning = "Warning"
)
