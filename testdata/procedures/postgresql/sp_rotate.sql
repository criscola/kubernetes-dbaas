CREATE OR REPLACE FUNCTION sp_rotate(k8sName text) RETURNS TABLE("key" text, value text) AS $func$
	DECLARE 
		username varchar(40);
		fqdn varchar(64);
		port varchar(5);
		new_pwd TEXT;
	BEGIN
		new_pwd := md5(random()::text);
		-- drop table or it will generate an error on subsequent executions if reusing connection pools
		DROP TABLE IF EXISTS t;
		CREATE TEMPORARY TABLE t("key" text, value text);
		INSERT INTO t VALUES('username', username);
		INSERT INTO t VALUES('password', new_pwd);
		INSERT INTO t VALUES('fqdn', fqdn);
		INSERT INTO t VALUES('port', port);
		UPDATE databases SET password = new_pwd WHERE dbname = k8sName;
		RETURN query EXECUTE 'select * from t';
	END;
$func$ LANGUAGE plpgsql;
