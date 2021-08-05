DELIMITER $
CREATE OR REPLACE PROCEDURE sp_rotate(k8sName TEXT) 
BEGIN 
	DECLARE newpwd CHAR(30);
	DECLARE username VARCHAR(40);
	DECLARE fqdn VARCHAR(64);
	DECLARE port VARCHAR(5);
	DECLARE ttimestamp VARCHAR(40);

	SET ttimestamp = (SELECT DATE_FORMAT(NOW(), '%m/%d/%Y %H:%i'));
	
	DROP TEMPORARY TABLE IF EXISTS t;
	CREATE TEMPORARY TABLE t(
		`key` text, 
		value text
	);
	SET newpwd = (SELECT LEFT(UUID(), 30));

	SELECT _databases.username, _databases.fqdn, _databases.port INTO @username, @fqdn, @port FROM _databases WHERE _databases.dbName = k8sName;

	INSERT INTO t VALUES('username', @username);
	INSERT INTO t VALUES('password', newpwd);
	INSERT INTO t VALUES('dbName', k8sName);
	INSERT INTO t VALUES('fqdn', @fqdn);
	INSERT INTO t VALUES('port', @port);
	INSERT INTO t VALUES('lastRotation', ttimestamp);
	
	UPDATE _databases SET password = newpwd WHERE dbName = k8sName;

	SELECT * FROM t;
END $
DELIMITER ;