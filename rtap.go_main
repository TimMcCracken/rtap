package main

import (
	"fmt"
	"gorm.io/gorm"
  	"github.com/glebarez/sqlite"
//	"io/fs"
	"log"
	"os"
//	"path/filepath"
	"rtap/realm"
//	"time"
)



// we need a realm object. Or do we want to call it the domain?  This object will
// store things that must be common across all data stores. This includes maps 
// for snapshots and contexts.

/* ---------------------------------------------------------------------------
SYSTEM is the root object of the system It is a singleton within any GO
program. It is fiarly small supporting the following functions.

CreateRealm()   	[only used during a 'SysGen' operation]
DropRealm()			[rarely used]
UseRealm() 			[called one time at system startup to initialize the RTDSMS]
CreateDataStore()	[typically only used during 'SysGen']
DropDataStore()		[should be rarely used]
GetDataStore()		Called anytime a program needs an object tor the APIs to
					operate on.
Start()				starts all the datastores
Stop()				stops all the datastores

---------------------------------------------------------------------------- */

type rtap struct {
	realms				[]*realm.realm
	realms_map 			map[string]int
}



// the singleton SYSTEM
var System rtap




// -----------------------------------------------------------------------------
//Function to initialize the "System" structure at statup. It functions by 
// looping through all the folders within the "base" folder that is declared
// as the environment variable.
// -----------------------------------------------------------------------------
func init(){

	System.realms_map = make(map[string]int)
	base := os.Getenv("rtdsms_base")

	// Open the directory
	d, err := os.Open(base)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	// Read the directory contents
	files, err := d.Readdir(-1) // -1 means no limit on number of files
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the files in the directory
	for _, file := range files {
		// Check if it's a directory and not the current or parent directory
		if file.IsDir() && file.Name() != "." && file.Name() != ".." {
			// Print the name of the subdirectory
			fmt.Println(file.Name())

			// create a realm in 'System'
			realm := realm.realm{	
				path: base + "/" + file.Name(),
				filename: base + "/" + file.Name() + "/" + file.Name() + ".realm",
				} 
			realm.descriptor.realm_name = file.Name()

			System.realms_map[file.Name()]= len(System.realms_map)
			System.realms = append(System.realms, & realm)

		}
	}
}

/*
descriptor 		realm_descriptor
	filename		string
	path			string
	db				*gorm.DB
	domains			[]*domain
	domains_map 	map[string]int
	requests 		chan( * request)
*/



// -----------------------------------------------------------------------------
// ConstructRealms() construct either all reamls unconditionally, or the realms
// that are active.
// -----------------------------------------------------------------------------


// -----------------------------------------------------------------------------
// ConstructRealm() construct a single named realm. 'Construct' means to 
// allocate space for all the datastores in the realm.
// -----------------------------------------------------------------------------

/*
func (system * rtdms ) ConstructRealm(realm_name string) error {

		// Validate the parameters
		err := ValidateObjectName(realm_name)
		if err != nil {
			return err
		}
	
		// ---------------------------------------------------------------------
		// Loop through all the domains
		// ---------------------------------------------------------------------



		return nil
}
		*/







/*******************************************************************************
Data Definition Language
*******************************************************************************/


/* ---------------------------------------------------------------------------
CreateRealm()
---------------------------------------------------------------------------- */
func (system * rtdms ) CreateRealm(realm_name string) error {
	
	// Validate the parameters
	err := ValidateObjectName(realm_name)
	if err != nil {
		return err
	}
	
	base := os.Getenv("rtdsms_base")
	path := base + "/" + realm_name
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	filename := path + "/" + realm_name + ".realm"

	// Create the sqlite '.realm' file
	db, err := gorm.Open(sqlite.Open(filename + "?_pragma=foreign_keys(1)" ), &gorm.Config{})
	if err != nil {
	  return fmt.Errorf("failed to connect database using path [%s]", filename)
	}
	sqlDB, err := db.DB()
    if err != nil {
		return fmt.Errorf("Internal error: failed to obtain database.DB in CreateRealm()")
    }
	defer sqlDB.Close()
	

	for _, sql := range sql_create_schema {
		db.Exec(sql)

		// the below will not allow the file to be deleted
		if db.Error != nil {
			_, err := db.DB()
			if(err != nil) {
				return fmt.Errorf("error fetching db.DB %s", err)
			}
			os.Remove(filename)
			return db.Error
	  	}
	}

	// -------------------------------------------------------------------------
	// make() the map, if it hasnt been already.
	// -------------------------------------------------------------------------
	if System.realms_map == nil {
		System.realms_map = make(map[string]int)
	}

	// -------------------------------------------------------------------------
	// Initialize the values in the realm object
	// -------------------------------------------------------------------------
	realm := realm{	
				path: base + "/" + realm_name,
				filename: base + "/" + realm_name + "/" + realm_name + ".realm",
				} 
				realm.descriptor.realm_name = realm_name
	System.realms_map[realm_name]= len(System.realms_map)
	System.realms = append(System.realms, & realm)

	return nil
}


/* ---------------------------------------------------------------------------
DropRealm()
---------------------------------------------------------------------------- */
func (system * rtdms ) DropRealm(realm_name string) error {
	return nil
}



/* ---------------------------------------------------------------------------
GetRealm()
---------------------------------------------------------------------------- */
func (system * rtdms) GetRealm(realm_name string) ( * realm , error) {

	// Validate the parameters
	err := ValidateObjectName(realm_name)
	if err != nil {
		return nil, err
	}

	i, ok := system.realms_map[realm_name]

	if !ok {
		return nil, fmt.Errorf("Realm [%s] does not exist.")
	} 

	return system.realms[i], nil

/*
	path := os.Getenv("rtdsms_base") + "/" + realm_name

	system.activeRealm = realm_name
	system.activeRealmPath  = os.Getenv("rtdsms_base") + "/" + realm_name + "/"
	system.activeRealmFile = path + "/" + realm_name + ".realm"
*/	

}






















func (system * rtdms) PrintSchema() {
/*
	fmt.Printf("\nRealm: %s\n", system.activeRealm)

	for i, ds := range system.datastores {
		fmt.Printf("  Datastore #: %v Name: %v\n", i, ds.name)
		ds.PrintSchema()
	}
*/
}



