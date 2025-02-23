package rtdsms

import (

	"bytes"
	"encoding/binary"	
	"fmt"
//		"math"
	"os"
//	"reflect"
//		"regexp"
//		"runtime"
//	"strings"
//		"time"
)

/* -----------------------------------------------------------------------------
These functions perfrom various action that enable persistence of the real time
objects including tables, key-value groups, arraus and matrices and perhaps
others in the future.  

Some applications might also need to "snapshot" files...
------------------------------------------------------------------------------*/


// This structure is the data supplied for the peristence requests. Not all
// fileds are used for every request.

type persist struct {
	schema_name		string
	table_name		string
	snapshot_name	string		// name of the snapshot used as the target 
	context_name	string		// name of the context us as the target
	path_file		string		// path and file name

}

type takeSnapshotFunc func(tbl *table, snapshot_index uint16, col_index uint16) error
type restoreSnapshotFunc func(tbl *table, snapshot_index uint16, col_index uint16) error

type writeSnapshotFunc func(tbl *table, snapshot_index uint16, col_index uint16, file * os.File) error
type readSnapshotFunc func(tbl *table, snapshot_index uint16, col_index uint16, file * os.File) error

type clearSnapshotFunc func(tbl *table, snapshot_index uint16, col_index uint16) error
type clearRealtimeFunc func(tbl *table, col_index uint16) error


var takeSnapshotFuncs		[max_column_types]takeSnapshotFunc
var restoreSnapshotFuncs	[max_column_types]restoreSnapshotFunc
var	writeSnapshotFuncs 		[max_column_types]writeSnapshotFunc
var	readSnapshotFuncs 		[max_column_types]readSnapshotFunc
var clearSnapshotFuncs		[max_column_types]clearSnapshotFunc
var clearRealtimeFuncs		[max_column_types]clearRealtimeFunc



func init(){

	// Initialize the restoreSnapshotFuncs table
	for i := 0; i < max_column_types; i++{
		takeSnapshotFuncs[i] = takeSnapshotUndefined
	}
	takeSnapshotFuncs[COL_TYPE_VAR_UTF8] 	= takeSnapshotUTF8 
	takeSnapshotFuncs[COL_TYPE_BLOB_UTF8] 	= takeSnapshotUTF8
	takeSnapshotFuncs[COL_TYPE_I64] 		= takeSnapshotInt64
	takeSnapshotFuncs[COL_TYPE_F64] 		= takeSnapshotFloat64


	// Initialize the restoreSnapshotFuncs table
	for i := 0; i < max_column_types; i++{
		restoreSnapshotFuncs[i] = restoreSnapshotUndefined
	}
	restoreSnapshotFuncs[COL_TYPE_VAR_UTF8]		= restoreSnapshotUTF8
	restoreSnapshotFuncs[COL_TYPE_BLOB_UTF8]	= restoreSnapshotUTF8
	restoreSnapshotFuncs[COL_TYPE_I64] 			= restoreSnapshotInt64
	restoreSnapshotFuncs[COL_TYPE_F64] 			= restoreSnapshotFloat64


	// Initialize the writeSnapshotFuncs table
	for i := 0; i < max_column_types; i++{
		writeSnapshotFuncs[i] = writeSnapshotUndefined
	}
	writeSnapshotFuncs[COL_TYPE_VAR_UTF8]	= writeSnapshotVarUTF8
	writeSnapshotFuncs[COL_TYPE_BLOB_UTF8]	= writeSnapshotBlobUTF8
	writeSnapshotFuncs[COL_TYPE_I64]		= writeSnapshotInt64
	writeSnapshotFuncs[COL_TYPE_F64]		= writeSnapshotFloat64

	// Initialize the readSnapshotFuncs table
	for i := 0; i < max_column_types; i++{
		readSnapshotFuncs[i] = readSnapshotUndefined
	}
	readSnapshotFuncs[COL_TYPE_VAR_UTF8] 	= readSnapshotVarUTF8
	readSnapshotFuncs[COL_TYPE_BLOB_UTF8] 	= readSnapshotBlobUTF8
	readSnapshotFuncs[COL_TYPE_I64] 		= readSnapshotInt64
	readSnapshotFuncs[COL_TYPE_F64] 		= readSnapshotFloat64

	// Initialize the eraseSnapshotFuncs table
	for i := 0; i < max_column_types; i++{
		clearSnapshotFuncs[i] = clearSnapshotUndefined
	}
	clearSnapshotFuncs[COL_TYPE_VAR_UTF8] 	= clearSnapshotUTF8
	clearSnapshotFuncs[COL_TYPE_BLOB_UTF8] 	= clearSnapshotUTF8
	clearSnapshotFuncs[COL_TYPE_I64] 		= clearSnapshotInt64
	clearSnapshotFuncs[COL_TYPE_F64] 		= clearSnapshotFloat64

	// Initialize the eraseRealtimeFuncs table
	for i := 0; i < max_column_types; i++{
		clearRealtimeFuncs[i] = clearRealtimeUndefined
	}
	clearRealtimeFuncs[COL_TYPE_VAR_UTF8] 	= clearRealtimeUTF8
	clearRealtimeFuncs[COL_TYPE_BLOB_UTF8] 	= clearRealtimeUTF8
	clearRealtimeFuncs[COL_TYPE_I64] 		= clearRealtimeInt64
	clearRealtimeFuncs[COL_TYPE_F64] 		= clearRealtimeFloat64
}
	



