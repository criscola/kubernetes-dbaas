CREATE OR REPLACE FUNCTION sp_delete(k8sName text) RETURNS void AS $func$
BEGIN

IF EXISTS (SELECT 1 FROM pg_database WHERE datname = k8sName) THEN
	PERFORM dblink_exec('dbname=' || current_database(), 'DROP DATABASE ' || quote_ident(k8sName) || ' WITH (FORCE)');
END IF;

IF exists(SELECT id FROM databases WHERE dbName = k8sName) THEN
	DELETE FROM databases WHERE dbName = k8sName
END IF;

END;
$func$ LANGUAGE plpgsql;