/*******************************************************************************
rtds_dn.go

Purpose: data manipulation server functions

Todo: complete writing a set of validation functions

*******************************************************************************/

package rtdsms

import (
	"fmt"
	"math"
	"time"
	"regexp"
//	"reflect"
)




// TODO: figure out where to put a single copy of this. It is kind of
// scattered all around
func validateObjectName(name string) error {

	// Validate the non-type dependent parameters in the descriptor
	if len(name) > max_object_name_length {
		return	fmt.Errorf("Object name is too long: %s", name)
	}

	// Check if the name contains any disallowed characters
	regex := regexp.MustCompile(object_name_format)
	if ! regex.MatchString(name) {
		return 	fmt.Errorf("Object Name contains invalid characters: %s", name)
	}	 
	
	return nil
}





type dmFunc func(tbl *table, row_index uint64, variant *variant) error
type deleteFunc func(tbl *table, row_index uint64, columnIndex uint16) error

// validateFuncs are used by the insert and update functions to check values
// that are to be written to the database.
type validateFunc func(col_desc * ColumnDescriptor, variant *variant) error


var selectFuncs		[max_column_types]dmFunc
var insertFuncs		[max_column_types]dmFunc
var	updateFuncs 	[max_column_types]dmFunc
var	deleteFuncs 	[max_column_types]deleteFunc
var	validateFuncs 	[max_column_types]validateFunc


func init() {
	// Initialize the selectFuncs array
	for i := 0; i < max_column_types; i++{
		selectFuncs[i] = selectUndefined
	}
	selectFuncs[COL_TYPE_VAR_UTF8] 	= selectVarUTF8
	selectFuncs[COL_TYPE_I64] 		= selectInt64
	selectFuncs[COL_TYPE_F64] 		= selectFloat64

	// Initialize the insertFuncs array
	for i := 0; i < max_column_types; i++{
		insertFuncs[i] = insertUndefined
	}
	insertFuncs[COL_TYPE_VAR_UTF8] 	= insertVarUTF8
	insertFuncs[COL_TYPE_I64] 		= insertInt64
	insertFuncs[COL_TYPE_F64] 		= insertFloat64


	// Initialize the updateFuncs array
	for i := 0; i < max_column_types; i++{
		updateFuncs[i] = updateUndefined
	}
	updateFuncs[COL_TYPE_VAR_UTF8] 	= updateVarUTF8
	updateFuncs[COL_TYPE_I64] 		= updateInt64
	updateFuncs[COL_TYPE_F64] 		= updateFloat64

	// Initialize the deleteFuncs array
	for i := 0; i < max_column_types; i++{
		deleteFuncs[i] = deleteUndefined
	}
	deleteFuncs[COL_TYPE_VAR_UTF8] 	= deleteVarUTF8
	deleteFuncs[COL_TYPE_I64] 		= deleteInt64
	deleteFuncs[COL_TYPE_F64] 		= deleteFloat64

	// Initialize the validateFuncs array
	for i := 0; i < max_column_types; i++{
		validateFuncs[i] = validateUndefined
	}
	validateFuncs[COL_TYPE_VAR_UTF8] 	= validateVarUTF8
	validateFuncs[COL_TYPE_I64] 		= validateInt64
	validateFuncs[COL_TYPE_F64] 		= validateFloat64

}






