package dbapi

import (
	"fmt"
	"rtap/rtdbms"
	"strings"
)

var RTDBMS *rtdbms.RTDBMS

// ----------------------------------------------------------------------------
//
//	realtime database functions
//
// ----------------------------------------------------------------------------
type dbapi_rt struct {
	db_name string
}

func (db *dbapi_rt) Init(conn_str string) error {
	db.db_name = conn_str
	return nil
}

func (db *dbapi_rt) Get(schema_name string, table_name string, columns *[]string,
	where string, orderby *[]string, out *[]any) error {

	// Get a pointer to the table
	table, err := RTDBMS.GetTableByName(db.db_name, schema_name, table_name)
	if err != nil {
		return err
	}

	//extract the row_name from the where clause
	splits := strings.Split(where, "=")
	if splits[0] != "rowname" {
		return fmt.Errorf("invalid column name")
	}

	// Read the requested columns
	return table.ReadRow(splits[1], columns, out)

}
