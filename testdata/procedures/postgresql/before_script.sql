create extension if not exists dblink;
CREATE TABLE databases (
	id 			SERIAL PRIMARY KEY,
	username	VARCHAR(40),
	password	TEXT NOT NULL,
	dbName		VARCHAR(63) UNIQUE NOT NULL,
	port		VARCHAR(5) NOT NULL,
	fqdn      	VARCHAR(64) NOT NULL
)