create proc sp_delete (@k8sName varchar(max))
as
declare @sql varchar(max)
declare @dbname varchar(255)
set @dbname=CONCAT('mydbtest', (SELECT LEFT(CONVERT(varchar(255), @k8sName),8)))
IF COUNT((DB_ID(@dbname))) > 0 
	BEGIN
		set @sql = CONCAT('drop database ',@dbname)
		exec (@sql)
	END