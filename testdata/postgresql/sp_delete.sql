create or replace function sp_delete(k8sName text) returns void as $func$
begin

if exists (select 1 from pg_database where datname = k8sName) then
	perform dblink_exec('dbname=' || current_database(), 'DROP DATABASE ' || quote_ident(k8sName) || ' WITH (FORCE)');
end if;

end;
$func$ language plpgsql;