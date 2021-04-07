// Package typeutil defines various enums and messages for conditions, events and logging purposes.
package typeutil

const (
	// Condition types
	TypeReady = "Ready"

	// UpperCamelCase reasons enumerable, generic format is <ReasonDescription>[Verb]<Outcome (e.g. "Ready", "InProgress", "Failed"...)>
	ReasonCreate                   = "DatabaseReady"
	ReasonDatabaseGetFailed        = "DatabaseGetFailed"
	ReasonDatabaseCreateFailed     = "DatabaseCreateFailed"
	ReasonDatabaseDeleteFailed     = "DatabaseDeleteFailed"
	ReasonDatabaseUpdateFailed     = "DatabaseUpdateFailed"
	ReasonDatabaseCreateInProgress = "DatabaseCreateInProgress"
	ReasonDbmsConfigReadFailed     = "DbmsConfigReadFailed"
	ReasonDatabaseClassGetFailed   = "DatabaseClassGetFailed"

	// Human-readable messages
	MessageCreateSuccess            = "Database instance provisioned successfully on DBMS endpoint"
	MessageDatabaseNotFound         = "Database resource not found. Ignoring since object must be deleted"
	MessageDatabaseGetFailed        = "Database resource get failed"
	MessageDatabaseCreateFailed     = "Failed to provision Database resource"
	MessageDatabaseDeleteFailed     = "Failed to delete Database resource"
	MessageDatabaseUpdateFailed     = "Error updating resource after removing finalizer"
	MessageDatabaseCreateInProgress = "Database instance is being provisioned on DBMS endpoint"
	MessageDbmsConfigReadFailed     = "Could not retrieve dbms list from operator config"
	MessageDatabaseClassGetFailed   = "Could not retrieve databaseclass name from dbms config: %s"

	// Event types
	Normal  = "Normal"
	Warning = "Warning"
)
