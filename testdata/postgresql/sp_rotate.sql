CREATE OR REPLACE FUNCTION sp_rotate(k8sName text) RETURNS TABLE("key" text, value text) AS $$
	DECLARE 
		new_pwd TEXT;
	BEGIN
		new_pwd := md5(random()::text);
		-- drop table or it will generate an error on subsequent executions if reusing connection pools
		DROP TABLE IF EXISTS t;
		CREATE TEMPORARY TABLE t("key" text, value text);
		INSERT INTO t VALUES('password', new_pwd);
		UPDATE databases SET password = new_pwd WHERE dbname = k8sName;
		RETURN query EXECUTE 'select * from t';
	END;
$$ LANGUAGE plpgsql;