// --------------------------------------------------------------------------------------------------------------------
// select_() attempts to select one or more records in the database
// --------------------------------------------------------------------------------------------------------------------
func select_ (ds * Datastore, request * Request) {

	// Get a pointer to the records
	records, ok := request.data.( []*Record)
	if ok == false {
		request.err = fmt.Errorf("internal error - record type assertion failed.") 
		return
	}
	
	for _, record := range  records {

		if record.indexesValid == false	{
			// ---------------------------------------------------------------------
			// Get the index to the schema, which also validates the schema
			// ---------------------------------------------------------------------
			schema_index, ok := ds.schemasMap[record.schemaKey]
			if ok == false {
				request.err =  fmt.Errorf("Schema [%s] not found in datastore [%s]", record.schemaKey, ds.descriptor.datastoreName ) 
				return
			}
			record.schemaIndex = uint64(schema_index)
		
			// ---------------------------------------------------------------------
			// Get the index to the table, which also validates the table
			// ---------------------------------------------------------------------
			table_index, ok := ds.schemas[schema_index].tablesMap[record.tableKey]
			if ok == false {
				request.err =  fmt.Errorf("Table [%s] not found in datastore [%s] schema [%s]", 
								record.tableKey, ds.descriptor.datastoreName, record.schemaKey ) 
				return
			}
			record.tableIndex = uint64(table_index)

			// ---------------------------------------------------------------------
			// use these variable to make later code easier to read - and maybe a 
			// little faster
			// ---------------------------------------------------------------------
			table_desc := ds.schemas[record.schemaIndex].tables[record.tableIndex].descriptor 
			table := ds.schemas[schema_index].tables[table_index]
	
			// ---------------------------------------------------------------------
			// Get an index for the column we use to look up the row
			// ---------------------------------------------------------------------
			lookup_columnIndex, ok := table_desc.columnsMap[record.lookupColumnKey]
			if ok == false {
				request.err = fmt.Errorf("lookup column name [%s] does not exist in table [%s]", 
										record.lookupColumnKey, record.tableKey)
				return
			}
			record.lookupColumnIndex = lookup_columnIndex

			// ---------------------------------------------------------------------
			// Get an index for the row_lookupKey
			// ---------------------------------------------------------------------

//		fmt.Printf("Lookup %d %s\n", lookup_columnIndex, record.lookup_rowKey)

			row_index, ok := table.columns[lookup_columnIndex].index[record.lookupRowKey]
			if ok == false {
				request.err = fmt.Errorf("row name [%s] does not exist in table [%s]", 
										record.lookupRowKey, record.tableKey)
				return
			}
			record.lookupRowIndex = uint64(row_index)

			// ---------------------------------------------------------------------
			// check that all the variant columns exist
			// ---------------------------------------------------------------------
			for _, variant := range record.variants{	
				columnIndex, ok := table_desc.columnsMap[variant.columnKey]
				if ok == false {
					request.err = fmt.Errorf("column key [%s] not found in table [%s]", 
										variant.columnKey, table.descriptor.tableName) 
					return
				}
				variant.columnIndex = uint16(columnIndex)		
			}
			record.indexesValid = true
		}
	}

	// -------------------------------------------------------------------------
	// everything validated, so copy the data and finish		
	// -------------------------------------------------------------------------
	for _, record := range records {

		table := ds.schemas[record.schemaIndex].tables[record.tableIndex]
		column_descriptors := table.descriptor.ColumnDescriptors
		row_index := record.lookupRowIndex

		// copy each of the variants from the database
		for _, variant := range record.variants{	
		
			typeConst := column_descriptors[variant.columnIndex].typeConstant

			if typeConst > max_column_types {
				request.err = fmt.Errorf("column type [%d] > max_column_types [%d].", typeConst, max_column_types)
			}

			err := selectFuncs[typeConst](table, row_index, variant)
			if err != nil{
				request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
				return
			}		
		}
		table.lastSelect = time.Now()
		table.selectCounter++
	}

	request.err = nil
}




