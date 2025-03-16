/*******************************************************************************
hmi.go

hmi.go is the main entry point of the HMI system.  There is one hmi object
contained within each domain. 

hmi_loop recives messages from and sends messages to the message_q. Messages
recieved from the message_q may be either broadcast messages (normally from
metronome) or unicast messages, normally received from DAC or the RTDSMS.

hmi_server maintains the web sockets server, waiting for connections requests
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
	"encoding/json"
//	"flag"
	"fmt"
	"html/template"
	"log"
	"maps"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
)



var upgrader = websocket.Upgrader{} // use default options

/*
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}
	*/


type HMIMessage struct {
	Command		string `json:"command"`
	TargetID	string	`json:"targetID"`
	Data		map[string]string `json:"data"`
}



func setDocumentTitle( conn *websocket.Conn, title string) {
	var msg HMIMessage 

	msg.Command = "SetDocumentTitle"
	msg.TargetID = ""
	msg.Data = make(map[string]string)
	msg.Data["title"] = title

	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return
	}
}

func setValue( conn *websocket.Conn, targetID string, value string) {
	var msg HMIMessage 

	msg.Command = "SetValue"
	msg.TargetID = targetID
	msg.Data = make(map[string]string)
	msg.Data["value"] = value

	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return
	}
}


func setAttributes( conn *websocket.Conn, targetID string, attributes map[string]string) {
	var msg HMIMessage 

	msg.Command = "SetAttributes"
	msg.TargetID = targetID
	msg.Data = attributes

	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return
	}
}

func appendElement( conn *websocket.Conn, targetID string, tag string, attributes map[string]string) {
	var msg HMIMessage 

	msg.Command = "AppendElement"
	msg.TargetID = targetID
	msg.Data = make(map[string]string)
	msg.Data["tag"]=tag
	maps.Copy( msg.Data, attributes)

	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return
	}
}

func setStyle( conn *websocket.Conn, targetID string, properties map[string]string) {
	var msg HMIMessage 

	msg.Command = "SetStyle"
	msg.TargetID = targetID
	msg.Data = properties
	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return
	}
}




func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Temporary code for testing	
	setDocumentTitle(conn, "The damn title!")

	attributes := make(map[string]string)
	attributes["title"] = "This is an another add in"
	attributes["tag"] = "input"
	appendElement(conn, "body","input", attributes)
	clear(attributes)

	attributes["title"] = "Hello Again from the world"
	setAttributes(conn, "a1", attributes)
	clear(attributes)


	attributes["background-color"] = "lightgrey"
	attributes["color"] = "red"
	attributes["font-weight"] = "bold"

	setStyle(conn, "a1", attributes)

	for {

	//	attributes["value"] = time.Now().Format("2006-01-02 15:04:05")
		setValue(conn, "a1", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(1 * time.Second)
	}

	// end temporary code
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}



func HmiServer() {
	
//	flag.Parse()
//	log.SetFlags(0)

	
	http.HandleFunc("/", home)
	http.HandleFunc("/ws", wsHandler)

	serverAddress := ":8080"
	fmt.Println("WebSocket server started on", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}


//go:embed home.html
var homeHTML string

var homeTemplate = template.Must(template.New("").Parse(homeHTML))
