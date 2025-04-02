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
//	"rtap/domain"	
//	bp "rtap/buffer_pool"
//	mq "rtap/message_q"

	//	"rtap/hmi/widget"
)


todo: combine HMI and HMI dispatcher to one struct

//var upgrader = websocket.Upgrader{} // use default options


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}


type HMI struct {

}




func (hmi * HMI) Start( serverAddress string){
	go hmi.HMIServerTask(serverAddress)
}



func homeStub(w http.ResponseWriter, r *http.Request) {

//	w.Write([]byte("WTF"))
	authenticateTemplate.Execute(w, "ws://" +r.Host+ "/ws/authenticate")
}

func chooserStub(w http.ResponseWriter, r *http.Request) {
	chooserTemplate.Execute(w, "ws://" +r.Host+ "/ws/chooser")
}

func displayStub(w http.ResponseWriter, r *http.Request) {
	displayTemplate.Execute(w, "ws://" +r.Host+ "/ws/display")
}


func domtermHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(jsDomterm))
}


// -----------------------------------------------------------------------------
// HMIServerTask is the web sockets listener
// -----------------------------------------------------------------------------
func (hmi * HMI) HMIServerTask(serverAddress string) {
	
	// Build the private messageQueu that we will use
	// to send messages to the websocket goroutines

	http.HandleFunc("/", homeStub)
	http.HandleFunc("/ws/authenticate", wsAuthenticateHandler)


	http.HandleFunc("/chooser", chooserStub)
	http.HandleFunc("/ws/chooser", wsChooserHandler)


	http.HandleFunc("/display", displayStub)
	http.HandleFunc("/ws/display", wsDisplayHandler)


	http.HandleFunc("/js/domterm.js", domtermHandler)


	// Start the HMI Server loop.
	fmt.Println("WebSocket server started on", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))

}





// The following commented lines are go directives, not a comment!

//go:embed authenticate.html
var authenticate 	string

//go:embed chooser.html
var chooser			string

//go:embed display.html
var display_html			string

//go:embed domterm.js
var  jsDomterm			string


var authenticateTemplate 	= template.Must(template.New("").Parse(authenticate))
var chooserTemplate 		= template.Must(template.New("chooser").Parse(chooser))
var displayTemplate 		= template.Must(template.New("display").Parse(display_html))