// --------------------------------------------------------------------------------------------------------------------
// insert() attempts to insert one or more records in the database
// --------------------------------------------------------------------------------------------------------------------
func insert (ds * Datastore, request * Request) {

	// Get a pointer to the records
	records, ok := request.data.([]*Record)
	if ok == false {
		request.err = fmt.Errorf("internal error - record type assertion failed.") 
		return
	}


	// Validate each of the records
	for cntr1, record := range records {


		if record.indexesValid == false {
			// ---------------------------------------------------------------------
			// Get the index to the schema, which also validates the schema
			// ---------------------------------------------------------------------
			schema_index, ok := ds.schemasMap[record.schemaKey]
			if ok == false {
				request.err =  fmt.Errorf("Schema [%s] not found in datastore [%s]", record.schemaKey, ds.descriptor.datastoreName ) 
				return
			}
			record.schemaIndex = uint64(schema_index)

			// ---------------------------------------------------------------------
			// Get the index to the table, which also validates the table
			// ---------------------------------------------------------------------
			table_index, ok := ds.schemas[schema_index].tablesMap[record.tableKey]
			if ok == false {
				request.err =  fmt.Errorf("Table [%s] not found in datastore [%s] schema [%s]", 
								record.tableKey, ds.descriptor.datastoreName, record.schemaKey ) 
				return
			}
			record.tableIndex = uint64(table_index)
		}


		// ---------------------------------------------------------------------
		// Check that the 1st variant column key is 'row_id'
		// ---------------------------------------------------------------------
		if record.variants[0].columnKey != "row_id" {
			request.err = fmt.Errorf("first variant column key myst be [row_id] for insert(). [%s] was supplied.", 
									record.variants[0].columnKey) 
			return
		}
		record.variants[0].columnIndex=0 // probably not a necessary statement

		// ---------------------------------------------------------------------
		// Assert that the row name aupp[ied by the user is a string
		// ---------------------------------------------------------------------
		row_id_value, ok := record.variants[0].value.(string)
		if ok == false{
			request.err = fmt.Errorf("Column [%s] value in variant [%v] must be string.", 
							record.variants[0].columnKey, record.variants[0].value)
			return
		} 

		// ---------------------------------------------------------------------
		// Check that the row id value  aupp[ied by the caller is a valid object name
		// ---------------------------------------------------------------------
		err := validateObjectName(row_id_value) 
		if err != nil{
			request.err = fmt.Errorf("row_id value [%s] is not a valid object name: %s", record.variants[0].value, err) 
			return
		}

		// ---------------------------------------------------------------------
		// check if the table for this record is full
		// ---------------------------------------------------------------------
		if ds.schemas[record.schemaIndex].tables[record.tableIndex].rowsEmpty == 0{
			request.err = fmt.Errorf("table [%s] is full", record.tableKey) 
			return
		}

		// ---------------------------------------------------------------------
		// Check the table key and the new row_id in this record against all
		// the remaining records to look for duplications to prevent 
		// inserting two new rows in the same table with the same key.
		// ---------------------------------------------------------------------
		for cntr2 := cntr1+1; cntr2 < len(records); cntr2++ {
			if  (record.tableKey == (records)[cntr2].tableKey) &&
				(record.variants[0].value ==  (records)[cntr2].variants[0].value){
					request.err = fmt.Errorf("duplicate record row_id values found: table [%s] & row [%s]", 
										record.tableKey, record.variants[0].value ) 
					return
			}			
		}
			
		// ---------------------------------------------------------------------
		// use this variable to make later code easier to read 
		// ---------------------------------------------------------------------
		tbl_desc := ds.schemas[record.schemaIndex].tables[record.tableIndex].descriptor 
		table := ds.schemas[record.schemaIndex].tables[record.tableIndex]

		// ---------------------------------------------------------------------
		// Get an index for the column in each of the variants
		// ---------------------------------------------------------------------
		for _, variant := range record.variants{	
	
			columnIndex, ok := tbl_desc.columnsMap[variant.columnKey]
			if ok == false {
				request.err = fmt.Errorf("Column key [%s] does not exist in table [%s].", 
										variant.columnKey, record.tableKey)
				return
			}
			variant.columnIndex = columnIndex

			// TODO: Check the length of the supplied value against the column length



			// -----------------------------------------------------------------
			// if the column is indexed, check that the column type is a string and that
			// and that the variant value is a string and that the supplied value is not 
			// already used.
			// -----------------------------------------------------------------
			if tbl_desc.ColumnDescriptors[columnIndex].indexed == true {
				if tbl_desc.ColumnDescriptors[columnIndex].typeConstant != COL_TYPE_VAR_UTF8 {
					request.err = fmt.Errorf("Column [%s] must be string type to be indexed.", 
										variant.columnKey)
					return
				}

				value, ok := variant.value.(string)
				if ok == false{
					request.err = fmt.Errorf("Column [%s] is indexed and value in variant [%v] must be string.", 
										variant.columnKey, variant.value)
					return
				} 
				
				_, ok = table.columns[columnIndex].index[value] 
				if ok == true{
					request.err = fmt.Errorf("Column [%s] is indexed and value [%s] already esists.", 
										variant.columnKey, variant.value)
					return
				} 
			}
		}

		// ---------------------------------------------------------------------
		// do a type assertion on the row_id value to make sure it is a string.
		// ---------------------------------------------------------------------
		row_id, ok := record.variants[0].value.(string)
		if ok == false {
			request.err = fmt.Errorf("internal error - variant type assertion failed on row_id value.") 
			return
		} 		
		
		// TODO: Do a validate object name?

		// ---------------------------------------------------------------------
		// Check if the row_id is already used. This was done above so maybe dont need it?
		// ---------------------------------------------------------------------

//		fmt.Printf("Inserting ROW_ID [%s]\n", row_id)

		_, ok = table.columns[0].index[row_id]
		if ok == true {
			request.err = fmt.Errorf("row id [%s] already exists in table [%s]", 
									record.variants[0].value, record.tableKey)
			return
		}
		
		record.indexesValid = true
	}
	

	// ------------------------------------------------------------------------
	// Finally we are ready to insert the new rows into the database. We set
	// the table index in the record above.
	// ------------------------------------------------------------------------
	
	for _, record := range records {
		schema_index, _ := ds.schemasMap[record.schemaKey]
		table := ds.schemas[schema_index].tables[record.tableIndex]
		//total_rows := table.descriptor.total_rows
		
		// TODO: As the tables get larger, we need to find a way to make this 
		// faster. Possibly with a linked list or a queue of empty rows or 
		// something.
		// 'row_index' will hold the index of an empty row in the table when the lopp breaks

		/*
		var row_index uint64
		for  row_index = uint64(0); row_index < total_rows; row_index++  {
			if table.row_empty[row_index] == true {
				break
			}
		}*/

		// Here is the new code to track rows. MUCH faster as the table grows
		row_index := table.firstEmptyRow
		table.firstEmptyRow = table.nextEmptyRow[row_index]
		table.nextEmptyRow[row_index] = math.MaxUint64

		// copy each of the variants to the database
		for _, variant := range record.variants{	
	
			typeConst := table.descriptor.ColumnDescriptors[variant.columnIndex].typeConstant
						
			err := insertFuncs[typeConst](table, row_index, variant)
			if err != nil{
				request.err = fmt.Errorf("Error taking snapshot func[%v].", err)
				return
			}	
		}

		table.rowEmpty[row_index] = false
		table.rowsEmpty = table.rowsEmpty - 1	
		table.lastInsert = time.Now()
		table.insertCounter++

		//	tbl.rowKeys[i] = record.rowKey
	//	tbl.rowMap[record.rowKey] = i
	}

	request.err = nil
	
}



