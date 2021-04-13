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
	RsnOpValuesCreateFail   = "OperationValuesCreateFailed"
	RsnOpRenderFail         = "OperationRenderFailed"
	RsnDbmsEndpointConnFail = "DbmsEndpointConnectFailed"
	RsnDbDeleteFail			= "DatabaseDeleteFailed"
	RsnDbCreateFail			= "DatabaseCreateFailed"
	RsnSecretCreateFail 	= "SecretCreateFailed"

	// Human-readable messages
	MsgDbProvisionSucc      = "Database instance provisioned successfully on DBMS endpoint"
	MsgDbDeleted            = "Database resource not found. Ignoring since object must be deleted"
	MsgDbGetFail            = "Database resource get failed"
	MsgDbProvisionFail      = "Failed to provision Database resource"
	MsgDbUpdateFail         = "Could not update resource after removing finalizer"
	MsgDbCreateInProg       = "Database instance is being provisioned on DBMS endpoint"
	MsgDbmsConfigGetFail    = "Could not retrieve DBMS list from operator config"
	MsgDbcConfigGetFail     = "Could not retrieve DatabaseClass name from dbms config"
	MsgDbcGetFail           = "DatabaseClass resource get failed"
	MsgOpNotSupported       = "Operation is not supported for DatabaseClass"
	MsgDbDeleteInProg       = "Database instance is being deleted from DBMS endpoint"
	MsgOpValuesCreateFail   = "Could not get create new operation values from resource"
	MsgOpRenderFail         = "Could not render operation values"
	MsgDbmsEndpointConnFail = "Could not get endpoint from pool"
	MsgDbDeleteFail			= "Could not delete database instance from DBMS endpoint"
	MsgDbCreateFail			= "Could not create database instance on DBMS endpoint"
	MsgSecretCreateFail 	= "Could not create Secret resource for Database resource"

	// Event types
	Normal  = "Normal"
	Warning = "Warning"
)
