/*******************************************************************************
dacc.go

This is the main file for 'DACC', Data Acquisition, Command & Control.



*******************************************************************************/

package dacc

import (

	"fmt"
	"time"
	proto 	"github.com/golang/protobuf/proto"

	bp 		"rtap/buffer_pool"
	mq 		"rtap/message_q"
	anlg 	"rtap/dacc/dacc_events/analog.pb"

)




type DACC struct {

}


func (dacc * DACC) Start (bp * bp.BufferPool, mq * mq.MessageQ){

	go dacc.daccTask(bp, mq)
}



func (dacc * DACC) daccTask (bp * bp.BufferPool, mq * mq.MessageQ){

	var ac anlg.AnalogChange

	// the following is just test code used to develop a data path
	// through the message queue  to the HMI

	x := float64(100.0)

	for {

		ac.Value = x

		// -------------------------------------------------------------------------
			// Get a buffer from the pool
		// -------------------------------------------------------------------------
		buf_ptr := bp.Get(1024).(*[]byte)  // fixed size for now. May change later
		
		// -------------------------------------------------------------------------
		// Marshall the data and copy it into the buffer
		// -------------------------------------------------------------------------
		data, err := proto.Marshal(&ac)
		if err != nil {
			// TODO: fix this
    		fmt.Printf("Failed to encode address book:", err)
		}

		// -------------------------------------------------------------------------
		// Ensure buffer is large enough
		// -------------------------------------------------------------------------
		if len(data) > cap(*buf_ptr) {
			fmt.Printf("Buffer is too small.")
		}

		// -------------------------------------------------------------------------
		// Copy marshaled data into the buffer to minimize allocations
		// -------------------------------------------------------------------------
		copy(*buf_ptr, data)		

		// -------------------------------------------------------------------------
		// Send the data to the message_q.
		// -------------------------------------------------------------------------
		destinations := []string {"HMI"}
		mq.Send(destinations, buf_ptr)



//		fmt.Printf("!")
		x += 100
		time.Sleep(1 * time.Second)
	}


}