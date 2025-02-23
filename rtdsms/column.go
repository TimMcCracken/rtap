
/*******************************************************************************
column.go



Rev Date     By  Reason
--- -------- --- ---------------------------------------------------------------

*******************************************************************************/

package rtdsms



type column_intfc interface {
	take_snapshot()
	restore_snapshot()
	write_snapshot()
	read_snapshot()
	clear_snapshot()
}


type column struct {
//	descriptor		*ColumnDescriptor
	realtime		any
	index			map[string]uint64
	snapshots		[]any // data is copied here for disk storage and 
						  //replication purposes, capturing post disturbance analysis
						  // data, etc
	contexts		[]any // for advanced apps that run simulations, etc
						 // and operate other tha on real time data
}