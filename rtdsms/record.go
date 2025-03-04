package rtdsms

import (
	"fmt"
	"math"
)


// Record is used to hold variants and other data used by the select, insert
// and update functions.
type Record struct {
		schemaKey				string
		schemaIndex				uint64
		tableKey				string
		tableIndex				uint64
		lookupColumnKey			string 		// The name of the column that we will use to lookup the row
		lookupColumnIndex		uint16		// the index of the column that we will use to lookkup the row
		lookupRowKey			string
		lookupRowIndex			uint64		// the index of the row we looked up
		variants				[]*variant
		indexesValid			bool
		cacheTimeout			uint64		// time in milliseconds
} 
	

func (record * Record) init() {
	
	record.schemaIndex 			= math.MaxUint64
	record.tableIndex 			= math.MaxUint64
	record.lookupColumnIndex	= math.MaxUint16
	record.lookupRowIndex 		= math.MaxUint64

	record.indexesValid = false
	record.variants = nil
	record.cacheTimeout = 0
}



// ----------------------------------------------------------------------------
// NewRecord()
// ----------------------------------------------------------------------------
func NewRecord(schema_key string, table_key string, lookup_column_key string, 
				lookup_row_key string, fields []string, cache_timeout uint64) * Record {
		
	record := new(Record)
	record.init()

	record.schemaKey = schema_key
	record.tableKey = table_key
	record.lookupColumnKey = lookup_column_key
	record.lookupRowKey = lookup_row_key
	record.cacheTimeout = cache_timeout

	for _, field := range fields {
		variant := new(variant)
		variant.columnKey = field
		record.variants = append(record.variants, variant)			
	}
	return record
}

// ----------------------------------------------------------------------------
// Reset() allows a record to be reused with a different data source
// ----------------------------------------------------------------------------
func (record * Record) Reset(schema_key string, table_key string, 
								lookup_column_key string, lookup_row_key string, 
								fields []string, cache_timeout uint64) {

	record.init()

	record.schemaKey = schema_key
	record.tableKey = table_key
	record.lookupColumnKey = lookup_column_key
	record.lookupRowKey = lookup_row_key
	record.cacheTimeout = cache_timeout
	
	for _, field := range fields {
		variant := new(variant)
		variant.columnKey = field
		record.variants = append(record.variants, variant)			
	}
}




func (record * Record)AppendVariant(column_key string, value any){
		variant := new(variant)
		variant.columnKey = column_key
		variant.value = value
		record.variants = append(record.variants, variant)
}
	


func (record * Record)GetValue(index int) (error, any){
	
	if index < 0 {
		return fmt.Errorf("Index value [%d] not allowed", index), nil
	}
	if index > len(record.variants) {
		return fmt.Errorf("Index value [%d] too big. Only have [%d] variants.", index, len(record.variants)), nil
	}

	return nil, record.variants[index].value
}


func (record * Record)SetValue(index int, value any) (error){

	if index < 0 {
		return fmt.Errorf("Index value [%d] not allowed", index)
	}
	if index > len(record.variants) {
		return fmt.Errorf("Index value [%d] too big. Only have [%d] variants.", index, len(record.variants))
	}

	record.variants[index].value = value
	return nil
}


func (record * Record)GetValues() (error, []any){
	
	values := make([]any, 0, len(record.variants))  

	return nil, values
}


func (record * Record)SetValues(values []any) (error){

	// that the length of the values matches the number of variants
	if len(values) != len(record.variants) {
		return fmt.Errorf("Number of values supplied [%d] must match number of variants[%d]", len(values), len(record.variants))
	}
	
	//record.variants[index].value = value
	return nil
}