// --------------------------------------------------------------------------------------------------------------------
// update() attempts to update one or more records in the database
// --------------------------------------------------------------------------------------------------------------------
func update (ds * Datastore, request * Request) {

	
	// Get a pointer to the records
	records, ok := request.data.( []*Record)
	if ok == false {
		request.err = fmt.Errorf("internal error - record type assertion failed.") 
		return
	}


	for _, record := range  records {
		if record.indexesValid == false	{

			// ---------------------------------------------------------------------
			// Get the index to the schema, which also validates the schema
			// ---------------------------------------------------------------------
			schema_index, ok := ds.schemasMap[record.schemaKey]
			if ok == false {
				request.err =  fmt.Errorf("Schema [%s] not found in datastore [%s]", record.schemaKey, ds.descriptor.datastoreName ) 
				return
			}
			record.schemaIndex = uint64(schema_index)
		
			// ---------------------------------------------------------------------
			// Get the index to the table, which also validates the table
			// ---------------------------------------------------------------------
			table_index, ok := ds.schemas[schema_index].tablesMap[record.tableKey]
			if ok == false {
				request.err =  fmt.Errorf("Table [%s] not found in datastore [%s] schema [%s]", 
								record.tableKey, ds.descriptor.datastoreName, record.schemaKey ) 
				return
			}
			record.tableIndex = uint64(table_index)


			// ---------------------------------------------------------------------
			// use these variable to make later code easier to read - and maybe a 
			// little faster
			// ---------------------------------------------------------------------
			table_desc := ds.schemas[record.schemaIndex].tables[record.tableIndex].descriptor 
			table := ds.schemas[schema_index].tables[table_index]

			// ---------------------------------------------------------------------
			// Get an index for the column we use to look up the row
			// ---------------------------------------------------------------------
			lookup_column_index, ok := table_desc.columnsMap[record.lookupColumnKey]
			if ok == false {
				request.err = fmt.Errorf("lookup column name [%s] does not exist in table [%s]", 
										record.lookupColumnKey, record.tableKey)
				return
			}
			record.lookupColumnIndex = lookup_column_index

			// ---------------------------------------------------------------------
			// Get an index for the row_lookupKey
			// ---------------------------------------------------------------------

//		fmt.Printf("Lookup %d %s\n", lookup_columnIndex, record.lookup_rowKey)

			row_index, ok := table.columns[lookup_column_index].index[record.lookupRowKey]
			if ok == false {
				request.err = fmt.Errorf("row name [%s] does not exist in table [%s]", 
										record.lookupRowKey, record.tableKey)
				return
			}
			record.lookupRowIndex = uint64(row_index)

			// ---------------------------------------------------------------------
			// check that all the variant columns exist
			// ---------------------------------------------------------------------
			for _, variant := range record.variants{	
						columnIndex, ok := table_desc.columnsMap[variant.columnKey]
				if ok == false {
					request.err = fmt.Errorf("column key [%s] not found in table [%s]", 
									variant.columnKey, table.descriptor.tableName) 
					return
				}
				variant.columnIndex = uint16(columnIndex)		
			}
			record.indexesValid = true
		}
	}

	// ------------------------------------------------------------------------
	// Finally we are ready to update the rows into the database. We set
	// the table index in the record above.
	// ------------------------------------------------------------------------
	
	for _, record := range records {
		table := ds.schemas[record.schemaIndex].tables[record.tableIndex]
		
		// copy each of the variants to the database
		for _, variant := range record.variants{	
	
			typeConst := table.descriptor.ColumnDescriptors[variant.columnIndex].typeConstant
			
			err := updateFuncs[typeConst](table, record.lookupRowIndex, variant)
			if err != nil{
				request.err = fmt.Errorf("Error updating table[%v].", err)
				return
			}	
		}
		table.lastUpdate = time.Now()
		table.updateCounter++

	}

	request.err = nil
	
}