func takeSnapshot (ds * datastore, request * request) {
	
	// -------------------------------------------------------------------------
	// Get a pointer to the request data
	// -------------------------------------------------------------------------
	persist, ok := request.data.(*persist)
	if ok == false {
		request.err = fmt.Errorf("internal error - 'persist' type assertion failed.") // random message
		return
	}


	// -------------------------------------------------------------------------
	// Get the schema_index, if any
	// -------------------------------------------------------------------------
	var schema_index uint64
	var table_index uint64
	if persist.schema_name != "*" {
		schema_index, ok = ds.schemas_map[persist.schema_name]
		if ok != true {
			request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
			return
		}

		// --------------------------------------------------------------------
		// Get the table_index, if any
		// --------------------------------------------------------------------
		if  persist.table_name != "*" {
			table_index, ok = ds.schemas[schema_index].tables_map[persist.table_name]
			if ok != true {
				request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
				return
			}
		}
	}

	// -------------------------------------------------------------------------
	// Get the snapshot index
	// -------------------------------------------------------------------------
	snapshot_index, ok := ds.snapshot_descriptors_map[persist.snapshot_name]
	if ok != true {
		request.err = fmt.Errorf("Snapshot Name [%s] not found.", persist.snapshot_name)
		return
	}

	// -------------------------------------------------------------------------
	// Do the snapshot
	// -------------------------------------------------------------------------
	if persist.schema_name == "*" {
		for _, schema := range ds.schemas {		
			for _, table := range schema.tables {
				for col_index, _ := range table.columns {

					type_constant := table.descriptor.column_descriptors[col_index].type_constant
					if type_constant >= max_column_types {
						request.err = fmt.Errorf("Column type_constant [%d] exceeds max_column_type [%d].", 
									type_constant, max_column_types)
						return
					}

					err := takeSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index))
					if err != nil{
						request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
						return
					}
				}
			}
		}
	} else {

		table := ds.schemas[schema_index].tables[table_index]
		for col_index, _ := range table.columns {
			type_constant := table.descriptor.column_descriptors[col_index].type_constant
			err := takeSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index))
			if err != nil{
				request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
				return
			}
		}
	}
	request.err = nil
}


