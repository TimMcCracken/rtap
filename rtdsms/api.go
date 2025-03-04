package rtdsms


import (
	"fmt"
//	"math"
//	"os"
//	"reflect"
//	"regexp"
//	"runtime"
//	"time"
)


/* Other planned RTAP APIs
	Time services inlcuding timers and schedulers

	Batch services - primarily aimed at housekeeping and reports

	Message Queuing - for event driven inter go-routine communications

*/



/* ----------------------------------------------------------------------------
datastore APIs

There are 3 data store APIs. 
	DC_API Data Control API is used to start and stop databases, and in the
	future will control security settings.

	DD_API Data Definiton API - used to create and drop various objects
	in the data store

	DM_API Data Manipulation API performs the typical CRUD functions on the
	objects in the datastore

The datastore manages several different kinds of data:
	Tables - consisting of multiple columns. Each column can be a different
			data type, but every value in a column is the same type.

	Key-Value Group - a key is a string and the value may be of any type.
			multiple key values are stored in a single group.
	
	Arrays = a list of values of the same type. Arrays can be keyed or not.

	Matrix - A two dimensional array with all values of the same type

	Record Management Server - provides record based disk storage. Also used
	for persistence for all the in-memory objects

	Managed File Server - a managed version of the standard file services
	provided by *nix systems, enable replication and backup to be handled by
	the data store

-----------------------------------------------------------------------------*/

// request is used by all 3 APIs.

type Request struct {
	functionID	 		uint16
	dbKey 				string 
	objectKey 			string 
	responseChannel 	*chan(error)
	data				any
	err					error
//	info				string
}

// ----------------------------------------------------------------------------
// DM_API
// ----------------------------------------------------------------------------
// structures used by the DM API



// ----------------------------------------------------------------------------
// The Select(), Update() and Insert() methods are all called with a 'tuples'
// pointer argument that is used for both the input and the output as needed.
// ----------------------------------------------------------------------------

func (ds * Datastore) Select( records []*Record) (error) {

	// Send the request to the datastore server go routine via a channel
	if ds.requests == nil {
		return fmt.Errorf("datastore is nil.")
	}
	
	// Create a channel for the response
	resp_chan := make( chan(error) )
	
	// Create a request
	req := Request{	functionID: func_id_select, 
						responseChannel: &resp_chan,
						data: records} 
		
	// Send the request to the datastore server go routine via a channel
	ds.requests <- &req 
	
	// wait on the response to come back and return to the calling program
	_ = <- resp_chan
	
	if req.err != nil {
		return req.err
	}
	return nil	
}


func (ds * Datastore) Update( records []*Record) (error) {

		// Send the request to the datastore server go routine via a channel
		if ds.requests == nil {
			return fmt.Errorf("datastore request channel is nil.")
		}
	
		// Create a channel for the response
		resp_chan := make( chan(error) )
	
		// Create a request
		req := Request{	functionID: func_id_update, 
						responseChannel: &resp_chan,
						data: records} 
		
		// Send the request to the datastore server go routine via a channel
		ds.requests <- &req 
	
		// wait on the response to come back and return to the calling program
		_ = <- resp_chan
	
		if req.err != nil {
			return req.err
		}
		return nil
}

func (ds * Datastore) Insert( records []*Record) (error) {

	// Send the request to the datastore server go routine via a channel
	if ds.requests == nil {
		return fmt.Errorf("datastore request channel is nil.")
	}

	// Create a channel for the response
	resp_chan := make( chan(error) )

	// Create a request
	req := Request{	functionID: func_id_insert, 
					responseChannel: &resp_chan,
					data: records} 
	
	// Send the request to the datastore server go routine via a channel
	ds.requests <- &req 

	// wait on the response to come back and return to the calling program
	_ = <- resp_chan

	if req.err != nil {
		return req.err
	}
	return nil
}


func (ds * Datastore) Delete( records []*Record) (error) {

	// Send the request to the datastore server go routine via a channel
	if ds.requests == nil {
		return fmt.Errorf("datastore is nil.")
	}

	// Create a channel for the response
	resp_chan := make( chan(error) )

	// Create a request
	req := Request{	functionID: func_id_delete, 
					responseChannel: &resp_chan,
					data: records} 
	
	// Send the request to the datastore server go routine via a channel
	ds.requests <- &req 

	// wait on the response to come back and return to the calling program
	_ = <- resp_chan

	if req.err != nil {
		return req.err
	}
	return nil
}








/******************************************************************************
persistence API functions
******************************************************************************/

