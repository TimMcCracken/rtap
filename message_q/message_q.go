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
)



type message struct{
	source			string
	destination		string
	paylaod			[]byte
	channel			*chan []byte	// used only for registering an application with the 
							// message quque
}



// We may want to have an API or interface that hides all these?
var messageQueue 	chan message
var messagePoolOut	chan []byte
var messagePoolIn	chan []byte


func init(){

	messageQueue 	= make(chan message, 256 )
	messagePoolOut	= make(chan []byte, 256 )
	messagePoolIn	= make(chan []byte, 256 )

	fmt.Println("message queue is initd")
}

func Start() {
	go MessageQueueLoop()
}

func Send( destination string, data *any) error { 

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get caller information")
	}

	msg := message {
		source : filepath.Base(filename),
		destination : destination,
	}

	messageQueue <- msg

	return nil
}

func Receive( ) (data *any, err error) { 
  
	return nil, nil
}


func MessageQueueLoop(){

	fmt.Println("message queue loop started")


	 // Loop forever, receiving from the channel
	 for {
        msg := <- messageQueue
        fmt.Printf("Received msg: source [%s] destination [%s]\n", msg.source, msg.destination)
    }

}