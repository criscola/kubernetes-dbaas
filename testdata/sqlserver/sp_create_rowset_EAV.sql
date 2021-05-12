create proc sp_create_rowset_EAV (@k8sName varchar(max))
as
declare @sql varchar(max)

IF COUNT((DB_ID(@k8sName))) = 0
    BEGIN
        set @sql = CONCAT('create database ',@k8sName)
        exec (@sql)
    END

declare @t table([key] varchar(max), value varchar(max))
insert @t values('username', 'testuser')
insert @t values('password', 'testpassword')
insert @t values('dbName', @k8sName)
insert @t values('fqdn', 'localhost')
insert @t values('port', '1433')

select * from @t