// --------------------------------------------------------------------------------------------------------------------
// delete() attempts to delete one or more records from the database
// --------------------------------------------------------------------------------------------------------------------
func delete (ds * Datastore, request * Request) {
	
	// -------------------------------------------------------------------------
	// Get a pointer to the records
	// -------------------------------------------------------------------------
	records, ok := request.data.( []*Record)
	if ok == false {
		request.err = fmt.Errorf("internal error - record type assertion failed.") 
		return
	}

	// -------------------------------------------------------------------------
	// All we are going to use the records for is to lookup the row to delete
	// -------------------------------------------------------------------------
	for _, record := range  records {
		if record.indexesValid == false	{

			// ---------------------------------------------------------------------
			// Get the index to the schema, which also validates the schema
			// ---------------------------------------------------------------------
			schema_index, ok := ds.schemasMap[record.schemaKey]
			if ok == false {
				request.err =  fmt.Errorf("Schema [%s] not found in datastore [%s]", record.schemaKey, ds.descriptor.datastoreName ) 
				return
			}
			record.schemaIndex = uint64(schema_index)
		
			// ---------------------------------------------------------------------
			// Get the index to the table, which also validates the table
			// ---------------------------------------------------------------------
			table_index, ok := ds.schemas[schema_index].tablesMap[record.tableKey]
			if ok == false {
				request.err =  fmt.Errorf("Table [%s] not found in datastore [%s] schema [%s]", 
								record.tableKey, ds.descriptor.datastoreName, record.schemaKey ) 
				return
			}
			record.tableIndex = uint64(table_index)


			// ---------------------------------------------------------------------
			// use these variable to make later code easier to read - and maybe a 
			// little faster
			// ---------------------------------------------------------------------
			table_desc := ds.schemas[record.schemaIndex].tables[record.tableIndex].descriptor 
			table := ds.schemas[schema_index].tables[table_index]

			// ---------------------------------------------------------------------
			// Get an index for the column we use to look up the row
			// ---------------------------------------------------------------------
			lookup_column_index, ok := table_desc.columnsMap[record.lookupColumnKey]
			if ok == false {
				request.err = fmt.Errorf("lookup column name [%s] does not exist in table [%s]", 
										record.lookupColumnKey, record.tableKey)
				return
			}
			record.lookupColumnIndex = lookup_column_index

			// ---------------------------------------------------------------------
			// Get an index for the row_lookupKey
			// ---------------------------------------------------------------------

			row_index, ok := table.columns[lookup_column_index].index[record.lookupRowKey]
			if ok == false {
				request.err = fmt.Errorf("row name [%s] does not exist in table [%s]", 
										record.lookupRowKey, record.tableKey)
				return
			}
			record.lookupRowIndex = uint64(row_index)
		}
		record.indexesValid = false //because we are deleting this record
	}

	
	// -------------------------------------------------------------------------
	// set the value of all the fields in all the rows to their zero value
	// -------------------------------------------------------------------------
	for _, record := range records {
		table := ds.schemas[record.schemaIndex].tables[record.tableIndex]		
		for col_index, _ := range(table.descriptor.ColumnDescriptors) {
	
			typeConst := table.descriptor.ColumnDescriptors[col_index].typeConstant
			
			err := deleteFuncs[typeConst](table, record.lookupRowIndex,uint16(col_index))
			if err != nil{
				request.err = fmt.Errorf("Error deleting [%v].", err)
				return
			}	
		}
		table.rowsEmpty++
		

		table.lastDelete = time.Now()
		table.deleteCounter++
	}

	request.err = nil
}



