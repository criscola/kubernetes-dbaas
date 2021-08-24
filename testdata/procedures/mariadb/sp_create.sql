DELIMITER $
CREATE OR REPLACE PROCEDURE sp_create_db_rowset_eav(k8sName TEXT) 
BEGIN 
	
	EXECUTE IMMEDIATE CONCAT("CREATE DATABASE IF NOT EXISTS `", k8sName, "`"); 

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
	INSERT INTO t VALUES('lastRotation', '');

	IF NOT EXISTS(SELECT id FROM _databases WHERE dbName = k8sName) THEN
		INSERT INTO _databases (username, password, dbName, port, fqdn) VALUES('testuser', 'testpassword', k8sName, '3306', 'localhost');
	END IF;

	SELECT * FROM t;
END $
DELIMITER ;