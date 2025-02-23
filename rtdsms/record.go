package rtdsms

import (
	"fmt"
	"math"
)


// Record is used to hold variants and other data used by the select, insert
// and update functions.
type Record struct {
		schema_key				string
		schema_index			uint64
		table_key				string
		table_index				uint64
		lookup_column_key		string 		// The name of the column that we will use to lookup the row
		lookup_column_index		uint16		// the index of the column that we will use to lookkup the row
		lookup_row_key			string
		lookup_row_index		uint64		// the index of the row we looked up
		variants				[]*variant
		indexes_valid			bool
		cache_timeout			uint64		// time in milliseconds
} 
	

func (record * Record) init() {
	
	record.schema_index 		= math.MaxUint64
	record.table_index 			= math.MaxUint64
	record.lookup_column_index	= math.MaxUint16
	record.lookup_row_index 	= math.MaxUint64

	record.indexes_valid = false
	record.variants = nil
	record.cache_timeout = 0
}



// ----------------------------------------------------------------------------
// NewRecord()
// ----------------------------------------------------------------------------
func NewRecord(schema_key string, table_key string, lookup_column_key string, 
				lookup_row_key string, fields []string, cache_timeout uint64) * Record {
		
	record := new(Record)
	record.init()

	record.schema_key = schema_key
	record.table_key = table_key
	record.lookup_column_key = lookup_column_key
	record.lookup_row_key = lookup_row_key
	record.cache_timeout = cache_timeout

	for _, field := range fields {
		variant := new(variant)
		variant.column_key = field
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

	record.schema_key = schema_key
	record.table_key = table_key
	record.lookup_column_key = lookup_column_key
	record.lookup_row_key = lookup_row_key
	
	for _, field := range fields {
		variant := new(variant)
		variant.column_key = field
		record.variants = append(record.variants, variant)			
	}
	record.cache_timeout = cache_timeout

}




func (record * Record)AppendVariant(column_key string, value any){
		variant := new(variant)
		variant.column_key = column_key
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
