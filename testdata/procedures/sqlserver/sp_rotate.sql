CREATE OR ALTER PROCEDURE sp_rotate @k8sName varchar(max)
AS 
DECLARE @sql varchar(max)
DECLARE @newpwd varchar(255)
DECLARE @ttimestamp datetime = getDate()
DECLARE @fqdn varchar(max)
DECLARE @port varchar(max)
DECLARE @username varchar(max)

SELECT @username = databases.username, @fqdn = databases.fqdn, @port = databases.port 
FROM databases 
WHERE databases.dbName = @k8sName

SET @newpwd = CONVERT(varchar(255), NEWID())

DECLARE @t TABLE([key] varchar(max), value varchar(max))
INSERT @t VALUES('username', @username)
INSERT @t VALUES('password', @newpwd)
INSERT @t VALUES('dbName', @k8sName)
INSERT @t VALUES('fqdn', @fqdn)
INSERT @t VALUES('port', @port)
INSERT @t VALUES('lastRotation', @ttimestamp)

UPDATE databases SET password = @newpwd WHERE dbname = @k8sName

SELECT * FROM @t