// -----------------------------------------------------------------------------
// TakeSnapshot copies data from the real-time data sources to a snapshot
// buffer. If the schema parameter is "*" then all tables in all schemas are
// snapshotted.  If the schema is not "*", then it must be the name of an
// existing schema. Then if table is "*", all tables in that schema are 
// snapshotted.
// -----------------------------------------------------------------------------
func (ds * Datastore ) TakeSnapshot(snapshot_name string, 
									schema_name string, 
									table_name string) error {
	
	// -------------------------------------------------------------------------
	// Create a channel for the response
	// -------------------------------------------------------------------------
	resp_chan := make( chan(error) )

	// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	snapshotName : snapshot_name,
						schemaName : schema_name,
						tableName  : table_name}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := Request{	functionID: func_id_take_snapshot, 
					responseChannel: &resp_chan,
					data : &data} 
	req.responseChannel = &resp_chan

	// -------------------------------------------------------------------------
	// Send the request to the datastore server go routine via a channel
	// -------------------------------------------------------------------------
	ds.requests <- &req 

	// -------------------------------------------------------------------------
	// wait on the response to come back and return to the calling program
	// -------------------------------------------------------------------------
	_ = <- resp_chan

	if req.err != nil {
		return req.err
	}
	return nil
}

func (ds * Datastore ) RestoreSnapshot(	snapshot_name string, 
										schema_name string, 
										table_name string) error {
	// -------------------------------------------------------------------------
	// Create a channel for the response
	// -------------------------------------------------------------------------
	resp_chan := make( chan(error) )

	// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	snapshotName : snapshot_name,
						schemaName : schema_name,
						tableName  : table_name}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := Request{	functionID: func_id_restore_snapshot, 
					responseChannel: &resp_chan,
					data : &data} 
	req.responseChannel = &resp_chan

	// -------------------------------------------------------------------------
	// Send the request to the datastore server go routine via a channel
	// -------------------------------------------------------------------------
	ds.requests <- &req 

	// -------------------------------------------------------------------------
	// wait on the response to come back and return to the calling program
	// -------------------------------------------------------------------------
	_ = <- resp_chan

	if req.err != nil {
		return req.err
	}
	return nil
}

func (ds * Datastore ) WriteSnapshot(	snapshot_name string, 
										schema_name string, 
										table_name string,
										path_file string) error {
	// Create a channel for the response
	resp_chan := make( chan(error) )

	// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	snapshotName : snapshot_name,
						schemaName : schema_name,
						tableName  : table_name,
						pathFile : path_file}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := Request{	functionID: func_id_write_snapshot, 
					responseChannel: &resp_chan,
					data : &data} 
	req.responseChannel = &resp_chan

	// -------------------------------------------------------------------------
	// Send the request to the datastore server go routine via a channel
	// -------------------------------------------------------------------------
	ds.requests <- &req 

	// -------------------------------------------------------------------------
	// wait on the response to come back and return to the calling program
	// -------------------------------------------------------------------------
	_ = <- resp_chan

	if req.err != nil {
		return req.err
	}
	return nil
}


func (ds * Datastore ) ReadSnapshot(snapshot_name string, 
									schema_name string, 
									table_name string,
									path_file string) error {
	// -------------------------------------------------------------------------
	// Create a channel for the response
	// -------------------------------------------------------------------------
	resp_chan := make( chan(error) )

	// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	snapshotName : snapshot_name,
						schemaName : schema_name,
						tableName  : table_name,
						pathFile : path_file}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := Request{	functionID: func_id_read_snapshot, 
					responseChannel: &resp_chan,
					data : &data} 
	req.responseChannel = &resp_chan

	// -------------------------------------------------------------------------
	// Send the request to the datastore server go routine via a channel
	// -------------------------------------------------------------------------
	ds.requests <- &req 

	// -------------------------------------------------------------------------
	// wait on the response to come back and return to the calling program
	// -------------------------------------------------------------------------
	_ = <- resp_chan

	if req.err != nil {
		return req.err
	}
	return nil
}

func (ds * Datastore ) ClearSnapshot(	snapshot_name string, 
										schema_name string, 
										table_name string) error {

	// -------------------------------------------------------------------------
		// Create a channel for the response
	// -------------------------------------------------------------------------
	resp_chan := make( chan(error) )

	// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	snapshotName : snapshot_name,
						schemaName : schema_name,
						tableName  : table_name}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := Request{	functionID: func_id_clear_snapshot, 
					responseChannel: &resp_chan,
					data : &data} 
	req.responseChannel = &resp_chan

	// -------------------------------------------------------------------------
	// Send the request to the datastore server go routine via a channel
	// -------------------------------------------------------------------------
	ds.requests <- &req 

	// -------------------------------------------------------------------------
	// wait on the response to come back and return to the calling program
	// -------------------------------------------------------------------------
	_ = <- resp_chan

	if req.err != nil {
		return req.err
	}
	return nil
}

func (ds * Datastore ) ClearRealtime(	schema_name string, 
										table_name string) error {
	// -------------------------------------------------------------------------
	// Create a channel for the response
	// -------------------------------------------------------------------------
	resp_chan := make( chan(error) )

// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	schemaName : schema_name,
						tableName  : table_name}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := Request{	functionID: func_id_clear_snapshot, 
					responseChannel: &resp_chan,
					data : &data} 
	req.responseChannel = &resp_chan

	// -------------------------------------------------------------------------
	// Send the request to the datastore server go routine via a channel
	// -------------------------------------------------------------------------
	ds.requests <- &req 

	// -------------------------------------------------------------------------
	// wait on the response to come back and return to the calling program
	// -------------------------------------------------------------------------
	_ = <- resp_chan

	if req.err != nil {
		return req.err
	}
	return nil
}