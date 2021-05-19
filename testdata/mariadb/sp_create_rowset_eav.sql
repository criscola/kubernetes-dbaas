DELIMITER $
CREATE OR REPLACE PROCEDURE sp_create_db_rowset_eav(k8sName TEXT) 
BEGIN 
	IF (SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = k8sName) = 0 THEN
		EXECUTE IMMEDIATE CONCAT("CREATE DATABASE ", k8sName); 
	END IF;

	CREATE TEMPORARY TABLE t(
		`key` text, 
		value text
	);

	INSERT INTO t VALUES('username', 'testuser');
	INSERT INTO t VALUES('password', 'testpassword');
	INSERT INTO t VALUES('dbName', k8sName);
	INSERT INTO t VALUES('fqdn', 'localhost');
	INSERT INTO t VALUES('port', '3306');
	
	SELECT * FROM t;
END $
DELIMITER ;