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






*******************************************************************************/


package hmi

import (
	_ "embed"
//	"encoding/json"
//	"flag"
	"fmt"
	"html/template"
	"log"
//	"maps"
	"net/http"
//	"time"
	"github.com/gorilla/websocket"
//	"rtap/hmi/domterm"
//	"rtap/domain"/	
bp "rtap/buffer_pool"
mq "rtap/message_q"

	//	"rtap/hmi/widget"
)



//var upgrader = websocket.Upgrader{} // use default options


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}


type HMI struct {

}


func (hmi * HMI) Start( bp * bp.BufferPool, mq * mq.MessageQ, serverAddress string){

	go hmi.HMILoop(bp, mq)
	go hmi.HMIServer(serverAddress)
}



func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/testpage")
}


// -----------------------------------------------------------------------------
// HMIServer is the web sockets listener
// -----------------------------------------------------------------------------
func (hmi * HMI) HMIServer(serverAddress string) {
	
	// Build the private messageQueu that we will use
	// to send messages to the websocket goroutines

	http.HandleFunc("/", home)
	http.HandleFunc("/testpage", testPageHandler)

	// Start the HMI Server loop.
	fmt.Println("WebSocket server started on", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))

}


// -----------------------------------------------------------------------------
// HMILoop listens for messages from the MessageQ as well as from the HMIWorker
// goroutines that are spawned by HMI loop. When it receives a message from the
// message_q, it sends it to the worker go routines. When it receives a msg from
// a worker, it .....
// -----------------------------------------------------------------------------
func (hmi * HMI) HMILoop( bp * bp.BufferPool, mq * mq.MessageQ) {

	// TODO FIX THIS ERROR HANDLER
	ch, err := mq.Register("HMI") 
	if err != nil {
		fmt.Printf("ERROR! Invalid object name [HMI]\n")
		return
	}
	
	//var msg mq.Message
	
	for {
		msg := <- ch

		fmt.Printf("Got an HMI msg.\n")

		msg = msg
	
	}


}





// The following line is a go directive, not a comment!
//go:embed home.html
var homeHTML string

var homeTemplate = template.Must(template.New("").Parse(homeHTML))
