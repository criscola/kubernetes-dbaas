create or replace function sp_delete(k8sName text) returns void as $func$
begin
	
-- drop any active connection to the db to avoid error	
select pg_terminate_backend(pg_stat_activity.pid)
from pg_stat_activity
where pg_stat_activity.datname = 'target_db' and pid <> pg_backend_pid();

if exists (select 1 from pg_database where datname = k8sName) then
	perform dblink_exec('dbname=' || current_database(), 'DROP DATABASE ' || quote_ident(k8sName));
end if;
return;
end;
$func$ language plpgsql;