func restoreSnapshot (ds * datastore, request * request) {
	
	// -------------------------------------------------------------------------
	// Get a pointer to the request data
	// -------------------------------------------------------------------------
	persist, ok := request.data.(*persist)
	if ok == false {
		request.err = fmt.Errorf("internal error - 'persist' type assertion failed.") // random message
		return
	}

	// -------------------------------------------------------------------------
	// Get the schema_index, if any
	// -------------------------------------------------------------------------
	var schema_index uint64
	var table_index uint64
	if persist.schema_name != "*" {
		schema_index, ok = ds.schemas_map[persist.schema_name]
		if ok != true {
			request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
			return
		}

		// --------------------------------------------------------------------
		// Get the table_index, if any
		// --------------------------------------------------------------------
		if  persist.table_name != "*" {
			table_index, ok = ds.schemas[schema_index].tables_map[persist.table_name]
			if ok != true {
				request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
				return
			}
		}
	}

	// -------------------------------------------------------------------------
	// Get the snapshot index
	// -------------------------------------------------------------------------
	snapshot_index, ok := ds.snapshot_descriptors_map[persist.snapshot_name]
	if ok != true {
		request.err = fmt.Errorf("Snapshot Name [%s] not found.", persist.snapshot_name)
		return
	}

	// -------------------------------------------------------------------------
	// Do the restore snapshot
	// -------------------------------------------------------------------------
	if persist.schema_name == "*" {
		for _, schema := range ds.schemas {		
			for _, table := range schema.tables {
				for col_index, _ := range table.columns {

					type_constant := table.descriptor.column_descriptors[col_index].type_constant
					if type_constant >= max_column_types {
						request.err = fmt.Errorf("Column type_constant [%d] exceeds max_column_type [%d].", 
									type_constant, max_column_types)
						return
					}

					err := restoreSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index))
					if err != nil{
						request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
						return
					}
				}
			}
		}
	} else {

		table := ds.schemas[schema_index].tables[table_index]
		for col_index, _ := range table.columns {
			type_constant := table.descriptor.column_descriptors[col_index].type_constant
			err := restoreSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index))
			if err != nil{
				request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
				return
			}
		}
	}
	request.err = nil
}

// -----------------------------------------------------------------------------
// ClearSnapshot
// -----------------------------------------------------------------------------
func clearSnapshot (ds * datastore, request * request) {
	// -------------------------------------------------------------------------
	// Get a pointer to the request data
	// -------------------------------------------------------------------------
	persist, ok := request.data.(*persist)
	if ok == false {
		request.err = fmt.Errorf("internal error - 'persist' type assertion failed.") // random message
		return
	}

	// -------------------------------------------------------------------------
	// Get the schema_index, if any
	// -------------------------------------------------------------------------
	var schema_index uint64
	var table_index uint64
	if persist.schema_name != "*" {
		schema_index, ok = ds.schemas_map[persist.schema_name]
		if ok != true {
			request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
			return
		}

		// --------------------------------------------------------------------
		// Get the table_index, if any
		// --------------------------------------------------------------------
		if  persist.table_name != "*" {
			table_index, ok = ds.schemas[schema_index].tables_map[persist.table_name]
			if ok != true {
				request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
				return
			}
		}
	}

	// -------------------------------------------------------------------------
	// Get the snapshot index
	// -------------------------------------------------------------------------
	snapshot_index, ok := ds.snapshot_descriptors_map[persist.snapshot_name]
	if ok != true {
		request.err = fmt.Errorf("Snapshot Name [%s] not found.", persist.snapshot_name)
		return
	}

	// -------------------------------------------------------------------------
	// Do the restore snapshot
	// -------------------------------------------------------------------------
	if persist.schema_name == "*" {
		for _, schema := range ds.schemas {		
			for _, table := range schema.tables {
				for col_index, _ := range table.columns {

					type_constant := table.descriptor.column_descriptors[col_index].type_constant
					if type_constant >= max_column_types {
						request.err = fmt.Errorf("Column type_constant [%d] exceeds max_column_type [%d].", 
									type_constant, max_column_types)
						return
					}

					err := clearSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index))
					if err != nil{
						request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
						return
					}
				}
			}
		}
	} else {

		table := ds.schemas[schema_index].tables[table_index]
		for col_index, _ := range table.columns {
			type_constant := table.descriptor.column_descriptors[col_index].type_constant
			err := clearSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index))
			if err != nil{
				request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
				return
			}
		}
	}
	request.err = nil
}


