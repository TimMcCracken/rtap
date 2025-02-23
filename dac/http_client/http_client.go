


package http_client


import (

	"fmt"
	"rtap/rtdsms"
)


func init() {
	
}



func Build(realm_name string) error {

	var system = &rtdsms.System
	realm, err := system.GetRealm(realm_name)


	snapshot_descriptor1, err := rtdsms.NewSnapshotDescriptor("diskcopy")
	snapshot_descriptor2, err := rtdsms.NewSnapshotDescriptor("standby")
	snapshot_descriptor_arr := []*rtdsms.SnapshotDescriptor{snapshot_descriptor1, snapshot_descriptor2}

	err = realm.CreateDatastore("prod", "http_client", snapshot_descriptor_arr, nil)
	if err != nil {
		return fmt.Errorf("Error creating http_client datastore: %s", err)
	}

	err = realm.CreateSchema("prod", "http_client", "config")
	if err != nil {
		return fmt.Errorf("Error creating http_client schema: %s", err)
	}
	
	
	//  Build stuff to use to create tables
	
	cd1, err := rtdsms.NewColumnDescriptor("script_name", "varUTF8", 0 /*rtdsms.INDEXED*/, uint64(16))
	if err != nil {
		return fmt.Errorf("Error NewColumnDescriptor http_client column: [%s]", err)
	}
	cd2, err := rtdsms.NewColumnDescriptor("frequency", "int64", 0)
	if err != nil {
		return fmt.Errorf("Error NewColumnDescriptor http_client column: [%s]", err)
	}
	cd3, err := rtdsms.NewColumnDescriptor("offset", "int64", 0 )
	if err != nil {
		return fmt.Errorf("Error NewColumnDescriptor http_client column: [%s]", err)
	}


	cd_arr := []*rtdsms.ColumnDescriptor{cd1, cd2, cd3}
	
	td, err := rtdsms.NewTableDescriptor("dev_scada", "prod", "http_client", "config", "polling", 100, cd_arr)
	if err != nil {
		return fmt.Errorf("Error creating table descriptor: %v\n", err)
	}


	
	// -------------------------------------------------------------------------
	// Create a table
	// -------------------------------------------------------------------------

	err = realm.CreateTable("prod", "http_client", "config", td)
	if err != nil {
		return fmt.Errorf("Error creating table: %v\n", err)
	} 

	return nil
}