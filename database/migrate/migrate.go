package migrate

var Stmt = `CREATE TABLE IF NOT EXISTS ohcl (
		unix bigint,
		symbol text,
		open double precision,
		high double precision,
		low double precision,
		close double precision);`
