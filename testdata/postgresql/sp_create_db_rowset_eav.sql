create extension if not exists dblink;

create or replace function sp_create_db_rowset_eav(k8sName text)
  returns table("key" text, value text) as
$func$
begin

if not exists (select 1 from pg_database where datname = k8sName) then
   perform dblink_exec('dbname=' || current_database(), 'CREATE DATABASE ' || quote_ident(k8sName));
end if;

-- drop table or it will generate an error on subsequent executions if reusing connection pools
drop table if exists t;
create temporary table t("key" text, value text);
insert into t values('username', 'testuser');
insert into t values('password', 'testpassword');
insert into t values('dbName', k8sName);
insert into t values('fqdn', 'localhost');
insert into t values('port', '5432');

INSERT INTO databases (username, password, dbname, fqdn, port) VALUES ('testuser', 'testpassword', k8sName, 'localhost', '5432');

return query execute 'select * from t';
end;
$func$ language plpgsql;