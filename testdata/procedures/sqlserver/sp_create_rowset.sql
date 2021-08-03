create proc sp_create_rowset (@k8sName varchar(max))
as
declare @sql varchar(max)
declare @tempDbName varchar(255)
set @tempDbName=CONCAT('mydbtest', (SELECT LEFT(CONVERT(varchar(255), @k8sName),8)))

IF COUNT((DB_ID(@tempDbName))) = 0
    BEGIN
        set @sql = CONCAT('create database ',@tempDbName)
        exec (@sql)
    END


declare @t table(username varchar(max), password varchar(max), dbName varchar(max), fqdn varchar(max), port varchar(max))
insert @t values('testuser', 'testpass', @tempDbName, 'localhost', '1433')
    
select * from @t