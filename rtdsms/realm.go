package rtdsms

import (
	"fmt"
	"gorm.io/gorm"
  	"github.com/glebarez/sqlite"
	"os"
	"time"
)


var sql_create_schema []string  

/*
`CREATE TABLE IF NOT EXISTS domain_states (
	state_id				integer primary key,
	state_name		text 	unique
	) STRICT`,

	`INSERT OR REPLACE INTO domain_states values (0, 'Offline')`,
	`INSERT OR REPLACE INTO domain_states values (1, 'Starting')`,
	`INSERT OR REPLACE INTO domain_states values (2, 'Online')`,
	`INSERT OR REPLACE INTO domain_states values (3, 'Stopping')`,
	`INSERT OR REPLACE INTO domain_states values (4, 'Paused')`,
*/

func init(){
	sql_create_schema = []string{


		`CREATE TABLE domains (
			domain_id		integer primary key,
			domain_name		text not null unique,
			created_at		integer not null,
			updated_at		integer	
			) STRICT`,
	
		`CREATE TABLE datastores (
			domain_id		integer references domains,
			datastore_id	integer not null,
			datastore_name	text not null,
			created_at		integer not null,
			updated_at		integer,	
			primary key 	(domain_id, datastore_id),
			UNIQUE			(domain_id, datastore_name)
			) STRICT`,		

		`CREATE TABLE snapshots(
			domain_id		integer references domains,
			datastore_id	integer not null,
			snapshot_id		integer not null,
			snapshot_name	text not null,
			created_at		integer not null,
			updated_at		integer,	
			foreign key 	(domain_id, datastore_id) references datastores,
			primary key 	(domain_id, datastore_id, snapshot_id),
			UNIQUE			(domain_id, datastore_id, snapshot_name)
			) STRICT`,
	

		`CREATE TABLE schemas(
			domain_id		integer,	
			datastore_id	integer,
			schema_id		integer not null,
			schema_name		text not null,
			created_at		integer not null,
			updated_at		integer,	
			foreign key 	(domain_id, datastore_id) references datastores,
			primary key 	(domain_id, datastore_id, schema_id),
			UNIQUE			(domain_id, datastore_id, schema_name)
			) STRICT`,
	
		`CREATE TABLE tables(
			domain_id		integer,	
			datastore_id	integer,
			schema_id		integer,
			table_id		integer  not null,
			table_name		text  not null,
			total_rows		integer not null,
			created_at		integer  not null,
			updated_at		integer,	
			foreign key 	(domain_id, datastore_id, schema_id) references schemas,
			primary key 	(domain_id, datastore_id, schema_id, table_id),
			UNIQUE			(domain_id, datastore_id, schema_id, table_name)
			) STRICT`,

		`CREATE TABLE columns (
			domain_id		integer,	
			datastore_id	integer,
			schema_id		integer,
			table_id		integer,
			column_id		integer  not null,
			column_name		text 	not null,
			type_name		text	not null,
			not_null		integer	not null,
			unique_			integer	not null,
			indexed			integer	not null,
			no_persist		integer	not null,
			default_value	text,
			check_			text,
			min_integer		integer,
			max_integer		integer,
			min_uint		integer,
			max_uint		integer,
			min_float		real,
			max_float		real,
			min_length		integer, 
			length			integer,
			created_at		integer,
			updated_at		integer	,
			foreign key 	(domain_id, datastore_id, schema_id, table_id) references tables,
			primary key 	(domain_id, datastore_id, schema_id, table_id, column_id),
			UNIQUE			(domain_id, datastore_id, schema_id, table_id, column_name)
			) STRICT`,

		`CREATE TABLE bitmap_names (
			domain_id		integer,	
			datastore_id	integer,
			schema_id		integer,
			table_id		integer,
			column_id		integer,
			bitfield_id		integer,
			bitfield_name	text 	not null,
			createa_at		integer,
			updated_at		integer,	
			foreign key (domain_id, datastore_id, schema_id, table_id, column_id) references columns,
			primary key (domain_id, datastore_id, schema_id, table_id, column_id, bitfield_id)
			) STRICT`,
				
	}
}

type realm struct {
	descriptor 				realm_descriptor
	filename				string
	path					string
	db						*gorm.DB
	domains					[]*domain
	domains_map 			map[string]int
	requests 				chan( * request)
}


