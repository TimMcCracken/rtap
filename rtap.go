package main


import (
	"fmt"
//	"math"
	"os"
//	"reflect"
//	"regexp"
//	"runtime"
	"time"

	mq "rtap/message_q"
	"rtap/metronome"
	"rtap/rtdsms"
	"rtap/dac/http_client"
)




func main() {

	mq.Start()

	mq.Send("ToHellAndBack", nil)

	metronome.Metronome()

	os.Exit(0)



	// Declare the variable for the RTDSMS. This is a 'global' variable
	// and typically there is only one in the system, although in very
	// special cases there might be more than one. 
	var system = &rtdsms.System
	
	err := system.CreateRealm("devscada")
	if err != nil {
		fmt.Printf("Error creating realm: %s", err)
	}
	
	devscada, err := system.GetRealm("devscada")
	if err != nil {
		fmt.Printf("Error getting realm: %s", err)
	}


	err = devscada.CreateDomain("prod")
	if err != nil {
		fmt.Printf("Error creating domain: %s", err)
	}



	sd1, err := rtdsms.NewSnapshotDescriptor("diskcopy")
	sd2, err := rtdsms.NewSnapshotDescriptor("standby")
	sd_arr := []*rtdsms.SnapshotDescriptor{sd1, sd2}

	err = devscada.CreateDatastore("prod", "scada", sd_arr, nil)
	if err != nil {
		fmt.Printf("Error creating datastore: %s", err)
	}



	err = devscada.CreateSchema("prod", "scada", "config")
	if err != nil {
		fmt.Printf("Error creating schema: %s", err)
	}
	
	
	//  Build stuff to use to create tables
	
	cd1, err := rtdsms.NewColumnDescriptor("column1", "varUTF8", 0 /*rtdsms.INDEXED*/, uint64(16))
	if err != nil {
		fmt.Printf("ERRORRRRRRRRRRR: %v\n", err)
	}
	cd2, err := rtdsms.NewColumnDescriptor("column2", "int64", 0)
	cd3, err := rtdsms.NewColumnDescriptor("column3", "float64", 0 )

	arr := []*rtdsms.ColumnDescriptor{cd1, cd2, cd3}
	
	td, err := rtdsms.NewTableDescriptor("dev_scada", "prod", "scada", "config", "my_table", 1000000 , arr)
	if err != nil {
		fmt.Printf("Error creating table descriptor: %v\n", err)
	}


	
	// -------------------------------------------------------------------------
	// Create a table
	// -------------------------------------------------------------------------

	err = devscada.CreateTable("prod", "scada", "config", td)
	if err != nil {
		fmt.Printf("Error creating table: %v\n", err)
	} 
	
	


	// ------------------------------------------------
	// http_client
	// ------------------------------------------------

	err = http_client.Build("devscada")
	if err != nil {
		fmt.Printf("Error building http_client: %v\n", err)
	} 












	// -------------------------------------------------------------------------
	// 'Construct' the database structures. This allocates all the space needed
	// and 'makes" anything that needs it. It does NOT save or load any data,
	// rather the data is created with zero values. The sqlite db for the 
	// specified realm is used to do the 'Construct' function
	// -------------------------------------------------------------------------
	err = devscada.Construct()
	if err != nil {
		fmt.Printf("Error constructing: %v\n", err)
	} 

	prod, err := devscada.GetDatastore("prod", "scada" )
	if err != nil {
		fmt.Printf("Error getting datastore: %s\n", err)
	} 
	
	prod.Start()
	time.Sleep(1 * time.Second)

	// Build stuff to insert
	fields := []string{"row_id", "column1", "column2", "column3"}

	// For an insert, the 'lookup_column_key' and 'lookup_row_key' will be ignored
	// as the first variant will be used as the 'row_id'
	records := []*rtdsms.Record{

		rtdsms.NewRecord("config", "my_table", "", "", fields, 0 ),
	//	rtdsms.NewRecord("config", "my_table", "", "", fields, 0 ),
	//	rtdsms.NewRecord("config", "my_table", "", "", fields, 0 ),
	}


	records[0].SetValue(0, "row_0")
	records[0].SetValue(1, "hello1")
	records[0].SetValue(2, int64(1234))
	records[0].SetValue(3, 1234.5678)
/*
	records[1].SetValue(0, "row2")
	records[1].SetValue(1, "hello2")
	records[1].SetValue(2, int64(234))
	records[1].SetValue(3, 234.5678)

	records[2].SetValue(0, "row_3")
	records[2].SetValue(1, "hello3")
	records[2].SetValue(2, int64(34))
	records[2].SetValue(3, 34.5678)
*/
	// ----------------------------
	// Insert the records
	// ----------------------------
	start := time.Now()

	for i:=0; i < 100000; i++ {
		records[0].SetValue(0, fmt.Sprintf("row_%d", i))
	
		err = prod.Insert(records)
		if err != nil {
			fmt.Printf("Error prod Inserting: %v\n", err)
		}
	} 

	elapsed := time.Since(start)
	fmt.Printf("Insert Elapsed time: %v\n", elapsed)


	// Insert the recrods
//	err = prod.Insert(records)
//	if err != nil {
//		fmt.Printf("Error prod Inserting: %v\n", err)
//	} 


	// read back the values
	records2 := []*rtdsms.Record{
		rtdsms.NewRecord("config", "my_table", "row_id", "row_751", fields, 0 ),
	//	rtdsms.NewRecord("config", "my_table", "column1", "hello2", fields, 0 ),
	//	rtdsms.NewRecord("config", "my_table", "column1", "hello1", fields, 0 ),
	}
	


	start = time.Now()

	for i:=0; i < 1000000; i++ {
		err = prod.Select(records2)
		if err != nil {
			fmt.Printf("Error selecting: %v\n", err)
		}
	}

	elapsed = time.Since(start)
	fmt.Printf("Select Elapsed time: %v\n", elapsed)



	err, v1 := records2[0].GetValue(0)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v2 := records2[0].GetValue(1)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v3 := records2[0].GetValue(2)
	if err != nil {
		fmt.Printf("%v",err)
	}
	
	err, v4 := records2[0].GetValue(3)
	if err != nil {
		fmt.Printf("%v",err)
	}
	fmt.Printf("Values before update:  %v  %v %v %v\n", v1, v2, v3, v4)


	// --------------------------------------------
	// update the values
	// --------------------------------------------

	fields3 := []string{"column1", "column2", "column3"}

	records3 := []*rtdsms.Record{
		rtdsms.NewRecord("config", "my_table", "row_id", "row_751", fields3, 0 ),
	//	rtdsms.NewRecord("config", "my_table", "column1", "hello2", fields, 0 ),
	//	rtdsms.NewRecord("config", "my_table", "column1", "hello1", fields, 0 ),
	}

	records3[0].SetValue(0, "world")
	records3[0].SetValue(1, int64(9876))
	records3[0].SetValue(2, 9876.5432)


	start = time.Now()

	for i:=0; i < 1000000; i++ {
		err = prod.Update(records3)
		if err != nil {
			fmt.Printf("Error selecting: %v\n", err)
		}
	}

	elapsed = time.Since(start)
	fmt.Printf("Update Elapsed time: %v\n", elapsed)



	//----------------------------------------
	// do another select to check the update
	//----------------------------------------
	// read back the values
	records4 := []*rtdsms.Record{
			rtdsms.NewRecord("config", "my_table", "row_id", "row_751", fields, 0 ),
		//	rtdsms.NewRecord("config", "my_table", "column1", "hello2", fields, 0 ),
		//	rtdsms.NewRecord("config", "my_table", "column1", "hello1", fields, 0 ),
	}

	err = prod.Select(records4)
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}

	
	err, v1 = records4[0].GetValue(0)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v2 = records4[0].GetValue(1)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v3 = records4[0].GetValue(2)
	if err != nil {
		fmt.Printf("%v",err)
	}
	
	err, v4 = records4[0].GetValue(3)
	if err != nil {
		fmt.Printf("%v",err)
	}
	fmt.Printf("Values after update:  %v  %v %v %v\n", v1, v2, v3, v4)




	// -----------------------------------
	// Take snapshot
	// -----------------------------------
	start = time.Now()
	err = prod.TakeSnapshot("standby", "config", "my_table" )
	if err != nil {
		fmt.Printf("Snapshot take Error:  %v\n",err)
	}

	elapsed = time.Since(start)
	fmt.Printf("Take Snapshot Elapsed time: %v\n", elapsed)


	// -----------------------------------
	// write snapshot
	// -----------------------------------
	start = time.Now()
	err = prod.WriteSnapshot("standby", "config", "my_table", "c:/rtap/test.snap")
	if err != nil {
		fmt.Printf("Snapshot write Error:  %v\n",err)
	}
	elapsed = time.Since(start)
	fmt.Printf("Write Snapshot Elapsed time: %v\n", elapsed)


	// We need a time delay because write snapshot is a goroutine
	time.Sleep(1 * time.Second)

	// -----------------------------------
	// Clear Realtime
	// -----------------------------------
	start = time.Now()
	err = prod.ClearRealtime("config", "my_table")
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}
	elapsed = time.Since(start)
	fmt.Printf("clear realtime Elapsed time: %v\n", elapsed)

	err = prod.Select(records4)
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}


	err, v1 = records4[0].GetValue(0)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v2 = records4[0].GetValue(1)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v3 = records4[0].GetValue(2)
	if err != nil {
		fmt.Printf("%v",err)
	}
	
	err, v4 = records4[0].GetValue(3)
	if err != nil {
		fmt.Printf("%v",err)
	}
	fmt.Printf("Values after clear realtime:  %v  %v %v %v\n", v1, v2, v3, v4)


	// -----------------------------------
	// Restore snapshot before clearing
	// -----------------------------------
	start = time.Now()
	err = prod.RestoreSnapshot("standby", "config", "my_table" )
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}
	elapsed = time.Since(start)
	fmt.Printf("restore snapshot Elapsed time: %v\n", elapsed)

	err = prod.Select(records4)
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}

	
	err, v1 = records4[0].GetValue(0)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v2 = records4[0].GetValue(1)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v3 = records4[0].GetValue(2)
	if err != nil {
		fmt.Printf("%v",err)
	}
	
	err, v4 = records4[0].GetValue(3)
	if err != nil {
		fmt.Printf("%v",err)
	}
	fmt.Printf("Values after restore before clearing snapshot:  %v  %v %v %v\n", v1, v2, v3, v4)




	// -----------------------------------
	// Clear Snapshots
	// -----------------------------------
	start = time.Now()
	err = prod.ClearSnapshot("standby", "config", "my_table")
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}
	elapsed = time.Since(start)
	fmt.Printf("clear realtime Elapsed time: %v\n", elapsed)