/******************************************************************************
select functions
******************************************************************************/
func selectUndefined (tbl *table, row_index uint64, variant *variant) error {

	err := fmt.Errorf("internal error - type const not supported.") 
	return err
}

func selectVarUTF8 (tbl *table, row_index uint64, variant *variant) error {

	data, ok := tbl.columns[variant.columnIndex].realtime.([]string)
	if ok == false {
		err := fmt.Errorf("internal error - variant type assertion failed.") 
		return err
	} 		
	variant.value = (data)[row_index] 
	return nil
}

func selectInt64 (tbl *table, row_index uint64, variant *variant) error {

	data, ok := tbl.columns[variant.columnIndex].realtime.([]int64)
	if ok == false {
		err := fmt.Errorf("internal error - variant type assertion failed.") 
		return err
	} 		
	variant.value = (data)[row_index] 
	return nil
}

func selectFloat64 (tbl *table, row_index uint64, variant *variant) error {

	data, ok := tbl.columns[variant.columnIndex].realtime.([]float64)
	if ok == false {
		err := fmt.Errorf("internal error - variant type assertion failed.") 
		return err
	} 		
	variant.value = (data)[row_index] 
	return nil

}



/******************************************************************************
insert functions
******************************************************************************/
func insertUndefined (table *table, row_index uint64, variant *variant) error {

	err := fmt.Errorf("internal error - type const not supported.") 
	return err
}