// -----------------------------------------------------------------------------
// clearRealtime
// -----------------------------------------------------------------------------
func clearRealtime (ds * datastore, request * request) {

	// -------------------------------------------------------------------------
	// Get a pointer to the request data
	// -------------------------------------------------------------------------
	persist, ok := request.data.(*persist)
	if ok == false {
		request.err = fmt.Errorf("internal error - 'persist' type assertion failed.") // random message
		return
	}

	// -------------------------------------------------------------------------
	// Get the schema_index, if any
	// -------------------------------------------------------------------------
	var schema_index uint64
	var table_index uint64
	if persist.schema_name != "*" {
		schema_index, ok = ds.schemas_map[persist.schema_name]
		if ok != true {
			request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
			return
		}

		// --------------------------------------------------------------------
		// Get the table_index, if any
		// --------------------------------------------------------------------
		if  persist.table_name != "*" {
			table_index, ok = ds.schemas[schema_index].tables_map[persist.table_name]
			if ok != true {
				request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
				return
			}
		}
	}

	// -------------------------------------------------------------------------
	// Clear the realtime data
	// -------------------------------------------------------------------------
	if persist.schema_name == "*" {
		for _, schema := range ds.schemas {		
			for _, table := range schema.tables {
				for col_index, _ := range table.columns {

					type_constant := table.descriptor.column_descriptors[col_index].type_constant
					if type_constant >= max_column_types {
						request.err = fmt.Errorf("Column type_constant [%d] exceeds max_column_type [%d].", 
									type_constant, max_column_types)
						return
					}

					err := clearRealtimeFuncs[type_constant](table, uint16(col_index))
					if err != nil{
						request.err = fmt.Errorf("Error clear realtime func[%v].", err)
						return
					}
				}
			}
		}
	} else {

		table := ds.schemas[schema_index].tables[table_index]
		for col_index, _ := range table.columns {
			type_constant := table.descriptor.column_descriptors[col_index].type_constant
			err := clearRealtimeFuncs[type_constant](table, uint16(col_index))
			if err != nil{
				request.err = fmt.Errorf("Error clear realtime func[%v].", err)
				return
			}
		}
	}
	request.err = nil


}



// -----------------------------------------------------------------------------
// WriteSnapshot
// -----------------------------------------------------------------------------
func writeSnapshot (ds * datastore, request * request) {

	// -------------------------------------------------------------------------
	// Get a pointer to the request data
	// -------------------------------------------------------------------------
	persist, ok := request.data.(*persist)
	if ok == false {
		request.err = fmt.Errorf("internal error - 'persist' type assertion failed.") // random message
		return
	}

	// -------------------------------------------------------------------------
	// Open the file that we will write to.
	// -------------------------------------------------------------------------
	file, err := os.OpenFile(persist.path_file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		request.err =  fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close(); 

	// -------------------------------------------------------------------------
	// Get the schema_index, if any
	// -------------------------------------------------------------------------
	var schema_index uint64
	var table_index uint64
	if persist.schema_name != "*" {
		schema_index, ok = ds.schemas_map[persist.schema_name]
		if ok != true {
			request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
			return
		}

		// --------------------------------------------------------------------
		// Get the table_index, if any
		// --------------------------------------------------------------------
		if  persist.table_name != "*" {
			table_index, ok = ds.schemas[schema_index].tables_map[persist.table_name]
			if ok != true {
				request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
				return
			}
		}
	}

	// -------------------------------------------------------------------------
	// Get the snapshot index
	// -------------------------------------------------------------------------
	snapshot_index, ok := ds.snapshot_descriptors_map[persist.snapshot_name]
	if ok != true {
		request.err = fmt.Errorf("Snapshot Name [%s] not found.", persist.snapshot_name)
		return
	}

	// -------------------------------------------------------------------------
	// Do the snapshot write
	// TODO: Convert the following code (and in readSnapshot) to be a gorouting
	// 
	// -------------------------------------------------------------------------
	if persist.schema_name == "*" {
		for _, schema := range ds.schemas {		
			for _, table := range schema.tables {
				for col_index, _ := range table.columns {

					type_constant := table.descriptor.column_descriptors[col_index].type_constant
					if type_constant >= max_column_types {
						request.err = fmt.Errorf("Column type_constant [%d] exceeds max_column_type [%d].", 
									type_constant, max_column_types)
						return
					}

					err := writeSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index), file)
					if err != nil{
						request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
						return
					}
				}
			}
		}
	} else {

		table := ds.schemas[schema_index].tables[table_index]
		for col_index, _ := range table.columns {
			type_constant := table.descriptor.column_descriptors[col_index].type_constant
			err := writeSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index), file)
			if err != nil{
				request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
				return
			}
		}
	}
	
	request.err = nil
}


// -----------------------------------------------------------------------------
// readSnapshot
// -----------------------------------------------------------------------------


