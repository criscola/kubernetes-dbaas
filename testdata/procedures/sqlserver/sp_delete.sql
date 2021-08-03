create proc sp_delete (@k8sName varchar(max))
as
declare @sql varchar(max)
IF COUNT((DB_ID(@k8sName))) > 0
	BEGIN
		set @sql = CONCAT('drop database ',@k8sName)
		exec (@sql)
	END