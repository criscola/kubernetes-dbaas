CREATE OR ALTER PROCEDURE sp_delete (@k8sName varchar(max))
AS
DECLARE @sql varchar(max)
IF COUNT((DB_ID(@k8sName))) > 0
BEGIN
	set @sql = CONCAT('drop database ',@k8sName)
	exec (@sql)
END

IF EXISTS (SELECT id FROM databases WHERE databases.dbName = @k8sName) 
BEGIN
	DELETE FROM databases WHERE dbName = @k8sName
END