//	err = prod.ClearSnapshot("standby", "config", "my_table")
//	if err != nil {
//		fmt.Printf("Error selecting: %v\n", err)
//	}


//	err = prod.Select(records4)
//	if err != nil {
//		fmt.Printf("Error selecting: %v\n", err)
//	}



	// -----------------------------------
	// Restore snapshot after clearing
	// -----------------------------------
	start = time.Now()
	err = prod.RestoreSnapshot("standby", "config", "my_table" )
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}
	elapsed = time.Since(start)
	fmt.Printf("restore snapshot Elapsed time: %v\n", elapsed)

	err = prod.Select(records4)
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}

	
	err, v1 = records4[0].GetValue(0)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v2 = records4[0].GetValue(1)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v3 = records4[0].GetValue(2)
	if err != nil {
		fmt.Printf("%v",err)
	}
	
	err, v4 = records4[0].GetValue(3)
	if err != nil {
		fmt.Printf("%v",err)
	}
	fmt.Printf("Values after restore CLEARED snapshot:  %v  %v %v %v\n", v1, v2, v3, v4)






	// -----------------------------------
	// read snapshot
	// -----------------------------------
	start = time.Now()
	err = prod.ReadSnapshot("standby", "config", "my_table", "c:/rtap/test.snap")
	if err != nil {
		fmt.Printf("Snapshot read Error:  %v\n",err)
	}
	elapsed = time.Since(start)
	fmt.Printf("Read Snapshot Elapsed time: %v\n", elapsed)


	// -----------------------------------
	// Restore snapshot
	// -----------------------------------
	start = time.Now()
	err = prod.RestoreSnapshot("standby", "config", "my_table" )
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}
	elapsed = time.Since(start)
	fmt.Printf("restore snapshot Elapsed time: %v\n", elapsed)

	err = prod.Select(records4)
	if err != nil {
		fmt.Printf("Error selecting: %v\n", err)
	}

	
	err, v1 = records4[0].GetValue(0)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v2 = records4[0].GetValue(1)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v3 = records4[0].GetValue(2)
	if err != nil {
		fmt.Printf("%v",err)
	}
	
	err, v4 = records4[0].GetValue(3)
	if err != nil {
		fmt.Printf("%v",err)
	}
	fmt.Printf("Values after restore snapshot:  %v  %v %v %v\n", v1, v2, v3, v4)




	os.Exit(0)

	