func (realm * realm) GetDatastore(domain_name string, datastore_name string) (*datastore, error)  {

	domain_index, ok := realm.domains_map[domain_name]
	if ! ok {
		return nil, fmt.Errorf("Domain [%s] does not exist in realm [%s].", domain_name, realm.descriptor.realm_name)
	}

	datastore_index, ok := realm.domains[domain_index].datastores_map[datastore_name]
	if ! ok {
		return nil, fmt.Errorf("Domain [%s] does not exist in realm [%s].", domain_name, realm.descriptor.realm_name)
	}

	return realm.domains[domain_index].datastores[datastore_index], nil
}


// -------------------------------------------------------------------------
// 'Construct' the database structures. This allocates all the space needed
// and 'makes" anything that needs it. It does NOT save or load any data,
// rather the data is created with zero values. The sqlite db for the 
// specified realm is used to do the 'Construct' function. 
//
// realm.Construct() calls domain.construct() which calls 
// datastore.construct() which calls schema.construct() which calls 
// table.construct(). 
// -------------------------------------------------------------------------
func (realm * realm) Construct() error  {

	if realm.domains_map == nil {
		realm.domains_map = make(map[string]int)
	}

//	fmt.Printf("Realm name: [%s]\n", realm.descriptor.realm_name)

	// Open the sqlite file
	db, err := gorm.Open(sqlite.Open(realm.filename + "?_pragma=foreign_keys(1)" ), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database using path [%s]", realm.filename)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Internal error: failed to obtain database.DB in CreateDomain()")
	}
	defer sqlDB.Close()
		
	rows, err := sqlDB.Query("SELECT domain_name FROM domains order by domain_id")
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
			dom := domain {
			}
			dom.descriptor.domain_name = name
			dom.descriptor.realm_name = realm.descriptor.realm_name
			realm.domains_map[name] = len(realm.domains_map)
			realm.domains = append(realm.domains, &dom)
		}
    }
    if err = rows.Err(); err != nil {
		return fmt.Errorf("%v\n", err)
    }

	for _, dom := range realm.domains {
		err = dom.construct(realm.filename)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}

	return nil
}



/* ---------------------------------------------------------------------------
CreateDomain
---------------------------------------------------------------------------- */
func (realm * realm) CreateDomain(domain_name string) error  {
	// Validate the parameters
	err := validateObjectName(domain_name) 
		if err != nil {
			return err
		}

	// ------------------------------------------------------------------------
	// Add it to the sqlite database
	// ------------------------------------------------------------------------
	// Open the sqlite file
	db, err := gorm.Open(sqlite.Open(realm.filename + "?_pragma=foreign_keys(1)" ), &gorm.Config{})
	if err != nil {
	  return fmt.Errorf("failed to connect database using path [%s]", realm.filename)
	}
	sqlDB, err := db.DB()
    if err != nil {
		return fmt.Errorf("Internal error: failed to obtain database.DB in CreateDomain()")
    }
	defer sqlDB.Close()

	// Get the highest primary key in the data stores table
	var domain_id int64
	result := db.Raw("SELECT count(*) FROM domains").Scan(&domain_id)
	if result.Error != nil {
		return result.Error
	}
	
	if domain_id == 0 {
		domain_id = 1
	} else {
		result = db.Raw("SELECT max(domain_id) FROM domains").Scan(&domain_id)
		if result.Error != nil {
			return result.Error
		}
		domain_id ++;
	}

	// add it to the sqlite3 database
	now := time.Now().UnixMilli()
	db.Exec("INSERT OR REPLACE INTO domains values (?, ?, ?, 0)", domain_id, domain_name, now)
	if db.Error != nil {
		return db.Error
	}

	// build the snapshot descriptors to use when creating the datastores
	sd1, err := NewSnapshotDescriptor("diskcopy")
	sd2, err := NewSnapshotDescriptor("standby")
	sd_arr := []*SnapshotDescriptor{sd1, sd2}

	// Create the two standard datastores
	err = realm.CreateDatastore(domain_name, "sys", sd_arr, nil)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	err = realm.CreateDatastore(domain_name, "system", sd_arr, nil)
	if err != nil {
		return fmt.Errorf("Error:%s", err)
	}

	return nil
}


/* ---------------------------------------------------------------------------
DropDomain
---------------------------------------------------------------------------- */
// TODO