func insertVarUTF8 (table *table, row_index uint64, variant *variant) error {

	data, ok := table.columns[variant.columnIndex].realtime.([]string)
	if ok == false {
		return  fmt.Errorf("internal error - variant type assertion failed.") 
	
	} 		
	(data)[row_index] = variant.value.(string) 


	value, ok := variant.value.(string)
	if ok == false{
		return fmt.Errorf("Column [%s] is indexed and value in variant [%v] must be string.", 
							variant.columnKey, variant.value)
	} 

	if table.descriptor.ColumnDescriptors[variant.columnIndex].indexed == true {	
		table.columns[variant.columnIndex].index[value] = row_index
	}

	return nil
}

func insertInt64 (table *table, row_index uint64, variant *variant) error {

	data, ok := table.columns[variant.columnIndex].realtime.([]int64)
	if ok == false {
		return fmt.Errorf("internal error - variant type assertion failed.") 
	} 		
	(data)[row_index] = variant.value.(int64) 

	return nil
}

func insertFloat64 (table *table, row_index uint64, variant *variant) error {

	data, ok := table.columns[variant.columnIndex].realtime.([]float64)
	if ok == false {
		return fmt.Errorf("internal error - variant type assertion failed.") 
	} 		
	(data)[row_index] = variant.value.(float64) 

	return nil
}



/******************************************************************************
update functions
******************************************************************************/
func updateUndefined (tbl *table, row_index uint64, variant *variant) error {

	err := fmt.Errorf("internal error - type const not supported.") 
	return err
}

func updateVarUTF8 (tbl *table, row_index uint64, variant *variant) error {
	
	data, ok := tbl.columns[variant.columnIndex].realtime.([]string)
	if ok == false {
		return fmt.Errorf("internal error - variant type assertion failed.") 
	} 		
	(data)[row_index] = variant.value.(string) 
	
	return nil
}

func updateInt64 (tbl *table, row_index uint64, variant *variant) error {
	data, ok := tbl.columns[variant.columnIndex].realtime.([]int64)
	if ok == false {
		return fmt.Errorf("internal error - variant type assertion failed.") 
	} 		
	(data)[row_index] = variant.value.(int64) 

	return nil
}

func updateFloat64 (tbl *table, row_index uint64, variant *variant) error {
	data, ok := tbl.columns[variant.columnIndex].realtime.([]float64)
	if ok == false {
		return fmt.Errorf("internal error - variant type assertion failed.") 
	}
	(data)[row_index] = variant.value.(float64) 

	return nil
}

/******************************************************************************
delete functions
******************************************************************************/
func deleteUndefined (tbl *table, row_index uint64, col_index uint16) error {

	err := fmt.Errorf("internal error - type const not supported.") 
	return err
}

func deleteVarUTF8 (tbl *table, row_index uint64, col_index uint16) error {

	data, ok := tbl.columns[col_index].realtime.([]string)
	if ok == false {
		return fmt.Errorf("internal error - variant type assertion failed.") 
	} 		
	(data)[row_index] = ""
	
	return nil
}

func deleteInt64 (tbl *table, row_index uint64, col_index uint16) error {
	data, ok := tbl.columns[col_index].realtime.([]int64)
	if ok == false {
		return fmt.Errorf("internal error - variant type assertion failed.") 
	} 		
	(data)[row_index] = int64(0)

	return nil
}

func deleteFloat64 (tbl *table, row_index uint64, col_index uint16) error {
	data, ok := tbl.columns[col_index].realtime.([]float64)
	if ok == false {
		return fmt.Errorf("internal error - variant type assertion failed.") 
	}
	(data)[row_index] = float64(0.0) 

	return nil
}



/******************************************************************************
validate functions
******************************************************************************/
func validateUndefined(col_desc * ColumnDescriptor, variant *variant) error {
	err := fmt.Errorf("internal error - type const not supported.") 
	return err
}

func validateVarUTF8(col_desc * ColumnDescriptor, variant *variant) error {
	return nil
}

func validateInt64(col_desc * ColumnDescriptor, variant *variant) error {
	return nil
}

func validateFloat64(col_desc * ColumnDescriptor, variant *variant) error {
	return nil
}