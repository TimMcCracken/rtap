/* -----------------------------------------------------------------------------
messageQueue

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


package messageQueue

import (
	"format"
)



// We may want to have an API or interface that hides all these?
MessageQueue 	chan
MessagePoolOut	chan
MessagePooiIn	chan


func init(){

	MessageQueue = make(channel , 256 )
	MessagePoolOut = make(channel , 256 )
	MessagePoolIn = make(channel , 256 )

}



func MessageQueueLoop(){

	 // Loop forever, receiving from the channel
	 for {
        val := <-ch
        fmt.Println("Received:", val)
    }

}