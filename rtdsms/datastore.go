

package rtdsms

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
//	"os"
//	"time"
)



/* ----------------------------------------------------------------------------
datastore is the functional equivelant of a 'database' in many environments.
We call it a data store since it will also have managed file system access,
key value pairs, and objects, in addition to the normal table based storage
of many database systems.

datastore is the interface that all the methods wrap around.
-----------------------------------------------------------------------------*/
type datastore struct {
	descriptor					datastoreDescriptor
//	name						string
	state 						uint8
	path						string
	requests 					chan( * request)
/*	tables						[]*table
	tableMap					map[string]uint64 */
	schemas						[]*schema
	schemas_map 				map[string]uint64
	snapshot_descriptors		[]*SnapshotDescriptor
	snapshot_descriptors_map	map[string]uint64
	context_descriptors			[]*ContextDescriptor
	context_descriptors_map		map[string]uint64
}


/* ---------------------------------------------------------------------------
Start()  
---------------------------------------------------------------------------- */
func (datastore * datastore ) Start() {
	
	go datastoreLoop(datastore)

}


func (datastore * datastore ) Stop() {
//	for _, ds := range rtdsms.datastores{
//		ds.Stop()
//	}
}

/* ----------------------------------------------------------------------------
 Each data store has a loop that runs continuoulsy to  service datastore 
 requests.
-----------------------------------------------------------------------------*/
func datastoreLoop(ds * datastore) {

	// Create an array of functions that provide data store services
	var funcs [100]func( *datastore, * request )           

	//First fill in the table with dummy funcs
	for i := 0; i < 100; i++ {
		funcs[i] = datastore_dummy
	}

	// -------------------------------------------------------------------------
	// Now load the array with implemented functions
	// -------------------------------------------------------------------------
	
	// data manipulation functions
	funcs[func_id_select] 		= select_ // 'select' is a GO key word
	funcs[func_id_insert] 		= insert
	funcs[func_id_update] 		= update
	funcs[func_id_delete] 		= delete

	// persistence functions
	funcs[func_id_take_snapshot] 		= takeSnapshot
	funcs[func_id_restore_snapshot] 	= restoreSnapshot
	funcs[func_id_clear_snapshot] 		= clearSnapshot
	funcs[func_id_clear_realtime] 		= clearRealtime
	funcs[func_id_write_snapshot] 		= writeSnapshot
	funcs[func_id_read_snapshot] 		= readSnapshot
//	funcs[func_id_write_context] 		= writeContext
//	funcs[func_id_read_context] 		= readContext

	// ... TODO: add the Data Control API functions

	ds.requests = make(chan * request)

	for true {
		request := <- ds.requests // Blocking receive
	
		if request.function_id > 100 {
			request.err = fmt.Errorf("Function_id [%d] is invalid - too large.", request.function_id)
		} else {
			funcs[request.function_id](ds, request)
		}

		// nil is ALWAYS returned through the channel, simply as a sync signal
		*request.response_channel <- nil
	}
}



// -----------------------------------------------------------------------------
// construct   () 
// -----------------------------------------------------------------------------
func (datastore * datastore) construct(filename string) error  {

//	fmt.Printf("Constructing datastore [%s] [%s] [%s]...\n", datastore.descriptor.realm_name, datastore.descriptor.domain_name, datastore.descriptor.datastore_name)

	if datastore.schemas_map == nil {
		datastore.schemas_map = make(map[string]uint64)
	}
	if datastore.snapshot_descriptors_map == nil {
		datastore.snapshot_descriptors_map = make(map[string]uint64)
	}



	// Open the sqlite file
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
					datastore.descriptor.domain_name).Scan(&domain_id)
	if result.Error != nil {
		return result.Error
	}

	// ------------------------------------------------------------------------
	// Get the datastore_id
	// ------------------------------------------------------------------------
	var datastore_id	int64
	result = db.Raw("SELECT datastore_id FROM datastores where domain_id = ? AND datastore_name = ?", 
					domain_id, datastore.descriptor.datastore_name).Scan(&datastore_id)
	if result.Error != nil {
		return result.Error
	}

	// ------------------------------------------------------------------------
	// Select the snapshots
	// ------------------------------------------------------------------------
	snapshot_rows, err := sqlDB.Query("SELECT snapshot_name FROM snapshots WHERE domain_id = ? AND datastore_id = ? order by snapshot_id", domain_id, datastore_id)
    if err != nil {
        return err
    }
    defer snapshot_rows.Close()

    // Loop through rows, using Scan to assign column data to struct fields.
    for snapshot_rows.Next() {
		var snapshot_name string
        err := snapshot_rows.Scan(&snapshot_name)
		if err != nil {
            return fmt.Errorf("%v\n", err)
        } else {
				snapshot := SnapshotDescriptor {
				name : snapshot_name,
			}
			datastore.snapshot_descriptors_map[snapshot_name] = uint64(len(datastore.snapshot_descriptors))
			datastore.snapshot_descriptors = append(datastore.snapshot_descriptors, &snapshot)
		}
	}

	


	// ------------------------------------------------------------------------
	// Select the schemas
	// ------------------------------------------------------------------------
	rows, err := sqlDB.Query("SELECT schema_name FROM schemas WHERE domain_id = ? AND datastore_id = ? order by schema_id", domain_id, datastore_id)
    if err != nil {
        return err
    }
    defer rows.Close()

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
		var name string
        err := rows.Scan(&name)
		if err != nil {
            return fmt.Errorf("%v\n", err)
        } else {
			schema := schema {
			}
			schema.descriptor.schema_name = name
			schema.descriptor.datastore_name = datastore.descriptor.datastore_name
			schema.descriptor.domain_name = datastore.descriptor.domain_name
			schema.descriptor.realm_name = datastore.descriptor.realm_name
			datastore.schemas_map[name] = uint64(len(datastore.schemas_map))
			datastore.schemas = append(datastore.schemas, &schema)
		}
    }
    if err = rows.Err(); err != nil {
		return fmt.Errorf("%v\n", err)
    }

	for _, schema := range datastore.schemas {
		err = schema.construct(filename)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}


//	fmt.Printf("datastore done.\n")

	return nil
}
