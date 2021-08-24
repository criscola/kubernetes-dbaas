CREATE OR REPLACE FUNCTION sp_rotate(k8sname text)
 RETURNS TABLE(key text, value text)
 LANGUAGE plpgsql
AS $function$
	DECLARE 
		username varchar(40);
		fqdn varchar(64);
		port varchar(5);
		new_pwd TEXT;
		ttimestamp TEXT;
	BEGIN
		new_pwd := md5(random()::text);
		-- drop table or it will generate an error on subsequent executions if reusing connection pools
		DROP TABLE IF EXISTS t;
		CREATE TEMPORARY TABLE t("key" text, value text);
	
		SELECT databases.username, databases.fqdn, databases.port INTO username, fqdn, port FROM databases WHERE databases.dbName = k8sName;
		SELECT to_char(CURRENT_TIMESTAMP(8), 'DD Mon YYYY HH:MI:SSPM') INTO ttimestamp;
		INSERT INTO t VALUES('username', username);
		INSERT INTO t VALUES('password', new_pwd);
		INSERT INTO t VALUES('dbName', k8sName);
		INSERT INTO t VALUES('fqdn', fqdn);
		INSERT INTO t VALUES('port', port);
		INSERT INTO t VALUES('lastRotation', ttimestamp);
		UPDATE databases SET password = new_pwd WHERE dbname = k8sName;
		RETURN query EXECUTE 'select * from t';
	END;
$function$
;
