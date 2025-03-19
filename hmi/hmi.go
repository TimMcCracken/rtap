/*******************************************************************************
hmi.go

hmi.go is the main entry point of the HMI system.  There is one hmi object
contained within each domain. 

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


// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

////go:build ignore
////+build ignore

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
//	"rtap/buffer_pool"
//	"rtap/message_q"

	//	"rtap/hmi/widget"
)



//var upgrader = websocket.Upgrader{} // use default options


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}





func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/testpage")
}



func HmiServer(server_address string) {
	
	// Build the private messageQueu that we will use
	// to send messages to the websocket goroutines

	http.HandleFunc("/", home)
	http.HandleFunc("/testpage", testPageHandler)

	// Start the HMI Server loop.
	fmt.Println("WebSocket server started on", server_address)
	log.Fatal(http.ListenAndServe(server_address, nil))

}


//go:embed home.html
var homeHTML string

var homeTemplate = template.Must(template.New("").Parse(homeHTML))