func readSnapshot (ds * datastore, request * request) {
	// -------------------------------------------------------------------------
	// Get a pointer to the request data
	// -------------------------------------------------------------------------
	persist, ok := request.data.(*persist)
	if ok == false {
		request.err = fmt.Errorf("internal error - 'persist' type assertion failed.") // random message
		return
	}

	// -------------------------------------------------------------------------
	// Open the file that we will readfrom
	// -------------------------------------------------------------------------
	file, err := os.OpenFile(persist.path_file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		request.err =  fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close(); 


	// -------------------------------------------------------------------------
	// Get the schema_index, if any
	// -------------------------------------------------------------------------
	var schema_index uint64
	var table_index uint64
	if persist.schema_name != "*" {
		schema_index, ok = ds.schemas_map[persist.schema_name]
		if ok != true {
			request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
			return
		}

		// --------------------------------------------------------------------
		// Get the table_index, if any
		// --------------------------------------------------------------------
		if  persist.table_name != "*" {
			table_index, ok = ds.schemas[schema_index].tables_map[persist.table_name]
			if ok != true {
				request.err = fmt.Errorf("Schema name [%s] not found.", persist.schema_name)
				return
			}
		}
	}

	// -------------------------------------------------------------------------
	// Get the snapshot index
	// -------------------------------------------------------------------------
	snapshot_index, ok := ds.snapshot_descriptors_map[persist.snapshot_name]
	if ok != true {
		request.err = fmt.Errorf("Snapshot Name [%s] not found.", persist.snapshot_name)
		return
	}

	// -------------------------------------------------------------------------
	// Do the snapshot write
	// -------------------------------------------------------------------------
	if persist.schema_name == "*" {
		for _, schema := range ds.schemas {		
			for _, table := range schema.tables {
				for col_index, _ := range table.columns {

					type_constant := table.descriptor.column_descriptors[col_index].type_constant
					if type_constant >= max_column_types {
						request.err = fmt.Errorf("Column type_constant [%d] exceeds max_column_type [%d].", 
									type_constant, max_column_types)
						return
					}

					err := readSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index), file)
					if err != nil{
						request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
						return
					}
				}
			}
		}
	} else {

		table := ds.schemas[schema_index].tables[table_index]
		for col_index, _ := range table.columns {
			type_constant := table.descriptor.column_descriptors[col_index].type_constant
			err := readSnapshotFuncs[type_constant](table, uint16(snapshot_index), uint16(col_index), file)
			if err != nil{
				request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
				return
			}
		}
	}
	request.err = nil
}


func writeContext (ds * datastore, request * request) {
	//	fmt.Printf("%s Request: Function: %d   DB: %s\n", ds.name, request.function_id, request.db_key)
		request.err = fmt.Errorf("Function_name [writeContext] under development .")
}

func readContext (ds * datastore, request * request) {
	//	fmt.Printf("%s Request: Function: %d   DB: %s\n", ds.name, request.function_id, request.db_key)
		request.err = fmt.Errorf("Function_name [readContext] under development .")
}



/******************************************************************************
takeSnapshot handlers 

TODO: move all the handlers to their own file when we add the rest of the types
******************************************************************************/

func takeSnapshotUndefined(tbl * table, snapshot_index uint16, col_index uint16) error {
	return fmt.Errorf("internal error - undefined column type [%d].", 
//				tbl.columns[col_index].descriptor.type_constant) 
				tbl.descriptor.column_descriptors[col_index].type_constant) 
}

// ----------------------------------------------------------------------------
// takeSnapshotUTF8() will will for both BLOBS and VARS. 
// ----------------------------------------------------------------------------
func takeSnapshotUTF8(tbl * table, snapshot_index uint16, col_index uint16) error {
	source, ok 	:= tbl.columns[col_index].realtime.([]string)
	if ok == false {
		return fmt.Errorf("internal error - source type assertion failed.") 
	} 					

	dest, ok := tbl.columns[col_index].snapshots[snapshot_index].([]string)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = source[rowCntr]
	}
	
	//fmt.Printf("Snapshot take string: \n")
	return nil
}


