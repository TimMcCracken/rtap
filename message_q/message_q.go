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
	"sync"

	"rtap/rtdsms"
)



type Message struct{
	buffer_size		int
	message_length	int
	Source			string
	Destinations	[]string
	Data			*[]byte
}

// -----------------------------------------------------------------------------
// MessageQ provides the API for the message queueing system. There is one 
// messageQ for each domain.
// -----------------------------------------------------------------------------
MessageQ {
	
	var messageQueue 	chan Message
	var scrubber		chan *[]byte
	var bufferPool		sync.Pool
	var receivers		map [string]chan Message
	var total_buffers	int
	var total_gets		int
	var total_puts 		int
	var total_scrubs	int
}



// -----------------------------------------------------------------------------
// Start() launches the goroutines that operate and support the message_queue
// system.
// -----------------------------------------------------------------------------
func (MessageQ * mq)Start() {

	if mq.MessageQ == nil {
		messageQueue 	= make(chan Message, 256 )
	}
	
	if mq.scrubber == nil {
		scrubber		= make(chan *[]byte, 256 )
	}

	if mq.receivers == nil {
		mq. receivers = make(map[string]chan Message, 256)
	}

	if bufferPool == nil {
		bufferPool = sync.Pool{
			New: func() interface{} {
				buf := make([]byte, 1024) // Allocate a 1KB buffer (adjust size as needed)
				total_buffers++
				return &buf
			},
		}
	}

	go messageQueueLoop()
	go scrubberLoop()

	fmt.Println("message queue is started")
}


// -----------------------------------------------------------------------------
// GeBuffer() returns a buffer from the buffer pool.
// -----------------------------------------------------------------------------
func GetBuffer() ( interface{}) {
	total_gets++
	fmt.Printf("Total Buffers, Gets, Puts, scrubs: %d  %d  %d  %d\n", total_buffers, total_gets, total_puts, total_scrubs)
	return bufferPool.Get().(*[]byte)
}

// -----------------------------------------------------------------------------
// PutBuffer() returns a buffer to the buffer pool/
// -----------------------------------------------------------------------------
func PutBuffer(buffer interface{}){
	total_puts++
	bufferPool.Put(buffer)
}


// -----------------------------------------------------------------------------
// Register allows a goroutine to register to recieve messages from other
// goroutines. If the "name" parameter is nil, it will derive a name from the
// calling functions file name.
// 
// The function checks if a channel already esists with that name. If yes, it
// returns the existing pointer, otherwise it creates a new pointer. The pointer
// is most l;kely to exist if the previous gorotuine crashed and a new one was
// started by the watch dog timer.
//
// -----------------------------------------------------------------------------
func Register( name string) ( chan Message, error) { 
	
	// validate the name
	err := rtdsms.ValidateObjectName(name)
	if err != nil {
		return nil, fmt.Errorf("Invalid object name [%s]\n", name)
	}

	// -------------------------------------------------------------------------
	// If the name already exists in the map, then return the pointer
	// -------------------------------------------------------------------------
	value, ok := receivers[name]
	if ok == true {
		return value, nil
	}

	// -------------------------------------------------------------------------
	// Otherwise create a new channel and return a pointer to it.
	// -------------------------------------------------------------------------
	new_channel := make(chan Message)
	receivers[name] = new_channel

	return new_channel, nil
}

// -----------------------------------------------------------------------------
// Send() sends a message to the message queuing system with a byte arrray
// parameter. Normally the byte array is generated as a protocol buffer. This
// facilitates connections by external systems.
// --------------------------------
// ---------------------------------------------
func Send( destinations []string, data *[]byte) error { 

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
	// send the message to the message_q
	// -------------------------------------------------------------------------
	messageQueue <- msg
	return nil
}


func Receive( ) (data []byte, err error) { 

	return nil, nil
}

//------------------------------------------------------------------------------
// The messageQueuLoop() runs continuosly from startup to shutdown. It receives
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

func messageQueueLoop(){

	fmt.Println("message queue loop started")

	//--------------------------------------------------------------------------
	// Loop forever, receiving from the channel and sending out to other
	// channels as needed 
	//--------------------------------------------------------------------------
	for {
        msg := <- messageQueue
    
		switch len(msg.Destinations) {

		case 0:
			fmt.Printf("Message queue receved message with 0 destinations.")
		case 1:

			if msg.Destinations[0] == "*" { // True if "broadcast"
				for _, ch := range receivers {
					// We need to make a copy of the message.
					msg2 := msg
					// get another buffer
					buf_ptr := bufferPool.Get().(*[]byte)
					// copy the buffers and assign the buffer ptr
					copy(*buf_ptr, *msg.Data)
					*msg2.Data = *buf_ptr
					ch <- msg2
				}	
				// return the original buffer to the pool
				scrubber <- msg.Data
			} else { // falls here if unicast
				ch, ok := receivers[msg.Destinations[0]]
				if ok == true {
					ch <- msg
				} else {
					fmt.Printf("Message queue unknown destination: [%s]", msg.Destinations)
				}
			}
		default: // true if multicast
			for  _, rcvr_key := range msg.Destinations {
				ch, ok := receivers[rcvr_key]
				if ok == true {
					ch <- msg
				} else {
					fmt.Printf("Message queue unknown destination: [%s]", msg.Destinations)
				}
			}
		}	
	}
}

//------------------------------------------------------------------------------
// scrubber loop zeroes out a buffer before returning it to the pool.
//------------------------------------------------------------------------------
func scrubberLoop(){

//	fmt.Println("scrubber loop started")
	for {
        data := <- scrubber
//		fmt.Printf("Scrubbing %d bytes\n", len(*data))
		for  i := range *data {
			(*data)[i] = 0
		}
		bufferPool.Put(data)
		total_scrubs++
	}
}


