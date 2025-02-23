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

type request struct {
	function_id 		uint16
	db_key 				string 
	object_key 			string 
	response_channel 	*chan(error)
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

func (ds * datastore) Select( records []*Record) (error) {

	// Send the request to the datastore server go routine via a channel
	if ds.requests == nil {
		return fmt.Errorf("datastore is nil.")
	}
	
	// Create a channel for the response
	resp_chan := make( chan(error) )
	
	// Create a request
	req := request{	function_id: func_id_select, 
						response_channel: &resp_chan,
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


func (ds * datastore) Update( records []*Record) (error) {

		// Send the request to the datastore server go routine via a channel
		if ds.requests == nil {
			return fmt.Errorf("datastore request channel is nil.")
		}
	
		// Create a channel for the response
		resp_chan := make( chan(error) )
	
		// Create a request
		req := request{	function_id: func_id_update, 
						response_channel: &resp_chan,
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

func (ds * datastore) Insert( records []*Record) (error) {

	// Send the request to the datastore server go routine via a channel
	if ds.requests == nil {
		return fmt.Errorf("datastore request channel is nil.")
	}

	// Create a channel for the response
	resp_chan := make( chan(error) )

	// Create a request
	req := request{	function_id: func_id_insert, 
					response_channel: &resp_chan,
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


func (ds * datastore) Delete( records []*Record) (error) {

	// Send the request to the datastore server go routine via a channel
	if ds.requests == nil {
		return fmt.Errorf("datastore is nil.")
	}

	// Create a channel for the response
	resp_chan := make( chan(error) )

	// Create a request
	req := request{	function_id: func_id_delete, 
					response_channel: &resp_chan,
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
func (ds * datastore ) TakeSnapshot(snapshot_name string, 
									schema_name string, 
									table_name string) error {
	
	// -------------------------------------------------------------------------
	// Create a channel for the response
	// -------------------------------------------------------------------------
	resp_chan := make( chan(error) )

	// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	snapshot_name : snapshot_name,
						schema_name : schema_name,
						table_name  : table_name}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := request{	function_id: func_id_take_snapshot, 
					response_channel: &resp_chan,
					data : &data} 
	req.response_channel = &resp_chan

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

func (ds * datastore ) RestoreSnapshot(	snapshot_name string, 
										schema_name string, 
										table_name string) error {
	// -------------------------------------------------------------------------
	// Create a channel for the response
	// -------------------------------------------------------------------------
	resp_chan := make( chan(error) )

	// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	snapshot_name : snapshot_name,
						schema_name : schema_name,
						table_name  : table_name}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := request{	function_id: func_id_restore_snapshot, 
					response_channel: &resp_chan,
					data : &data} 
	req.response_channel = &resp_chan

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

func (ds * datastore ) WriteSnapshot(	snapshot_name string, 
										schema_name string, 
										table_name string,
										path_file string) error {
	// Create a channel for the response
	resp_chan := make( chan(error) )

	// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	snapshot_name : snapshot_name,
						schema_name : schema_name,
						table_name  : table_name,
						path_file : path_file}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := request{	function_id: func_id_write_snapshot, 
					response_channel: &resp_chan,
					data : &data} 
	req.response_channel = &resp_chan

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


func (ds * datastore ) ReadSnapshot(snapshot_name string, 
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
	data := persist { 	snapshot_name : snapshot_name,
						schema_name : schema_name,
						table_name  : table_name,
						path_file : path_file}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := request{	function_id: func_id_read_snapshot, 
					response_channel: &resp_chan,
					data : &data} 
	req.response_channel = &resp_chan

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

func (ds * datastore ) ClearSnapshot(	snapshot_name string, 
										schema_name string, 
										table_name string) error {

	// -------------------------------------------------------------------------
		// Create a channel for the response
	// -------------------------------------------------------------------------
	resp_chan := make( chan(error) )

	// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	snapshot_name : snapshot_name,
						schema_name : schema_name,
						table_name  : table_name}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := request{	function_id: func_id_clear_snapshot, 
					response_channel: &resp_chan,
					data : &data} 
	req.response_channel = &resp_chan

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

func (ds * datastore ) ClearRealtime(	schema_name string, 
										table_name string) error {
	// -------------------------------------------------------------------------
	// Create a channel for the response
	// -------------------------------------------------------------------------
	resp_chan := make( chan(error) )

// -------------------------------------------------------------------------
	// allocate a persist structire
	// -------------------------------------------------------------------------
	data := persist { 	schema_name : schema_name,
						table_name  : table_name}

	// -------------------------------------------------------------------------
	// Create a request
	// -------------------------------------------------------------------------
	req := request{	function_id: func_id_clear_realtime, 
					response_channel: &resp_chan,
					data : &data} 
	req.response_channel = &resp_chan

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