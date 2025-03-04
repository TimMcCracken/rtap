

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
type Datastore struct {
	descriptor					datastoreDescriptor
//	name						string
	state 						uint8
	path						string
	requests 					chan( * Request)
/*	tables						[]*table
	tableMap					map[string]uint64 */
	schemas						[]*schema
	schemasMap 					map[string]uint64
	snapshotDescriptors			[]*SnapshotDescriptor
	snapshotDescriptorsMap		map[string]uint64
	contextDescriptors			[]*ContextDescriptor
	contextDescriptorsMap		map[string]uint64
}


// -----------------------------------------------------------------------------
// NewDatastore() 
// -----------------------------------------------------------------------------
func NewDatastore(realm_name string, domain_name string, datastore_name string) ( * Datastore, error) {


	// TODO: validate object names

	ds := Datastore {
	}
	ds.descriptor.datastoreName = datastore_name
	ds.descriptor.domainName = domain_name
	ds.descriptor.realmName = realm_name
	return &ds, nil
}





/* ---------------------------------------------------------------------------
Start()  
---------------------------------------------------------------------------- */
func (datastore * Datastore ) Start() {
	
	go datastoreLoop(datastore)

}




func (datastore * Datastore ) Stop() {
//	for _, ds := range rtdsms.datastores{
//		ds.Stop()
//	}
}

func datastore_dummy (ds * Datastore, request * Request) {
	//	fmt.Printf("%s Request: Function: %d   DB: %s\n", ds.name, request.functionID, request.db_key)
		request.err = fmt.Errorf("functionID [%d] is not implemented.", request.functionID)
}


// -----------------------------------------------------------------------------
// ping() is a diagnostic function that is used bu the WatchDog timer
// -----------------------------------------------------------------------------
func datastore_ping (ds * Datastore, request * Request) {
	request.err = nil
}


/* ----------------------------------------------------------------------------
 Each data store has a loop that runs continuoulsy to  service datastore 
 requests.
-----------------------------------------------------------------------------*/
func datastoreLoop(ds * Datastore) {

	// Create an array of functions that provide data store services
	var funcs [100]func( *Datastore, * Request )           

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

	ds.requests = make(chan * Request)

	for true {
		request := <- ds.requests // Blocking receive
	
		if request.functionID > 100 {
			request.err = fmt.Errorf("functionID [%d] is invalid - too large.", request.functionID)
		} else {
			funcs[request.functionID](ds, request)
		}

		// nil is ALWAYS returned through the channel, simply as a sync signal
		*request.responseChannel <- nil
	}
}



// -----------------------------------------------------------------------------
// construct   () 
// -----------------------------------------------------------------------------
func (ds * Datastore) Construct(filename string) error  {

//	fmt.Printf("Constructing datastore [%s] [%s] [%s]...\n", datastore.descriptor.realmName, datastore.descriptor.domainName, datastore.descriptor.datastoreName)

	if ds.schemasMap == nil {
		ds.schemasMap = make(map[string]uint64)
	}
	if ds.snapshotDescriptorsMap == nil {
		ds.snapshotDescriptorsMap = make(map[string]uint64)
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
	result := db.Raw("SELECT domain_id FROM domains where domainName = ?", 
					ds.descriptor.domainName).Scan(&domain_id)
	if result.Error != nil {
		return result.Error
	}

	// ------------------------------------------------------------------------
	// Get the datastore_id
	// ------------------------------------------------------------------------
	var datastore_id	int64
	result = db.Raw("SELECT datastore_id FROM datastores where domain_id = ? AND datastoreName = ?", 
					domain_id, ds.descriptor.datastoreName).Scan(&datastore_id)
	if result.Error != nil {
		return result.Error
	}

	// ------------------------------------------------------------------------
	// Select the snapshots
	// ------------------------------------------------------------------------
	snapshot_rows, err := sqlDB.Query("SELECT snapshotName FROM snapshots WHERE domain_id = ? AND datastore_id = ? order by snapshot_id", domain_id, datastore_id)
    if err != nil {
        return err
    }
    defer snapshot_rows.Close()

    // Loop through rows, using Scan to assign column data to struct fields.
    for snapshot_rows.Next() {
		var snapshotName string
        err := snapshot_rows.Scan(&snapshotName)
		if err != nil {
            return fmt.Errorf("%v\n", err)
        } else {
				snapshot := SnapshotDescriptor {
				name : snapshotName,
			}
			ds.snapshotDescriptorsMap[snapshotName] = uint64(len(ds.snapshotDescriptors))
			ds.snapshotDescriptors = append(ds.snapshotDescriptors, &snapshot)
		}
	}

	


	// ------------------------------------------------------------------------
	// Select the schemas
	// ------------------------------------------------------------------------
	rows, err := sqlDB.Query("SELECT schemaName FROM schemas WHERE domain_id = ? AND datastore_id = ? order by schema_id", domain_id, datastore_id)
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
			schema.descriptor.schemaName = name
			schema.descriptor.datastoreName = ds.descriptor.datastoreName
			schema.descriptor.domainName = ds.descriptor.domainName
			schema.descriptor.realmName = ds.descriptor.realmName
			ds.schemas = append(ds.schemas, &schema)
		}
    }
    if err = rows.Err(); err != nil {
		return fmt.Errorf("%v\n", err)
    }

	for _, schema := range ds.schemas {
		err = schema.construct(filename)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}


//	fmt.Printf("datastore done.\n")

	return nil
}
