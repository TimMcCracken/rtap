/* -----------------------------------------------------------------------------
message_q

messageQueue implements a message queuing system which provided an event driven
asynchronous communication system to allow go-routines to communicate send
messages between each other.  Messages are always in "simplec" (there is no
dequest-response mechanism). There is also a broadcast and multicast capability.

The message queuing system is mmuch like a post office. go-routines connect
to the messageQueue channel (a singleton) and send messages that are directed
at another well known channel that has also connected to the message queue.
When a messatge is received, it is forwarded to the another channel. 

Multicast messages are sent to a list of channeld determined by the sender.
Broadcast messages are sent to all receivers.

Broadcast and muulticast messages are copied (one to each receiver)


Rev Date     By  Reason
--- -------- --- ---------------------------------------------------------------


------------------------------------------------------------------------------*/


package message_q

import (
	"fmt"
	"path/filepath"
	"runtime"
	bp "rtap/buffer_pool"
//	"rtap/common"

//	"rtap/domain"
//	"rtap/rtdsms"
)



type Message struct{
	bufferSize		int
	messageLength	int
	Source			string
	Destinations	[]string
	Data			*[]byte
}

// -----------------------------------------------------------------------------
// MessageQ provides the API for the message queueing system. There is one 
// messageQ for each domain. 
// -----------------------------------------------------------------------------
type MessageQ struct {
	bp				* bp.BufferPool
	mq_chan		 	chan Message
	receivers		map [string]chan Message
}



// -----------------------------------------------------------------------------
// Start() launches the goroutines that operate and support the message_queue
// system.
// -----------------------------------------------------------------------------
func (mq * MessageQ )Start(bp * bp.BufferPool) {

	mq.bp = bp

	if mq.mq_chan == nil {
		mq.mq_chan 	= make(chan Message, 256 )
	}
	
	if mq.receivers == nil {
		mq. receivers = make(map[string]chan Message, 256)
	}

	go messageQueueTask(mq, bp)

	fmt.Println("message queue is started")
}



// -----------------------------------------------------------------------------
// Register allows a goroutine to register to recieve messages from other
// goroutines. If the "name" parameter is nil, it will derive a name from the
// calling functions file name.
// 
// The function checks if a channel already esists with that name. If yes, it
// returns the existing pointer, otherwise it creates a new pointer. The pointer
// is most likely to exist if the previous gorotuine crashed and a new one was
// started by the watch dog timer.
//
// -----------------------------------------------------------------------------
func (mq * MessageQ )Register( name string) ( chan Message, error) { 
	
	// validate the name may want to allow more characters here
//	err := common.ValidateObjectName(name)
//	if err != nil {
//		return nil, fmt.Errorf("Invalid object name [%s]\n", name)
//	}

	// -------------------------------------------------------------------------
	// If the name already exists in the map, then return the pointer
	// -------------------------------------------------------------------------
	value, ok := mq.receivers[name]
	if ok == true {
		return value, nil
	}

	// -------------------------------------------------------------------------
	// Otherwise create a new channel and return a pointer to it.
	// -------------------------------------------------------------------------
	new_channel := make(chan Message)
	mq.receivers[name] = new_channel

	return new_channel, nil
}

// -----------------------------------------------------------------------------
// Send() sends a message to the message queuing system with a byte arrray
// parameter. Normally the byte array is generated as a protocol buffer. This
// facilitates connections by external systems.
// --------------------------------
// ---------------------------------------------
func (mq * MessageQ )Send( destinations []string, data *[]byte) error { 


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
	// copy the data to data from the buffer pool
	// -------------------------------------------------------------------------
	buffer := mq.bp.Get(1024)
	copy(*buffer, *data)

	// -------------------------------------------------------------------------
	// build the message
	// -------------------------------------------------------------------------
	msg := Message {
		Source : name,
		Destinations : destinations,
		Data: buffer,
	}

	// -------------------------------------------------------------------------
	// send the message to the message_q
	// -------------------------------------------------------------------------
	mq.mq_chan <- msg
	return nil
}



func (mq * MessageQ )Receive( ch chan Message  ) (data []byte, err error) { 
	
	// Wait on a message
	msg := <- ch

	// copy the data from the message
	var buffer []byte
	copy(buffer, *msg.Data)

	// return the buffer to the pool
	mq.bp.Put(msg.Data)

	// return
	return buffer, nil
}


//------------------------------------------------------------------------------
// messageQueuLoop() runs continuosly from startup to shutdown. It receives
// messages on a well know channel, and forwards the messages as requested.
// There are 3 types of destinations: unicast, multicast and broadcast. 
//
// Unicast is sent to a single receiver.
// Multicast is sent to a slice of rregistered eceivers.
// Broadcast is sent to all registerd receivers.
//
// For multicast and braodcast, we send a copy of each message.
//
//------------------------------------------------------------------------------
func messageQueueTask(mq * MessageQ, bp * bp.BufferPool){

	fmt.Println("message queue loop started")


	//--------------------------------------------------------------------------
	// Loop forever, receiving from the channel and sending out to other
	// channels as needed 
	//--------------------------------------------------------------------------
	for {

        msg := <- mq.mq_chan

		for key, _ := range mq.receivers {
			fmt.Printf("Receiver: %s\n", key)
		}

		// If there are no listeners, then return the buffer to the pool
		if len(mq.receivers) == 0 {
			fmt.Println("No receivers")
			bp.Put(msg.Data)
			continue
		}

		switch len(msg.Destinations) {

			case 0:
				// todo: log a message
				fmt.Printf("Message queue receved message with 0 destinations.")
				bp.Put(msg.Data)
			case 1:

				if msg.Destinations[0] == "*" { // True if "broadcast"

					for _, ch := range mq.receivers {
						// We need to make a copy of the message.
						msg2 := msg
						// get another buffer
						buf_ptr := bp.Get(1024) //.(*[]byte)   // VAR IGNORED FOR NOW

						// copy the buffers and assign the buffer ptr
						copy(*buf_ptr, *msg.Data)
						*msg2.Data = *buf_ptr
						ch <- msg2
					}	
					// return the original buffer to the pool
					bp.Put(msg.Data)
				} else { // falls here if unicast
					ch, ok := mq.receivers[msg.Destinations[0]]
					if ok == true {
						ch <- msg
					} else {
						fmt.Printf("Message queue unknown destination: [%s]", msg.Destinations)
					}
				}
			default: // true if multicast
				for  _, rcvr_key := range msg.Destinations {
					ch, ok := mq.receivers[rcvr_key]
					if ok == true {
						ch <- msg
					} else {
						fmt.Printf("Message queue unknown destination: [%s]", msg.Destinations)
					}
				}
		}
	}
}



