package dbapi

import (
	"database/sql"
)

type dbapi_sql struct {
	pool *sql.DB
	conn string
}

func (db *dbapi_sql) Init(connection string) error {
	db.conn = connection
	return nil
}

func (db *dbapi_sql) Get(schema string, table string,
	columns *[]string, where string, orderby *[]string, out *[]any) error {
	return nil
}