func takeSnapshotInt64(tbl * table, snapshot_index uint16, col_index uint16) error {
	
	source, ok 	:= tbl.columns[col_index].realtime.([]int64)
	if ok == false {
		return fmt.Errorf("internal error - source type assertion failed.") 
	} 					

	dest, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]int64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = source[rowCntr]
	}
	
	//fmt.Printf("Snapshot take int64: \n")
	return nil
}

func takeSnapshotFloat64(tbl * table, snapshot_index uint16, col_index uint16) error {
	
	source, ok 	:= tbl.columns[col_index].realtime.([]float64)
	if ok == false {
		return fmt.Errorf("internal error - source type assertion failed.") 
	} 					

	dest, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]float64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = source[rowCntr]
	}
	
	//fmt.Printf("Snapshot take float64: \n")
	return nil
}


/******************************************************************************
restoreSnapshot handlers 

TODO: move all the handlers to their own file when we add the rest of the types
******************************************************************************/

func restoreSnapshotUndefined(tbl * table, snapshot_index uint16, col_index uint16) error {
	return fmt.Errorf("internal error - undefined column type [%d].", 
			tbl.descriptor.column_descriptors[col_index].type_constant) 
}


// ----------------------------------------------------------------------------
// takeSnapshotUTF8() will will for both BLOBS and VARS. 
// ----------------------------------------------------------------------------
func restoreSnapshotUTF8(tbl * table, snapshot_index uint16, col_index uint16) error {

	source, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]string)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	dest, ok 	:= tbl.columns[col_index].realtime.([]string)
	if ok == false {
		return fmt.Errorf("internal error - source type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = source[rowCntr]
	}
	
	//fmt.Printf("Snapshot restore string: \n")
	return nil
}

func restoreSnapshotInt64(tbl * table, snapshot_index uint16, col_index uint16) error {
	
	source, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]int64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	dest, ok 	:= tbl.columns[col_index].realtime.([]int64)
	if ok == false {
		return fmt.Errorf("internal error - source type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = source[rowCntr]
	}
	
	//fmt.Printf("Snapshot restore int64: \n")
	return nil
}

func restoreSnapshotFloat64(tbl * table, snapshot_index uint16, col_index uint16) error {
	
	source, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]float64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	dest, ok 	:= tbl.columns[col_index].realtime.([]float64)
	if ok == false {
		return fmt.Errorf("internal error - source type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = source[rowCntr]
	}
	
	//fmt.Printf("Snapshot restore float64: \n")
	return nil
}
	


/******************************************************************************
writeSnapshot handlers 

TODO: move all the handlers to their own file when we add the rest of the types
******************************************************************************/

func writeSnapshotUndefined(tbl * table, snapshot_index  uint16, col_index uint16, file * os.File) error {
	return fmt.Errorf("internal error - undefined column type [%d].", 
					tbl.descriptor.column_descriptors[col_index].type_constant) 
}


func writeSnapshotVarUTF8(tbl * table, snapshot_index  uint16, col_index uint16, file * os.File) error {

	source, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]string)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					


	// Copy the data into a []byte buffer
	buffer_length := tbl.descriptor.column_descriptors[col_index].length
	var buffer bytes.Buffer
	paddedByteArray := make([]byte, buffer_length)

	for _, datum := range source {

		paddedByteArray = bytes.Repeat([]byte{0}, int(buffer_length))

		byteArray := []byte(datum)
		copy(paddedByteArray, byteArray)

		_, err := buffer.Write( paddedByteArray )
		if err != nil {
			return fmt.Errorf("internal error - writing array of data to buffer [%v]", err) 
		}
	}

	// write the data to disk
	_, err := file.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("internal error - file write failed.") 
	}
	return nil
}

func writeSnapshotBlobUTF8(tbl * table, snapshot_index  uint16, col_index uint16, file * os.File) error {

	source, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]string)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	for _, datum := range source {
		// write the length to disk
		var buffer bytes.Buffer
		binary.Write(file, binary.LittleEndian, len(datum))
		_, err := file.Write(buffer.Bytes())
		if err != nil {
			return fmt.Errorf("file write failed.") 
		}
	

		// write the data to disk
		_, err = file.WriteString(datum)
		if err != nil {
		return fmt.Errorf("file write failed.") 
		}
	}
	return nil
}



