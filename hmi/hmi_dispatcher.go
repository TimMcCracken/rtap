/* -----------------------------------------------------------------------------
hmi_dispatcher

hmi_dispatcher implements a message queuing system similar to message_q, but
which provided an event driven asynchronous communication system. 
HMI_Dispatcher receives message from message_q, examines the data and forwards
the messages to worker goroutines that were spawned by hmi_server. Initially
it will forward all messages to all worker goroutines, but eventually this will
become a publish/subscribe system.


Rev Date     By  Reason
--- -------- --- ---------------------------------------------------------------


------------------------------------------------------------------------------*/


package hmi

import (
	"fmt"
//	"path/filepath"
//	"runtime"
//	"rtap/common"
	bp "rtap/buffer_pool"
	mq "rtap/message_q"
	"github.com/google/uuid"
//	"rtap/rtdsms"
)


/*
type Message struct{
	bufferSize		int
	messageLength	int
	Source			string
	Destinations	[]string
	Data			*[]byte
}*/





// -----------------------------------------------------------------------------
// HMIDispatcher provides the API for the message queueing system. There is one 
// HMIDispatcher for each domain. 
// -----------------------------------------------------------------------------
type HMIDispatcher struct {
	bp				bp.BufferPool	
	hd_chan	chan 	mq.Message
	receivers		map[uuid.UUID]chan mq.Message
}



// -----------------------------------------------------------------------------
// Start() launches the goroutines that operate and support the message_queue
// system.
// -----------------------------------------------------------------------------
func (hd * HMIDispatcher )Start() {

	if hd.hd_chan == nil {
		hd.hd_chan 	= make(chan mq.Message, 256 )
	}
	
	if hd.receivers == nil {
		hd. receivers = make(map[uuid.UUID]chan mq.Message, 256)
	}


	go hmiDispatcherLoop(hd)

	fmt.Println("HMIDispatcher is started")
}



// -----------------------------------------------------------------------------
// Register allows a display goroutine to register to recieve messages from 
// other HMIDispatcher. Each display goroutine register when in starts and
// unregister when it exits.
// 
//
// -----------------------------------------------------------------------------
func (hd * HMIDispatcher )Register() ( chan mq.Message, uuid.UUID, error) { 
	
	// -------------------------------------------------------------------------
	// Create a UUID to use for the key in the channel map
	// -------------------------------------------------------------------------
	uid := uuid.Must(uuid.NewRandom())

	// Todo check for duplicate and loop around until we get a unique on

	// -------------------------------------------------------------------------
	// Create a new channel and add it to the map
	// -------------------------------------------------------------------------
	new_channel := make(chan mq.Message)
	hd.receivers[uid] = new_channel


	return new_channel, uid, nil
}

// -----------------------------------------------------------------------------
// Send() sends a message to the message queuing system with a byte arrray
// parameter. Normally the byte array is generated as a protocol buffer. This
// facilitates connections by external systems.
// --------------------------------
// ---------------------------------------------

/* This function is probably not necessary or even desired
func (hd * HMIDispatcher )Send( destinations []string, data *[]byte) error { 

	// -------------------------------------------------------------------------
	// determine the senders name
	// -------------------------------------------------------------------------
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get caller information")
	}
	source := filepath.Base(filename)
	extension := filepath.Ext(filename)
	name := source[0:len(source)-len(extension)]

	// -------------------------------------------------------------------------
	// build the message
	// -------------------------------------------------------------------------
	msg := Message {
		Source : name,
		Destinations : destinations,
		Data: data,
	}

	// -------------------------------------------------------------------------
	// send the message to the hmi_dispatcher
	// -------------------------------------------------------------------------
	hd.hd_chan <- msg
	return nil
}*/


func Receive( ) (data []byte, err error) { 

	return nil, nil
}

//------------------------------------------------------------------------------
// The HMIDispatcherLoop() runs continuosly from startup to shutdown. It receives
// messages from the HMIDispatcher, and forwards the messages as needed.
//
//------------------------------------------------------------------------------

func hmiDispatcherLoop(hd * HMIDispatcher){

	fmt.Println("hmi dispatcher loop started")

	//--------------------------------------------------------------------------
	// Loop forever, receiving from the channel and sending out to other
	// channels as needed 
	//--------------------------------------------------------------------------
	for {
        msg := <- hd.hd_chan

		for _, ch := range hd.receivers {
			// We need to make a copy of the message.
			msg2 := msg
			// get another buffer
			buf_ptr := hd.bp.Get(1024).(*[]byte)  // fixed size for now, may change later
			// copy the buffers and assign the buffer ptr
			copy(*buf_ptr, *msg.Data)
			*msg2.Data = *buf_ptr
			ch <- msg2
		}	
		// return the orginal message to the pool
		hd.bp.Put(msg.Data)
	}
}