/* ---------------------------------------------------------------------------
CreateDataStore
Note: Context Descriptors are for future use 
---------------------------------------------------------------------------- */
func (realm * realm) CreateDatastore(	domain_name string, 
										datastore_name string, 
										snapshot_descriptors []*SnapshotDescriptor , 
										context_descriptors []*ContextDescriptor ) error {

	// check that a valid realm has been set in the system structure by
	// trying to open the database file
	db, err := gorm.Open(sqlite.Open(realm.filename + "?_pragma=foreign_keys(1)" ), &gorm.Config{})
	if err != nil {
	  return fmt.Errorf("failed to connect database using path [%s]", realm.filename)
	}
	sqlDB, err := db.DB()
    if err != nil {
		return fmt.Errorf("Internal error: failed to obtain database.DB in CreateDatastore()")
    }
	defer sqlDB.Close()

	// Validate the domain_name
	err = validateObjectName(domain_name) 
	if err != nil {
		return err
	}
	
	var domain_count int
	result := db.Raw("SELECT count(*) FROM domains where domain_name = ?", domain_name).Scan(&domain_count)
	if result.Error != nil {
		return result.Error
	}

	if domain_count == 0 {
		return fmt.Errorf("Domain [%s] does not exist.", domain_name)
	}
	

	// Validate the datastore_name
	err = validateObjectName(datastore_name) 
	if err != nil {
		return err
	}

	for _, sd := range snapshot_descriptors {
		err := validateObjectName(sd.name) 
		if err != nil {
			return err
		}
	}

	for _, cd := range context_descriptors {
		err := validateObjectName(cd.name) 
		if err != nil {
			return err
		}
	}


	// ------------------------------------------------------------------------
	// Get the domain_id
	// ------------------------------------------------------------------------
	var domain_id	int64
	
	result = db.Raw("SELECT domain_id FROM domains where domain_name = ?", domain_name).Scan(&domain_id)
	if result.Error != nil {
		return result.Error
	}


	// ------------------------------------------------------------------------
	// Add it to the sqlite database
	// ------------------------------------------------------------------------

	// Get the highest primary key in the data stores table
	var datastore_id int64
	db.Raw("SELECT count(*) FROM datastores").Scan(&datastore_id)
	if datastore_id == 0 {
		datastore_id = 1
	} else {
		db.Raw("SELECT max(datastore_id) FROM datastores").Scan(&datastore_id)
		datastore_id ++;
	}

	// add it to the sqlite3 database
	now := time.Now().UnixMilli()
	db.Exec("INSERT OR REPLACE INTO datastores values (?, ?, ?, ?, 0)", 
												domain_id, datastore_id, 
												datastore_name, now)
	if db.Error != nil {
		return db.Error
	}

					
	// ------------------------------------------------------------------------
	// Allocate space ( move "inside")
	// ------------------------------------------------------------------------
	
	/*
	ds := &datastore{	name: datastore_name, 
		requests : make(chan * request, default_channel_depth),
		tableMap : make(map[string]uint64),
		snapshotDescriptorsMap: make(map[string]uint64),
		contextDescriptorsMap: make(map[string]uint64),	}
	system.datastores[datastore_name] = ds
	*/

	// ------------------------------------------------------------------------
	// Create the disk directory where the data will be stored
	// ------------------------------------------------------------------------
	path := realm.path + "/" + domain_name + "/" + datastore_name  + "/"

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}


/*
	domain_id		integer,	
	datastore_id	integer,
	descriptor_id	integer not null,
	descriptor_name	text not null,
	created_at		integer not null,
	updated_at		integer,	
*/

	// ------------------------------------------------------------------------
	// Add the snapshot and context descriptors
	// ------------------------------------------------------------------------
	
	if snapshot_descriptors != nil {

		for i, sd := range snapshot_descriptors {
			
			db.Exec("INSERT INTO snapshots Values (?, ?, ?, ?, ?, 0)", 
					domain_id, datastore_id, i, sd.name, now)
			if db.Error != nil {
				return db.Error
			}

		}
	}

	/*
	if context_descriptors != nil {
		ds.contextDescriptors = contDescs
		
		for i, cd := range contDescs {
			ds.contextDescriptorsMap[cd.name] = uint64(i)
		}
	} */
	

	// By default, set the datastore state to 'online'
	//ds.state = DS_ONLINE

	return nil
}






/* ---------------------------------------------------------------------------
DropDataStore()
---------------------------------------------------------------------------- */
func (realm * realm) DropDataStore(domain_name string, datastore_name string) error {

	return nil
}




