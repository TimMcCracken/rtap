package rtdsms

import (
	"fmt"
	"gorm.io/gorm"
	"math"
//		"os"
//		"reflect"
//		"regexp"
//		"runtime"
	"github.com/glebarez/sqlite"
	"time"
)







/* ----------------------------------------------------------------------------
Table
-----------------------------------------------------------------------------*/
type table struct {
	descriptor		*TableDescriptor
//	path			string
	lastSelect		time.Time
	lastInsert		time.Time
	lastUpdate		time.Time
	lastDelete		time.Time
	selectCounter	uint64
	insertCounter	uint64
	updateCounter	uint64
	deleteCounter	uint64
	rowsEmpty		uint64
	rowDirty		[]bool
	rowEmpty		[]bool
	columns			[]*column

	// The following fields provide a doubly-linked list of all the empty rows
	// in the table
	firstEmptyRow			uint64		// list head
	lastEmptyRow			uint64		// list tail
	nextEmptyRow			[]uint64	// one entry per row
	prevEmptyRow			[]uint64	// one entry per row

}


func (table *table) construct (filename string)  error {
	
//	fmt.Printf("Constructing table [%s] [%s] [%s] [%s] [%s]\n", 
//				table.descriptor.realm_name, 
//				table.descriptor.domain_name, 
//				table.descriptor.datastore_name,
//				table.descriptor.schemaName,
//				table.descriptor.Table_name)


	if table.descriptor.columnsMap == nil {
		table.descriptor.columnsMap = make(map[string]uint16)
	}

	// ------------------------------------------------------------------------
	// Open the sqlite file
	// ------------------------------------------------------------------------
	db, err := gorm.Open(sqlite.Open(filename + "?_pragma=foreign_keys(1)" ), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database using path [%s]", filename)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Internal error: failed to obtain database.DB in CreateDomain()")
	}
	defer sqlDB.Close()

	// ------------------------------------------------------------------------
	// Get the domain_id
	// ------------------------------------------------------------------------
	var domain_id	int64
	result := db.Raw("SELECT domain_id FROM domains where domain_name = ?", 
					table.descriptor.domainName).Scan(&domain_id)
	if result.Error != nil {
		return result.Error
	}

	// ------------------------------------------------------------------------
	// Get the datastore_id_id
	// ------------------------------------------------------------------------
	var datastore_id	int64
	result = db.Raw("SELECT datastore_id FROM datastores WHERE domain_id = ? AND datastore_name = ?", 
					domain_id, table.descriptor.datastoreName).Scan(&datastore_id)
	if result.Error != nil {
		return result.Error
	}

	// ------------------------------------------------------------------------
	// Get the snapshot count
	// ------------------------------------------------------------------------
	var snapshot_count	int
	result = db.Raw("SELECT count(*) FROM snapshots WHERE domain_id = ? AND datastore_id = ?", 
					domain_id, datastore_id).Scan(&snapshot_count)
	if result.Error != nil {
		return result.Error
	}

	
	// ------------------------------------------------------------------------
	// Get the schema_id
	// ------------------------------------------------------------------------
	var schema_id	int64
	result = db.Raw("SELECT schema_id FROM schemas WHERE domain_id = ? AND datastore_id = ? AND schema_name = ?", 
					domain_id, datastore_id, table.descriptor.schemaName).Scan(&schema_id)
	if result.Error != nil {
		return result.Error
	}

	// ------------------------------------------------------------------------
	// Get the table_id, etc/
	// ------------------------------------------------------------------------
	var table_id	int64
	result = db.Raw("SELECT table_id FROM tables WHERE domain_id = ? AND " + 
					"datastore_id = ? AND schema_id = ? AND table_name = ?", 
					domain_id, datastore_id, schema_id,
					table.descriptor.tableName).Scan(&table_id)
	if result.Error != nil {
		return result.Error
	}

	table.rowDirty 		= make([]bool, table.descriptor.totalRows, table.descriptor.totalRows)
	table.rowEmpty 		= make([]bool, table.descriptor.totalRows, table.descriptor.totalRows)
	table.rowsEmpty		= table.descriptor.totalRows
	table.nextEmptyRow 	= make([]uint64, table.descriptor.totalRows, table.descriptor.totalRows)
	table.prevEmptyRow 	= make([]uint64, table.descriptor.totalRows, table.descriptor.totalRows)
	
	table.firstEmptyRow = 0	// not really necessary
	table.lastEmptyRow = table.descriptor.totalRows -1

	for i := uint64(0) ; i < table.descriptor.totalRows; i++ {
		table.nextEmptyRow[i] = i + 1
		table.prevEmptyRow[i] = math.MaxUint64
	}
	table.nextEmptyRow[table.descriptor.totalRows - 1] = math.MaxUint64


	// set all rows to 'empty = true'
	for i := uint64(0); i < table.descriptor.totalRows; i++ {
		table.rowEmpty[i] = true
	}

	// ------------------------------------------------------------------------
	// Get the columns for this table
	// ------------------------------------------------------------------------	
	rows, err := sqlDB.Query("SELECT column_name, type_name, not_null, " +
							"unique_, indexed, no_persist, default_value, " +
							"check_, min_integer, max_integer, min_uint, " +
							"max_uint, min_float, max_float, min_length, " +
							"length FROM columns " +  
							"WHERE domain_id = ? AND datastore_id = ? AND " +
							"schema_id = ? AND table_id = ? ORDER BY column_id", 
							domain_id, datastore_id, schema_id, table_id)
	if err != nil {
		return err
	}
	defer rows.Close()

	// ------------------------------------------------------------------------
	// Loop through the columns, fetching from the .schema (sqlite3) file 
	// ------------------------------------------------------------------------
	for rows.Next() {
		var column_name 	string
		var type_name 		string
		var not_null		bool
		var unique			bool
		var indexed			bool
		var no_persist 		bool
		var default_value	string
		var	check			string
		var min_integer 	int64
		var max_integer 	int64
		var min_uint 		int64 //sqlite dose not support unsigned types
		var max_uint 		int64 //sqlite dose not support unsigned types
		var min_float 		float64
		var max_float 		float64
		var min_length 		uint64
		var length 			uint64

        err := rows.Scan(&column_name, &type_name, &not_null, &unique, 
						 &indexed, &no_persist, &default_value, &check, 
						 &min_integer, &max_integer, &min_uint, &max_uint, 
						 &min_float, &max_float, &min_length,  &length)
		if err != nil {
			return fmt.Errorf("Error constructing table: [%v]\n", err)
		} else {
			var options 	uint64
//			fmt.Printf("Column: [%s] [%s]\n", column_name, type_name,  )
			
			//-----------------------------------------------------------------
			// change the options back to a bitfield
			// -----------------------------------------------------------------
			if not_null 	{ options |= NOT_NULL}
			if unique 		{ options |= UNIQUE}
			if indexed 		{ options |= INDEXED}
			if no_persist 	{ options |= NO_PERSIST}
	

/*
type column struct {
	descriptor		*ColumnDescriptor
	realtime		any
	index			map[string]uint64
	snapshots		[][]any // data is copied here for disk storage and 
						  //replication purposes, capturing post disturbance analysis
						  // data, etc
	contexts		[][]any // for advanced apps that run simulations, etc
						 // and operate other tha on real time data
} */			

			cd, err := NewColumnDescriptor(column_name, type_name, options, length, min_length)
			if err != nil {
				return fmt.Errorf("Error constructing column: [%v]\n", err)
			}
			table.descriptor.ColumnDescriptors = append(table.descriptor.ColumnDescriptors, cd)

			var data any 
			var snapshots []any = make([]any, snapshot_count)


			switch (type_name) {
			case "varUTF8":
				data = make([]string, table.descriptor.totalRows, table.descriptor.totalRows)
				for i := 0; i < snapshot_count; i++  {
					snapshots[i] = make([]string, table.descriptor.totalRows, table.descriptor.totalRows)
				} 

		
			case "int64": // TODO: Add parameters
				data = make([]int64, table.descriptor.totalRows, table.descriptor.totalRows)
				for i := 0; i < snapshot_count; i++  {
					snapshots[i] = make([]int64, table.descriptor.totalRows, table.descriptor.totalRows)
				} 

			case "float64": // TODO: Add parameters
				data = make([]float64, table.descriptor.totalRows, table.descriptor.totalRows)
				for i := 0; i < snapshot_count; i++  {
					snapshots[i] = make([]float64, table.descriptor.totalRows, table.descriptor.totalRows)
				} 

			}

			column := column{
				realtime: data,
				snapshots: snapshots,
			}
			
			if indexed == true {
				column.index = make(map[string]uint64)					
			}
			
			table.descriptor.columnsMap[column_name] = uint16(len(table.columns))
			table.columns = append(table.columns, &column)

//			fmt.Printf("Columns: %v\n", len(table.descriptor.Column_descriptors))
		}
	}
	return nil
}



/*

func (tbl *table) PrintColumnDescriptors(){

	fmt.Println("")
	for i, col_desc := range tbl.descriptor.Column_descriptors {
		fmt.Printf(" %03d %16s %-8s\n", i, col_desc.name, col_desc.typeName)
	}
	fmt.Println("")

}
	*/


