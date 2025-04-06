/*******************************************************************************
hmi.go

hmi.go is ...

HMI uses a publish/subscribe system that operates similar to the message_q
but allows each of the websocket goroutines to subscribe to the objects that
then need to keep their displays upto date.


hmiLoop recives messages from the domain message_q. Messages recieved from the 
message_q may be either broadcast messages (normally from metronome) or unicast
messages, normally received from DACC or the RTDSMS.

hmiServer maintains the web sockets server, waiting for connections requests
from clients. Each connection request creates a display_loop, that runs 
until the client connection is closed - which normally happens when a 
http client is terminated (browser tab closed). 

hmi_loop also recives message from the display_loop established for each
active connection. These consist of subscribe/unsubscribe messages. The HMI
maintains a map of subscribed points for each display_loop, which is used to
track all subscribed data points. JSON messages are exchanged over the 
web socket connection between the client and the hti_loop.


	hmiServer (1 webserver that serves all realms/domains)
	hmiTask   (1 task per realm/domain that receives requests from HMI server 
				and publishes to )
	hmiWorker (spawned by hmiServer and subscribes on hmiTask )




*******************************************************************************/


package hmi

import (
	_ "embed"
//	"encoding/json"
//	"flag"
//	"fmt"
//	"html/template"
//	"log"
//	"maps"
//	"net/http"
//	"time"
//	"github.com/google/uuid"
//	"net/http"
//	"github.com/gorilla/websocket"
//	"rtap/hmi/domterm"
//	"rtap/domain"	
//	bp "rtap/buffer_pool"
//	mq "rtap/message_q"

	//	"rtap/hmi/widget"
)





//var upgrader = websocket.Upgrader{} // use default options




/*
type HMI struct {

	
	bp				* bp.BufferPool	
	hmiTaskChan		chan mq.Message
	receivers		map[uuid.UUID]chan []byte
	
	
}
*/


/*
func (hmi * HMI) Start( serverAddress string){
	
} */


