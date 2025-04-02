/*******************************************************************************
domterm.go

domterm is a lowlevel systme for controling the Domain Object Model (DOM) in 
much the same way as an ANSI "dumb" terminal.  It creates JSON messages that
are sebt vua websockets to the browser window that are interpreted and executed
by javascript. Each message describes a Javascript operations such as setValue,
appendChild, SetDocumentTitle, etc.

Widgets and Gadgets use domterm internally to manag their display objects.
Widgets are simple display elements such as TimeOfDay, AnalogValue, 
CounterValue, etc.  Gadgets are more complex display elements such as charts,
etc.


*******************************************************************************/


package domterm


import (
	_ "embed"
	"encoding/json"
//	"flag"
//	"fmt"
//	"html/template"
	"log"
	"maps"
//	"net/http"
//	"time"
	"github.com/gorilla/websocket"
)


// -----------------------------------------------------------------------------
// message is the sturctire that is converted to JSON and sent to the browser
// for all DOM commands.
// -----------------------------------------------------------------------------
type message struct {
	Command		string `json:"command"`
	TargetID	string	`json:"targetID"`
	Data		map[string]string `json:"data"`
}


// -----------------------------------------------------------------------------
// SetDocumentTitle() sends a command to set the title of the HTML document.
// -----------------------------------------------------------------------------
func SetDocumentTitle( conn *websocket.Conn, title string) error {
	var msg message 

	msg.Command = "SetDocumentTitle"
	msg.TargetID = ""
	msg.Data = make(map[string]string)
	msg.Data["title"] = title

	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// SetDocumentTitle() sends a command to set the "value" attribute of a DOM
// element.
// -----------------------------------------------------------------------------
func SetValue( conn *websocket.Conn, targetID string, value string) error {
	var msg message 

	msg.Command = "SetValue"
	msg.TargetID = targetID
	msg.Data = make(map[string]string)
	msg.Data["value"] = value

	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return err
	}
	return nil
}


// -----------------------------------------------------------------------------
// SetAttributes() sends a command to set multiple attributes of a DOM element.
// -----------------------------------------------------------------------------
func SetAttributes( conn *websocket.Conn, targetID string, attributes map[string]string) error {
	var msg message 

	msg.Command = "SetAttributes"
	msg.TargetID = targetID
	msg.Data = attributes

	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// SetAttributes() sends a command to create a new dom element and append it to
// an existing element.
// -----------------------------------------------------------------------------
func AppendElement( conn *websocket.Conn, targetID string, tag string, attributes map[string]string) error {
	var msg message 

	msg.Command = "AppendElement"
	msg.TargetID = targetID
	msg.Data = make(map[string]string)
	msg.Data["tag"]=tag
	maps.Copy( msg.Data, attributes)

	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return err
	}
	return nil
}

// -----------------------------------------------------------------------------
// SetAttributes() sends a command to set multiple style properties on an a DOM 
// element.
// -----------------------------------------------------------------------------
func SetStyle( conn *websocket.Conn, targetID string, properties map[string]string) error {
	var msg message 

	msg.Command = "SetStyle"
	msg.TargetID = targetID
	msg.Data = properties
	jmsg, err := json.Marshal(&msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}
	err = conn.WriteMessage(1, []byte(jmsg))
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return err
	}

	return nil
}

