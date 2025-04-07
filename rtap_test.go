

package rtap 

import (
	"fmt"
//	"rtap"
//	"rtap/message_q"
//	"rtap/metronome"
//	"rtap/hmi"
	"testing"
	"time"
)

const checkMark = "\u2713"
const ballotX 	= "\u2717"




func TestDownload(t *testing.T) {
	t.Log("rtap_test starting.")


	realm, err := RTAP.Realm("devscada")
	if err != nil {
		fmt.Printf("Error getting realm: %v", err)
	} 

	realm.Construct()
	if err != nil {
		fmt.Printf("Error constructing realm: %v", err)
	}

	

	domain, err := realm.Domain("prod")

	domain.Start()


	// Loop foreve
	for {
		//msg := <- ch
		//	fmt.Printf("X Received msg: source [%s] destination [%s] Payload length: %d\n", msg.Source, msg.Destinations, len(*msg.Data))
		//	err = proto.Unmarshall()


		time.Sleep(1 * time.Second) 
		
		
		// ---------------------------------------------------------------------
		// return the buffer to the pool
		// ---------------------------------------------------------------------
		//mq.PutBuffer(msg.Data)
	}
		
}