package rtdsms

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
)



type schema struct {
	descriptor 	schemaDescriptor
	tables		[]*table
	tables_map 	map[string]uint64
}



/// -----------------------------------------------------------------------------
// construct   () 
// -----------------------------------------------------------------------------
func (schema * schema) construct(filename string) error  {

//	fmt.Printf("Constructing schema [%s] [%s] [%s] [%s]\n", 
//				schema.descriptor.realm_name, 
//				schema.descriptor.domain_name, 
//				schema.descriptor.datastore_name,
//				schema.descriptor.schema_name)

	if schema.tables_map == nil {
		schema.tables_map = make(map[string]uint64)
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
					schema.descriptor.domain_name).Scan(&domain_id)
	if result.Error != nil {
		return result.Error
	}

	// ------------------------------------------------------------------------
	// Get the datastore_id_id
	// ------------------------------------------------------------------------
	var datastore_id	int64
	result = db.Raw("SELECT datastore_id FROM datastores WHERE domain_id = ? AND datastore_name = ?", 
					domain_id, schema.descriptor.datastore_name).Scan(&datastore_id)
	if result.Error != nil {
		return result.Error
	}

	// ------------------------------------------------------------------------
	// Get the schema_id
	// ------------------------------------------------------------------------
	var schema_id	int64
	result = db.Raw("SELECT schema_id FROM schemas WHERE domain_id = ? AND datastore_id = ? AND schema_name = ?", 
					domain_id, datastore_id, schema.descriptor.schema_name).Scan(&schema_id)
	if result.Error != nil {
		return result.Error
	}

	rows, err := sqlDB.Query("SELECT table_name, total_rows FROM tables WHERE domain_id = ? and datastore_id = ? and schema_id = ? order by table_id", 
					domain_id, datastore_id, schema_id)
    if err != nil {
        return err
    }
    defer rows.Close()

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
		var table_name 		string
		var total_rows 		uint64
		
        err := rows.Scan(&table_name, &total_rows)
		if err != nil {
            return fmt.Errorf("%v\n", err)
        } else {
			table := table {
			}
			table.descriptor, err = NewTableDescriptor(schema.descriptor.realm_name,
												schema.descriptor.domain_name,
												schema.descriptor.datastore_name,
												schema.descriptor.schema_name, 
												table_name, total_rows, nil) 

			schema.tables_map[table_name] = uint64(len(schema.tables_map))
			schema.tables = append(schema.tables, &table)
		}
    }
    if err = rows.Err(); err != nil {
		return fmt.Errorf("%v\n", err)
    }

	for _, table := range schema.tables {
		err = table.construct(filename)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}
//	fmt.Printf("schema done.\n")

	return nil
}


/*
NewTableDescriptor(		realm_name string,
						domain_name string,
						datastore_name string,
						schema_name string, 
						table_name string, 
						total_rows uint64, 
						column_descriptors []*ColumnDescriptor )
*/