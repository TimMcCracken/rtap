package rtdsms

import (
		"fmt"
//		"gorm.io/gorm"
//		"math"
//		"os"
//		"reflect"
//		"regexp"
//		"runtime"
//		"time"
//		"rtap/realm"
)






/* ----------------------------------------------------------------------------
definitions:
	rtdb_api = the object used by an application to operate on a goupr of
	datastores that are stored together on a single machine. 

	rtdb: a collection of data within schemas  that can be operated on 
	together

	schema: a collection of tables, KV pairs or other objectsthat utilize a
	common set of attributes. For example, all configuration data might be
	stored in a schema called 'config' because the objects in it don't
	change often and thereby don't need to be written to disk often 
	whereas frequently changing real-time data might be stored in a schem
	call "real-time" which is ferquently written to disk

--------------------------------------------------------------------------   */

// THese need to be moved to the 'options.go' file for each project
const max_datastores 	= 8		// max datastores per realm
const max_tables		= 100	// max tables per data store
const max_columns		= 100	// max columns per table
const max_snapshots		= 4		// max snapshots per datastore
const max_contexts		= 8		// max contexts per datastore

const default_channel_depth	= 16

const max_object_name_length	= 16
const object_name_format = "^[a-z][a-z0-9_]*$"

// Money type?
// time type(s)
// duration type(s)
// engineering units (float64 or float32 with EU)
// phasor (angle + magnitude)
// 2 dimension coordinates
// 3 dimension coordinates
// 3 dimension coordinates + velocity

const max_column_types = 64


const rtds_func_id_ping		= 0

// Data Definition API Functions
const FID_CREATE_DOMAIN		= 1
const FID_CREATE_DATASTORE	= 1
const FID_CREATE_SCHEMA		= 2
const FID_CREATE_TABLE		= 10

const func_id_drop_table		= 11
const func_id_create_kv_group	= 14
const func_id_drop_kv_group		= 15
const func_id_create_array		= 16
const func_id_drop_array		= 17
const func_id_create_matrix		= 18
const func_id_drop_matrix		= 19

// Data Manipulation API Functions
const func_id_select	= 30
const func_id_insert	= 31
const func_id_update	= 32
const func_id_delete	= 33
//const func_id_execute	= 34


// persistence functions
const func_id_take_snapshot		= 40	// copies the realtime datas to a snapshot
const func_id_restore_snapshot	= 41	// copies a snapshot to the realtime data
const func_id_clear_snapshot	= 42	// clears the data from a snapshot
const func_id_clear_realtime	= 43
const func_id_write_snapshot	= 44	// writes a snapshot to disk
const func_id_read_snapshot		= 45	// read a snapshot from disk
const func_id_write_context		= 46	// writes data from a context to disk
const func_id_read_context		= 47	// read from disk into a context 


// domain states (& datastroe states?)
const DS_OFFLINE	= 0
const DS_STARTING	= 1
const DS_ONLINE		= 2
const DS_STOPPING	= 3
const DS_PAUSED		= 4

// snapshot states
const SS_STATE_CLEAR	= 0
const SS_STATE_VALID	= 1
const SS_STATE_BUSY		= 2


/* Below here is exported. Some of the constants above probably need to
be unexported */







type SnapshotDescriptor struct {
	name	string
}

func NewSnapshotDescriptor(name string) (*SnapshotDescriptor, error ){

	err := validateObjectName(name)
	if err != nil {
		return nil, err
	}

	newDesc := new(SnapshotDescriptor)
	newDesc.name = name
	return newDesc, nil

}

func (sd * SnapshotDescriptor) Name() string {
	return sd.name
}






// ContextDescriptor is for future use
type ContextDescriptor struct {
	name	string
}




func GetState(ds * Datastore) string {
	switch(ds.state) {
	case 0:
		return "Offline"
	case 1:
		return "Starting"
	case 2:
		return "Online"
	case 3:
		return "Stopping"
	case 4:
		return "Paused"
	}
	return fmt.Sprintf("Unknown state: [%d]", ds.state)
}



/* not necessary ?
// ----------------------------------------------------------------------------
// realm loop
// ----------------------------------------------------------------------------
func realmLoop(ds * datastore) {

	// Create an array of functions that provide data store services
	var funcs [100]func( *realm, * request )           

	//First fill in the table with dummy funcs
	for i := 0; i < 100; i++ {
		funcs[i] = realm_dummy
	}


}
*/




/*
func (ds * datastore) PrintSchema() {


	for i, tbl := range ds.tables {
		fmt.Printf("    Table #: %d  Name: %s\n", i, tbl.descriptor.name)

		for j, column := range tbl.columns {
			fmt.Printf("      Column #: %5d   Name: %-16s  Type: %-8s\n", j, column.descriptor.name, column.descriptor.typeName)

		}
	}
}
*/




/* ----------------------------------------------------------------------------
  The following are the functions served by the data store loop for the 
  Datastore API and the Realm API
-----------------------------------------------------------------------------*/

/*
func realm_dummy (r * realm, request * request) {
	//	fmt.Printf("%s Request: Function: %d   DB: %s\n", ds.name, request.function_id, request.db_key)
		request.err = fmt.Errorf("Function_id [%d] is not implemented.", request.function_id)
}
*/



/*	
// ping() is a diagnostic function that is used bu the WatchDog timer
func realm_ping (ds * datastore, request * request) {
		request.err = nil
}
		*/
