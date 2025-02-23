package rtdsms

import (
	"fmt"
//		"math"
//		"os"
//		"reflect"
//		"regexp"
//		"runtime"
//	"time"
)


/* ----------------------------------------------------------------------------
TableDescriptor
-----------------------------------------------------------------------------*/
type table_descriptor struct {
	realm_name			string
	domain_name			string
	datastore_name		string
	schema_name			string
	table_name			string
	total_rows			uint64
	columns_map			map[string]uint16
	column_descriptors	[]*ColumnDescriptor	
}

func NewTableDescriptor(realm_name string,
						domain_name string,
						datastore_name string,
						schema_name string, 
						table_name string, 
						total_rows uint64, 
						column_descriptors []*ColumnDescriptor) (*table_descriptor, error ){

	// TODO: Add name validation, etc
	tblDesc := new (table_descriptor)

	tblDesc.realm_name 		= realm_name
	tblDesc.domain_name 	= domain_name
	tblDesc.datastore_name 	= datastore_name
	tblDesc.schema_name 	= schema_name
	tblDesc.table_name 		= table_name
	tblDesc.total_rows 		= total_rows

	// -------------------------------------------------------------------------
	// make the column map
	// -------------------------------------------------------------------------
	tblDesc.columns_map = make (map[string]uint16)


	// -------------------------------------------------------------------------
	// This is kind of a kluge but it works. As this function is called by both 
	// DDL (create table) and the construct function, you can end up with two 
	// row_id columns. But, the
	// DDL function call with column_descriptors defined while the construct
	// function calls with column descriptors set to nil. So.... 
	// -------------------------------------------------------------------------
	if column_descriptors != nil {
		row_id_cd, err := NewColumnDescriptor("row_id", "varUTF8", INDEXED, uint64(max_object_name_length))
		if err != nil {		// This should never happen!
			return nil, fmt.Errorf("Internal error: %s", err)
		}

		// add the ROW_ID  at the beginning of the column descriptors and column map
		tblDesc.columns_map["row_id"] = uint16(len(tblDesc.columns_map))
		tblDesc.column_descriptors = append(tblDesc.column_descriptors, row_id_cd)
	}

	// -------------------------------------------------------------------------
	// Build the remaining column map.
	// -------------------------------------------------------------------------
	for i, desc := range column_descriptors {
		tblDesc.columns_map[desc.name] = uint16(len(tblDesc.columns_map))
		tblDesc.column_descriptors = append(tblDesc.column_descriptors, column_descriptors[i])
		
		// fmt.Printf("Creating column %s\n", desc.name)
	}

	

	return tblDesc, nil
}

