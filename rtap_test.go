

package rtap 

import (
	"fmt"
//	"rtap/message_q"
	"rtap/metronome"
	"testing"
)

const checkMark = "\u2713"
const ballotX 	= "\u2717"




func TestDownload(t *testing.T) {
	t.Log("rtap_test starting.")


	realm, err := Realm("devscada")
	if err != nil {
		fmt.Printf("Error getting realm: %v", err)
	}

	realm.Construct()
	if err != nil {
		fmt.Printf("Error constructing realm: %v", err)
	}


	mq, err := MessageQueue("devscada", "prod")
	if err != nil {
		fmt.Printf("Error getting message queue: %v", err)
	}

	mq.Start()

	ch, err := mq.Register("mychannel")
	if err != nil {
		fmt.Printf("Error registering: %v", err)
	}
	ch = ch

	go metronome.Metronome(mq)


	// Loop forever, receiving from the channel
	for {
		msg := <- ch
		fmt.Printf("X Received msg: source [%s] destination [%s] Payload length: %d\n", msg.Source, msg.Destinations, len(*msg.Data))
		

	//	err = proto.Unmarshall()


		// ---------------------------------------------------------------------
		// return the buffer to the pool
		// ---------------------------------------------------------------------
		mq.PutBuffer(msg.Data)
	}
		
}