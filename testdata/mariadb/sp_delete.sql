DELIMITER $
CREATE OR REPLACE PROCEDURE sp_delete(k8sName text)
BEGIN
	IF (SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = k8sName) <> 0 THEN
		EXECUTE IMMEDIATE CONCAT("DROP DATABASE ", k8sName);
	END IF;
END $
DELIMITER ;