func writeSnapshotInt64(tbl * table, snapshot_index uint16, col_index uint16, file * os.File) error {

	source, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]int64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	// Copy the data into a []byte buffer
	var buffer bytes.Buffer

	for _, datum := range source {
		binary.Write(&buffer, binary.LittleEndian, datum)
		_, err := file.Write(buffer.Bytes())
		if err != nil {
			return fmt.Errorf("writing array of data to buffer [%v]", err) 
		}
	}

	// write the data to disk
	_, err := file.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("file write failed.") 
	}

	//fmt.Printf("Int64 bytes written %d\n", bytes_written)

	return nil
}


func writeSnapshotFloat64(tbl * table, snapshot_index uint16, col_index uint16, file * os.File) error {

	source, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]float64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	// Copy the data into a []byte buffer
	var buffer bytes.Buffer

	for _, datum := range source {
		err := binary.Write(&buffer, binary.LittleEndian, datum)
		if err != nil {
			return fmt.Errorf("internal error - writing array of data to buffer.") 
		}
	}

	// write the data to disk
//	bytes_written, err := file.Write(buffer.Bytes())
	_, err := file.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("internal error - file write failed.") 
	}

	//fmt.Printf("Float64 bytes written %d\n", bytes_written)
	return nil
}
	


/******************************************************************************
readSnapshot handlers 

TODO: move all the handlers to their own file when we add the rest of the types
******************************************************************************/


func readSnapshotUndefined(tbl * table, snapshot_index uint16, col_index uint16, file * os.File) error {
	return fmt.Errorf("internal error - undefined column type [%d].", 
	tbl.descriptor.column_descriptors[col_index].type_constant) 
}

// -----------------------------------------------------------------------------
// readSnapshotVarUTF8()
// -----------------------------------------------------------------------------
func readSnapshotVarUTF8(tbl * table, snapshot_index uint16, col_index uint16, file * os.File) error {

	dest, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]string)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	field_size := tbl.descriptor.column_descriptors[col_index].length

	// Make a []byte buffer to receive the data
	file_buffer := make([]byte, tbl.descriptor.total_rows * field_size)

	// read the data from disk
	_, err := file.Read(file_buffer)
//	bytes_read, err := file.Read(file_buffer)
	if err != nil {
		return fmt.Errorf("file read failed.") 
	}

	// copy the data to for each varchar, convert it to a string, and
	// put that in the destination
	
	for row_cntr := uint64(0); row_cntr < tbl.descriptor.total_rows; row_cntr++ {
		start := row_cntr * uint64(field_size)
		end := start + field_size
		str_buf := file_buffer[start : end]
		dest[row_cntr] = string(str_buf)
	}


	//fmt.Printf("Writing varchar\n")
	return nil
}


// -----------------------------------------------------------------------------
// readSnapshotBlobUTF8()
// -----------------------------------------------------------------------------
 unc readSnapshotBlobUTF8(tbl * table, snapshot_index  uint16, col_index uint16, file * os.File) error {

	dest, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]string)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	for row_cntr, _ := range dest {
	
		//----------------------------------------------------------------------
		// wread the length from disk
		//----------------------------------------------------------------------
		var length int
		err := binary.Read(file, binary.LittleEndian, length)
		if err != nil {
			return fmt.Errorf("file read failed.") 
		}
	
		//----------------------------------------------------------------------
		// read the data from disk
		//----------------------------------------------------------------------
		buffer := make([]byte, length)

		n, err := file.Read(buffer)
		if err != nil {
			return fmt.Errorf("file read failed.") 
		}


		if n!= length {
			return fmt.Errorf("file read failed.") 
		}
		
		dest[row_cntr] = string(buffer)
		//datum = datum
	}

	return nil
}



