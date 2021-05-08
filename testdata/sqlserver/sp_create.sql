create proc sp_create (@k8sName varchar(max), @username varchar(max) OUTPUT, @password varchar(max) OUTPUT, @dbName varchar(max) OUTPUT, @fqdn varchar(max) OUTPUT, @port varchar(max) OUTPUT)
as
declare @sql varchar(max)
declare @tempDbName varchar(255)
set @tempDbName=CONCAT('mydbtest', (SELECT LEFT(CONVERT(varchar(255), @k8sName),8)))

IF COUNT((DB_ID(@tempDbName))) = 0
    BEGIN
        set @sql = CONCAT('create database ',@tempDbName)
        exec (@sql)
    END
    
select @username='testuser', @password='testpass', @dbName=@tempDbName, @fqdn='localhost', @port='1433'