/*


	
	tuples3[0].SetValue(0, "world")
	tuples3[0].SetValue(1, int64(9876))
	tuples3[0].SetValue(2, 9876.5432)

	err = sys.Update(tuples3)
	if err != nil {
		fmt.Printf("Error update: %v\n", err)
	}


	err = sys.Select(tuples2)
	if err != nil {
		fmt.Printf("Error selecting after updating: %v\n", err)
	}

	err, v1 = tuples2[0].GetValue(0)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v2 = tuples2[0].GetValue(1)
	if err != nil {
		fmt.Printf("%v",err)
	}
	
	err, v3 = tuples2[0].GetValue(2)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v4 = tuples2[0].GetValue(3)
	if err != nil {
		fmt.Printf("%v",err)
	}


	fmt.Printf("Values after update:  %v %v %v %v\n", v1, v2, v3, v4)


	err = sys.TakeSnapshot("diskcopy")
	if err != nil {
		fmt.Printf("Snapshot take Error:  %v\n",err)
	}
*/


/*
	err = sys.TakeSnapshot("standby")
	if err != nil {
		fmt.Printf("Snapshot take Error:  %v\n",err)
	}

	err = sys.RestoreSnapshot("diskcopy")
	if err != nil {
		fmt.Printf("Snapshot restore Error:  %v\n",err)
	}


	err = sys.RestoreSnapshot("standby")
	if err != nil {
		fmt.Printf("Snapshot restore Error:  %v\n",err)
	}
*/