// -----------------------------------------------------------------------------
// readSnapshotInt64()
// -----------------------------------------------------------------------------
func readSnapshotInt64(tbl * table, snapshot_index uint16, col_index uint16, file * os.File) error {

	dest, ok := tbl.columns[col_index].snapshots[snapshot_index].([]int64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	// Make a []byte buffer to receive the data
	buffer := make([]byte, tbl.descriptor.total_rows * 8) // 8 is the size of an int64

	// read the data from disk
//	bytes_read, err := file.Read(buffer)
	_, err := file.Read(buffer)
	if err != nil {
		return fmt.Errorf("internal error - file read failed.") 
	}


	// move the data to the destination
	buf := bytes.NewReader(buffer)
	for _, datum := range dest {
		err := binary.Read(buf, binary.LittleEndian, &datum)
	
		if err != nil {
			return fmt.Errorf("internal error - reading array of data to buffer [%v]", err) 
		}
	}
	//fmt.Printf("Int64 bytes read %d\n", bytes_read)
	return nil
}


// -----------------------------------------------------------------------------
// readSnapshotFloat64()
// -----------------------------------------------------------------------------
func readSnapshotFloat64(tbl * table, snapshot_index uint16, col_index uint16, file * os.File) error {

	dest, ok := tbl.columns[col_index].snapshots[snapshot_index].([]float64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	// Make a []byte buffer to receive the data
	buffer := make([]byte, tbl.descriptor.total_rows * 8) // 8 is the size of an float64

	// read the data from disk
//	bytes_read, err := file.Read(buffer)
	_, err := file.Read(buffer)
	if err != nil {
		return fmt.Errorf("internal error - file read failed.") 
	}

// func (f *File) Read(b []byte) (n int, err error)
	buf := bytes.NewReader(buffer)

	for _, datum := range dest {
		err := binary.Read(buf, binary.LittleEndian, &datum)
	
		if err != nil {
			return fmt.Errorf("internal error - reading array of data to buffer [%v]", err) 
		}
	}
	//fmt.Printf("Int64 bytes read %d\n", bytes_read)
	return nil
}




/******************************************************************************
clearSnapshot handlers 

TODO: move all the handlers to their own file when we add the rest of the types
******************************************************************************/


func clearSnapshotUndefined(tbl * table, snapshot_index uint16, col_index uint16) error {
	return fmt.Errorf("internal error - undefined column type [%d].", 
				tbl.descriptor.column_descriptors[col_index].type_constant) 
}

func clearSnapshotUTF8(tbl * table, snapshot_index uint16, col_index uint16) error {

	dest, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]string)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = ""
	}
	
	//fmt.Printf("Snapshot take string: \n")
	return nil
}



func clearSnapshotInt64(tbl * table, snapshot_index uint16, col_index uint16) error {
	

	dest, ok 	:= tbl.columns[col_index].snapshots[snapshot_index].([]int64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = 0
	}
	
	//fmt.Printf("Snapshot take int64: \n")
	return nil
}

func clearSnapshotFloat64(tbl * table, snapshot_index uint16, col_index uint16) error {
	

	dest, ok := tbl.columns[col_index].snapshots[snapshot_index].([]float64)
	if ok == false {
		return fmt.Errorf("internal error - destination type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = 0.0
	}
	
	//fmt.Printf("Snapshot take float64: \n")
	return nil
}

/******************************************************************************
clearRealtime handlers 

TODO: move all the handlers to their own file when we add the rest of the types
******************************************************************************/

func clearRealtimeUndefined(tbl * table, col_index uint16) error {
	return fmt.Errorf("internal error - undefined column type [%d].", 
				tbl.descriptor.column_descriptors[col_index].type_constant) 
}

func clearRealtimeUTF8(tbl * table, col_index uint16) error {


	dest, ok 	:= tbl.columns[col_index].realtime.([]string)
	if ok == false {
		return fmt.Errorf("internal error - source type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = ""
	}
	
	//fmt.Printf("Snapshot restore string: \n")
	return nil
}

func clearRealtimeInt64(tbl * table, col_index uint16) error {
	

	dest, ok 	:= tbl.columns[col_index].realtime.([]int64)
	if ok == false {
		return fmt.Errorf("internal error - source type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = 0
	}
	
	//fmt.Printf("Snapshot restore int64: \n")
	return nil
}

func clearRealtimeFloat64(tbl * table, col_index uint16) error {
	

	dest, ok 	:= tbl.columns[col_index].realtime.([]float64)
	if ok == false {
		return fmt.Errorf("internal error - source type assertion failed.") 
	} 					

	for rowCntr := uint64(0); rowCntr < tbl.descriptor.total_rows; rowCntr++ {
		dest[rowCntr] = 0.0
	}
	
	//fmt.Printf("Snapshot restore float64: \n")
	return nil
}
	
