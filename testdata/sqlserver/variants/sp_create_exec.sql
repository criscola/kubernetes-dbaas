declare @usr varchar(255)
declare @pwd varchar(255)
declare @name varchar(255)
declare @myfqdn varchar(255)
declare @myport varchar(255)

EXEC sp_create @k8sName="mytest", @resourceUID="myid", @username=@usr, @password=@pwd, @dbName=@name, @fqdn=@myfqdn, @port=@myport