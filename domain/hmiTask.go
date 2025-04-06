/* -----------------------------------------------------------------------------
hmi_dispatcher

hmi_dispatcher implements a message queuing system similar to message_q, but
forwards messages to displayWorkers

HMI_Dispatcher receives message from message_q, examines the data and forwards
the messages to worker goroutines that were spawned by hmi_server. Initially
it will forward all messages to all worker goroutines, but eventually this will
become a publish/subscribe system.


Rev Date     By  Reason
--- -------- --- ---------------------------------------------------------------


------------------------------------------------------------------------------*/


package domain

import (
	"fmt"
//	"path/filepath"
//	"runtime"
//	"rtap/domain"

//	bp "rtap/buffer_pool"
//	mq "rtap/message_q"
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
/*
type HMIDispatcher struct {
	bp				* bp.BufferPool	
	hmiTaskChan		chan [] mq.Message
	receivers		map[uuid.UUID]chan []byte
}*/



// -----------------------------------------------------------------------------
// Start() launches the goroutines that operate and support the message_queue
// system.
// -----------------------------------------------------------------------------

/*
hmi.bp = bp

	if hmi.hmiTaskChan == nil {
		hmiTaskChan, err := mq.Register("hmi")
		if err != nil {
			fmt.Printf("Error creating hmiTaskChannel!\n") // todo log error
		}
		hmi.hmiTaskChan = hmiTaskChan
	}
	
	if hmi.receivers == nil {
		hmi. receivers = make(map[uuid.UUID]chan []byte, 256)
	}

	go hmi.hmiDispatcherTask()
	fmt.Println("HMIDispatcher is started")
*/

/*
func (domain * Domain )Start() {

	// This task receives messages from the MessageQ for this domain
	if domain.hmi_chan == nil {
		domain.hmi_chan 	= mq.Register()
	}
	
	// This task forwards messages to the hmi Workers
	if domain.hmiWorkers == nil {
		domain.hmiWorkers = make(map[uuid.UUID]chan []byte, 256)
	}


	go hmiTask(hd)

	fmt.Println("hmiTask() is started on %s %s\n", domain.Descriptor.realm_name, domain.Descriptor.domain_name)
}*/



// -----------------------------------------------------------------------------
// Register allows a displayWorker to register to recieve messages from 
// other HMIDispatcher. Each display goroutine register when in starts and
// unregister when it exits.
// 
//
// -----------------------------------------------------------------------------
func (domain * Domain )RegisterHmiWorker() ( chan []byte, uuid.UUID, error) { 
	
	// -------------------------------------------------------------------------
	// Create a UUID to use for the key in the channel map
	// -------------------------------------------------------------------------
	uid := uuid.Must(uuid.NewRandom())

	// Todo check for duplicate and loop around until we get a unique on

	// -------------------------------------------------------------------------
	// Create a new channel and add it to the map
	// -------------------------------------------------------------------------
	new_channel := make(chan []byte)
	domain.hmiWorkers[uid] = new_channel


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


func (domain * Domain )ReceiveHmiMessage( ) (data []byte, err error) { 

	return nil, nil
}

//------------------------------------------------------------------------------
// The hmiTask() runs continuosly from startup to shutdown. It receives
// messages from the MessageQ, and forwards the messages to the hmiWorkers
//
//------------------------------------------------------------------------------

func (domain * Domain) hmiTask(){

	fmt.Println("hmi task starting")


	// This task receives messages from the MessageQ for this domain
	if domain.hmiChannel == nil {
		hmiChannel, err 	:= domain.messageQueue.Register("hmi")
		if err != nil {
			fmt.Printf("Failed to register HMI with message_q") //tofo fix error log
		}
		domain.hmiChannel = hmiChannel
	}

	// This task forwards messages to the hmi Workers
	if domain.hmiWorkers == nil {
		domain.hmiWorkers = make(map[uuid.UUID]chan []byte, 256)
	}

	fmt.Println("hmiTask() is started on %s %s\n", domain.Descriptor.RealmName, domain.Descriptor.DomainName)


	//--------------------------------------------------------------------------
	// Loop forever, receiving from the channel and sending out to other
	// channels as needed 
	//--------------------------------------------------------------------------
	for {
		// Get a new message from the messageQ
        msg := <- domain.hmiChannel

		// Send the message to all the hmiWorkers. later this will be a publish
		// subscribe system to send messages only to the workers that need the
		// message.
		for _, ch := range domain.hmiWorkers {
			// We need to make a copy of the message.
			//msg2 := msg.Data

			// get another buffer
			buf_ptr := domain.bufferPool.Get(1024)// .(*[]byte)  // fixed size for now, may change later
			
			// copy the buffers and assign the buffer ptr
			copy(*buf_ptr, *msg.Data)
			
			//msg2 = buf_ptr
			ch <- *buf_ptr
		}	
		// return the orginal buffer to the pool
		domain.bufferPool.Put(msg.Data)

		fmt.Printf("hmiTask()\n")
	}
}

