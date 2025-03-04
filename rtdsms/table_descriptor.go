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
type TableDescriptor struct {
	realmName			string
	domainName			string
	datastoreName		string
	schemaName			string
	tableName			string
	totalRows			uint64
	columnsMap			map[string]uint16
	ColumnDescriptors	[]*ColumnDescriptor	
}

func NewTableDescriptor(realmName string,
						domainName string,
						datastoreName string,
						schemaName string, 
						tableName string, 
						totalRows uint64, 
						columnDescriptors []*ColumnDescriptor) (*TableDescriptor, error ){

	// TODO: Add name validation, etc
	tblDesc := new (TableDescriptor)

	tblDesc.realmName 		= realmName
	tblDesc.domainName 		= domainName
	tblDesc.datastoreName 	= datastoreName
	tblDesc.schemaName 		= tableName
	tblDesc.totalRows 		= totalRows

	// -------------------------------------------------------------------------
	// make the column map
	// -------------------------------------------------------------------------
	tblDesc.columnsMap = make (map[string]uint16)

	// -------------------------------------------------------------------------
	// This is kind of a kluge but it works. As this function is called by both 
	// DDL (create table) and the construct function, you can end up with two 
	// row_id columns. But, the
	// DDL function call with column_descriptors defined while the construct
	// function calls with column descriptors set to nil. So.... 
	// -------------------------------------------------------------------------
	if columnDescriptors != nil {
		row_id_cd, err := NewColumnDescriptor("row_id", "varUTF8", INDEXED, uint64(max_object_name_length))
		if err != nil {		// This should never happen!
			return nil, fmt.Errorf("Internal error: %s", err)
		}

		// add the ROW_ID  at the beginning of the column descriptors and column map
		tblDesc.columnsMap["row_id"] = uint16(len(tblDesc.columnsMap))
		tblDesc.ColumnDescriptors = append(tblDesc.ColumnDescriptors, row_id_cd)
	}

	// -------------------------------------------------------------------------
	// Build the remaining column map.
	// -------------------------------------------------------------------------
	for i, desc := range columnDescriptors {
		tblDesc.columnsMap[desc.name] = uint16(len(tblDesc.columnsMap))
		tblDesc.ColumnDescriptors = append(tblDesc.ColumnDescriptors, columnDescriptors[i])
		
		// fmt.Printf("Creating column %s\n", desc.name)
	}

	

	return tblDesc, nil
}


func (td * TableDescriptor) Name() string {
	return td.tableName
}

func (td * TableDescriptor) TotalRows() uint64 {
	return td.totalRows
}