/* ---------------------------------------------------------------------------
CreateSchema
Note: Context Descriptors are for future use 
---------------------------------------------------------------------------- */
func (realm * realm) CreateSchema(	domain_name string, 
										datastore_name string, 
										schema_name string) error {


	// check that a valid realm has been set in the system structure by
	// trying to open the database file
	db, err := gorm.Open(sqlite.Open(realm.filename + "?_pragma=foreign_keys(1)" ), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database using path [%s]", realm.filename)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Internal error: failed to obtain database.DB in CreateSchema()")
	}
	defer sqlDB.Close()

	// ------------------------------------------------------------------------
	// Validate the domain name
	// ------------------------------------------------------------------------
	err = validateObjectName(domain_name) 
	if err != nil {
		return err
	}

	var domain_count int
	result := db.Raw("SELECT count(*) FROM domains where domain_name = ?", 
						domain_name).Scan(&domain_count)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}

	if domain_count == 0 {
		return fmt.Errorf("Domain [%s] does not exist.", domain_name)
	}

	// ------------------------------------------------------------------------
	// Get the domain id
	// ------------------------------------------------------------------------
	var domain_id		int64
	result = db.Raw("SELECT domain_id FROM domains where domain_name = ? ", 
					domain_name).Scan(&domain_id)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}


	// ------------------------------------------------------------------------
	// Validate the datastore_name
	// ------------------------------------------------------------------------
	err = validateObjectName(datastore_name) 
	if err != nil {
		return err
	}

	var datastore_count int
	result = db.Raw("SELECT count(*) FROM datastores where domain_id = ? AND " + 
					"datastore_name = ?",  domain_id, 
					datastore_name).Scan(&datastore_count)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}
	if datastore_count == 0 {
		return fmt.Errorf("Datastore [%s] does not exist.", datastore_name)
	}

	// ------------------------------------------------------------------------
	// Get the datastore id
	// ------------------------------------------------------------------------
	var datastore_id	int64
	result = db.Raw("SELECT datastore_id FROM datastores WHERE domain_id = ?" + 
					"AND datastore_name = ? ", domain_id, 
					datastore_name).Scan(&datastore_id)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}

	// ------------------------------------------------------------------------
	// Add it to the sqlite database
	// ------------------------------------------------------------------------

	// Get the highest primary key in the data stores table
	var schema_id int64
	result = db.Raw("SELECT count(*) FROM schemas").Scan(&schema_id)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}
	if schema_id == 0 {
		schema_id = 1
	} else {
		db.Raw("SELECT max(schema_id) FROM schemas").Scan(&schema_id)
		schema_id ++;
	}

	// ------------------------------------------------------------------------
	// add it to the sqlite3 database
	// ------------------------------------------------------------------------
	now := time.Now().UnixMilli()
	db.Exec("INSERT OR REPLACE INTO schemas values (?, ?, ?, ?, ?, 0)", 
			domain_id, datastore_id, schema_id, schema_name, now)
	if db.Error != nil {
		return db.Error
	}

	// ------------------------------------------------------------------------
	// Create the disk directory where the data will be stored
	// ------------------------------------------------------------------------
	path := realm.path + "/" + domain_name + "/" + 
				datastore_name  + "/" + schema_name + "/"

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	// By default, set the datastore state to 'online'
	//ds.state = DS_ONLINE

	return nil
}



