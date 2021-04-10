// Package typeutil defines various enums and messages for conditions, events and logging purposes. It should stay
// backward-compatible because external systems (e.g. alerting systems) might rely upon those values.
package typeutil

const (
	// Condition types
	TypeReady = "Ready"

	// UpperCamelCase reasons enumerable, generic format is <Subject>[Verb]<Outcome (e.g. "Ready", "InProgress", "Failed"...)>
	RsnCreate                = "DatabaseReady"
	RsnDbGetFail             = "DatabaseGetFailed"
	RsnDbCreateFail          = "DatabaseCreateFailed"
	RsnDbDeleteFail          = "DatabaseDeleteFailed"
	RsnDbUpdateFail          = "DatabaseUpdateFailed"
	RsnDbCreateInProg        = "DatabaseCreateInProgress"
	RsnDbmsConfigGetFail     = "DbmsConfigGetFailed"
	RsnDbcConfigGetFail      = "DatabaseClassConfigGetFailed"
	RsnDbcGetFail            = "DatabaseClassGetFailed"
	RsnOperationNotSupported = "OperationNotSupported"
	RsnDbDeleteInProg		 = "DatabaseDeleteInProg"

	// Human-readable messages
	MsgDbProvisionSucc   = "Database instance provisioned successfully on DBMS endpoint"
	MsgDbDeleted         = "Database resource not found. Ignoring since object must be deleted"
	MsgDbGetFail         = "Database resource get failed"
	MsgDbProvisionFail   = "Failed to provision Database resource"
	MsgDbDeleteFail      = "Failed to delete Database resource"
	MsgDbUpdateFail      = "Error updating resource after removing finalizer"
	MsgDbCreateInProg    = "Database instance is being provisioned on DBMS endpoint"
	MsgDbmsConfigGetFail = "Could not retrieve DBMS list from operator config"
	MsgDbcConfigGetFail  = "Could not retrieve DatabaseClass name from dbms config"
	MsgDbcGetFail        = "DatabaseClass resource get failed"
	MsgOperationNotSupported = "Operation not supported for DatabaseClass"
	MsgDbDeleteInProg    = "Database instance is being deleted from DBMS endpoint"

	// Event types
	Normal  = "Normal"
	Warning = "Warning"
)
