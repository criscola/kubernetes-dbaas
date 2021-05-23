DELIMITER $
CREATE OR REPLACE PROCEDURE sp_create_db_rowset_eav(k8sName TEXT) 
BEGIN 
	
	EXECUTE IMMEDIATE CONCAT("CREATE DATABASE IF NOT EXISTS ", k8sName); 

	DROP TEMPORARY TABLE IF EXISTS t;
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