/*
	err = sys.WriteSnapshot("diskcopy")
	if err != nil {
		fmt.Printf("Snapshot write Error:  %v\n",err)
	}


	err = sys.ReadSnapshot("diskcopy")
	if err != nil {
		fmt.Printf("Snapshot read Error:  %v\n",err)
	}

	err = sys.RestoreSnapshot("diskcopy")
	if err != nil {
		fmt.Printf("Snapshot restore Error:  %v\n",err)
	}


	tuples2[0].SetValue(1, "")
	tuples2[0].SetValue(2, 0)
	tuples2[0].SetValue(3, 0.0)

	err, v21 := tuples2[0].GetValue(1)	
	err, v22 := tuples2[0].GetValue(2)	
	err, v23 := tuples2[0].GetValue(3)
	fmt.Printf("Values BEFORE SELECT:  [%v] %v %v \n", v21, v22, v23)




	err = sys.Select(tuples2)
	if err != nil {
		fmt.Printf("Error selecting after updating: %v\n", err)
	}

	err, v11:= tuples2[0].GetValue(0)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v12 := tuples2[0].GetValue(1)
	if err != nil {
		fmt.Printf("%v",err)
	}
	
	err, v13 := tuples2[0].GetValue(2)
	if err != nil {
		fmt.Printf("%v",err)
	}

	err, v14 := tuples2[0].GetValue(3)
	if err != nil {
		fmt.Printf("%v",err)
	}


	fmt.Printf("Values after read/write:  %v [%v] %v %v\n", v11, v12, v13, v14)


	system.PrintSchema()
*/
}