// -----------------------------------------------------------------------------
// create_table() 
// -----------------------------------------------------------------------------
func (realm * realm) CreateTable(	domain_name string, datastore_name string, 
								schema_name string, tbl_desc *table_descriptor) error {

	// Get a pointer to a table descriptor type
//	tbl_desc, ok := request.data.(*tableDescriptor)
//	if ok == false {
//		request.err = fmt.Errorf("internal error - type assertion failed.") // random message
//		return
//	}


	//The columns were validated when added to the table descrptor

	// check that a valid realm has been set in the system structure by
	// trying to open the database file
	db, err := gorm.Open(sqlite.Open(realm.filename + 
							"?_pragma=foreign_keys(1)" ), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database using path [%s]", 
							realm.filename)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Internal error: failed to obtain database.DB" + 
							"in CreateSchema()")
	}
	defer sqlDB.Close()

	// ------------------------------------------------------------------------
	// Validate the domain name
	// ------------------------------------------------------------------------
	err = validateObjectName(domain_name) 
	if err != nil {
		return err
	}

	var domain_count int
	result := db.Raw("SELECT count(*) FROM domains where domain_name = ?", 
						domain_name).Scan(&domain_count)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}
	if domain_count == 0 {
		return fmt.Errorf("Domain [%s] does not exist.", domain_name)
	}

	// ------------------------------------------------------------------------
	// Get the domain id
	// ------------------------------------------------------------------------
	var domain_id		int64
	result = db.Raw("SELECT domain_id FROM domains where domain_name = ? ", 
					domain_name).Scan(&domain_id)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}

	// ------------------------------------------------------------------------
	// Validate the datastore_name
	// ------------------------------------------------------------------------
	err = validateObjectName(datastore_name) 
	if err != nil {
		return err
	}

	var datastore_count int
	result = db.Raw("SELECT count(*) FROM datastores where domain_id = ? " + 
					"AND datastore_name = ?", domain_id, 
					datastore_name).Scan(&datastore_count)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}
	if datastore_count == 0 {
		return fmt.Errorf("Datastore [%s] does not exist.", datastore_name)
	}

	// ------------------------------------------------------------------------
	// Get the datastore id
	// ------------------------------------------------------------------------
	var datastore_id	int64
	db.Raw("SELECT datastore_id FROM datastores where domain_id = ? AND " + 
			"datastore_name = ? ", domain_id, 
			datastore_name).Scan(&datastore_id)

	// ------------------------------------------------------------------------
	// Validate the schema_name
	// ------------------------------------------------------------------------
	err = validateObjectName(schema_name) 
	if err != nil {
		return err
	}

	var schema_count int
	result = db.Raw("SELECT count(*) FROM schemas where domain_id = ? AND " +
					"datastore_id = ? AND schema_name = ?", domain_id, 
					datastore_id, schema_name).Scan(&schema_count)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}
	if schema_count == 0 {
		return fmt.Errorf("schema [%s] does not exist.", schema_name)
	}

	// ------------------------------------------------------------------------
	// Get the schema id
	// ------------------------------------------------------------------------
	var schema_id	int64
	result = db.Raw("SELECT schema_id FROM schemas where domain_id = ? AND " + 
					"datastore_id  = ? AND schema_name = ?", 
					domain_id, datastore_id, schema_name).Scan(&schema_id)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}

	// ------------------------------------------------------------------------
	// Add it to the sqlite database
	// ------------------------------------------------------------------------

	// Get the highest primary key in the 'tables' table 
	var tables_count int64
	result = db.Raw("SELECT count(*) FROM tables WHERE domain_id = ? AND " + 
					" datastore_id = ? AND schema_id =? ", 
					domain_id, datastore_id, schema_id).Scan(&tables_count)
	if result.Error != nil {
		return fmt.Errorf("SQL error: [?]", result.Error)
	}
	
	var table_id int64
	if tables_count == 0 {
		table_id = 1
	} else {
		result = db.Raw("SELECT max(table_id) FROM tables WHERE domain_id = ?" +
						" AND datastore_id = ? AND schema_id =? ",  
						domain_id, datastore_id, schema_id ).Scan(&table_id)
		if result.Error != nil {
			return result.Error
		}					
		table_id ++;
	}

	// start a transaction
	db.Begin()

	// add the table to the sqlite3 database
	now := time.Now().UnixMilli()
	db.Exec("INSERT INTO tables values (?, ?, ?, ?, ?, ?, ?, 0)", 
								domain_id, datastore_id, schema_id, 
								table_id, tbl_desc.table_name, 
								tbl_desc.total_rows, now, 0)
	if db.Error != nil {
		return fmt.Errorf("Error: %v\n", db.Error)
  	}

	// add the columns to the sqlite3 database
	for i, col_desc := range tbl_desc.column_descriptors {

		db.Exec("INSERT INTO columns values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?," +
				" ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)", 
		domain_id, datastore_id, schema_id, table_id, i, col_desc.name, 
		col_desc.type_name, col_desc.not_null, col_desc.unique, 
		col_desc.indexed, col_desc.no_persist, col_desc.default_value, 
		col_desc.check, col_desc.min_int, col_desc.max_int, 
		int64(col_desc.min_uint), int64(col_desc.max_uint), 
		col_desc.min_float, col_desc.max_float,  col_desc.min_length, 
		col_desc.length,  now ) 
		
		if db.Error != nil {
			return fmt.Errorf("Error: %v\n", db.Error) 
		}
	}

	// end transaction
	db.Commit()

	// ------------------------------------------------------------------------
	//  Create the table in the data store.
	// TODO: Make this if ds_state is 'online'
	// ------------------------------------------------------------------------
	/*
	
	tbl := new(table)

	tbl.init(tbl_desc) // most of the real work takes place in init.
	tbl.path = ds.path + tbl.descriptor.name + "/"

	ds.tableMap[tbl.descriptor.name] = uint64(len(ds.tables))
	ds.tables = append(ds.tables, tbl) 
	
	request.err = nil
	*/

	return nil
}
