package hmi

import (
	_ "embed"
	"encoding/json"
//	"flag"
	"fmt"
//	"html/template"
	"log"
//	"maps"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
	"rtap/hmi/domterm"
	mq "rtap/message_q"
	"rtap/hmi/widget"
)

type display struct {
	widgetMap map[string]widget.Widget
}


func (d * display) NewLabel(parent string, top int, left int, width int, height int, zIndex int, content string ) (widget.Label, error) {

	newLabel := widget.Label  {
		Parent 	: parent,
		Top 	: top,
		Left 	: left,
		Width 	: width,
		Height 	: height,
		ZIndex	: zIndex,
		Content : content,
	}

	// make up an ID
	newLabel.DisplayID = fmt.Sprintf("lbl_%d", len(d.widgetMap))
	// add to the map
	d.widgetMap[newLabel.DisplayID] = &newLabel
	// return the new label and error code
	return newLabel, nil
}

func (d * display) NewDigitalClock(parent string, top int, left int, width int, height int, zIndex int, content string ) (*widget.DigitalClock, error) {

	dc := new( widget.DigitalClock)

	// make up an ID
	displayID := fmt.Sprintf("dc_%d", len(d.widgetMap))

	err := dc.Init(displayID, parent, top, left, width, height, zIndex, content)
	if err != nil {
		return dc, err
	}

	// add to the map
	d.widgetMap[dc.DisplayID] = dc

	// return the new label and error code
	return dc, nil
}




// -----------------------------------------------------------------------------
// Show() should be called after creating all the display objects
// -----------------------------------------------------------------------------
func (d * display) Show(conn *websocket.Conn) {

	for _, widget := range d.widgetMap {
		widget.Show(conn)
	}

}


// -----------------------------------------------------------------------------
// wsDisplayHandler() is a goroutine that is spawned by the HMI listener to 
// handle create, execute and close displays.
// -----------------------------------------------------------------------------
func wsDisplayHandler(w http.ResponseWriter, r *http.Request) {

	display := display{
		widgetMap : make(map[string]widget.Widget),
	}


	// authenticate the connection

	// Get the realm, if any. make smart enought to default if only 1

	// Get the domain, if any. make smart enought tp default if only 1

	// is the user authorized?

	// Do we need to determine the type of display?

	// Show the display



	fmt.Printf("Starting display handler.")
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()
	
	// -------------------------------------------------------------------------
	// hmiChan is a channel to receive messages from the hmiTask that processes 
	// messages from the messageQ
	// -------------------------------------------------------------------------
	hmiChan := make(chan mq.Message)
	// -------------------------------------------------------------------------
	// clientChan is a channel to receive messages from the client via the 
	// sweb socket connetion..
	// -------------------------------------------------------------------------
	clientChan 	:= make(chan []byte)

	// -------------------------------------------------------------------------
	// The following is an anonymous go routine that recives messages from the
	// web socket and foreads them to the client channel.
	// -------------------------------------------------------------------------
	go func(){
		for {
			msgtype, payload, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
			} else {
				fmt.Printf("Message Type %v\n", msgtype)
				payload = payload
				msgtype = msgtype
	
				clientChan <- payload
			}
		}
	}()




	// -------------------------------------------------------------------------
	// The following is test code while we more fully develop the dHMI sytem
	// -------------------------------------------------------------------------
	// set the page style
	attributes := make(map[string]string)
	attributes["background-color"] = "Black"
	attributes["color"] = "White"
	domterm.SetStyle( conn, "html", attributes)
	clear(attributes)


	display.NewLabel("body", 50, 50, 200, 0, 0, "Local" )
	display.NewLabel("body", 50, 300, 200, 0, 0, "America/New_York" )
	display.NewLabel("body", 50, 550, 200, 0, 0, "UTC" )

	dc1, _ := display.NewDigitalClock("body", 100, 50, 200, 0, 1, "" )
	dc1.Timezone = "Local"
	dc2, _ := display.NewDigitalClock("body", 100, 300, 200, 0, 1, "" ) 
	dc1.Timezone = "New York"
	dc3, _ := display.NewDigitalClock("body", 100, 550, 200, 0, 1, "" ) 
	dc1.Timezone = "UTC"




	// -------------------------------------------------------------------------
	// Send all the display objects to the client browser
	// -------------------------------------------------------------------------
	display.Show(conn)






	// ------------------------------------
	// the following goroutine is test code
	// ------------------------------------
	go func() {

		for {
			dc1.Update(conn)
			dc2.Update(conn)
			dc3.Update(conn)
			time.Sleep(1 * time.Second)
		}
	}()

	
	// -------------------------------------------------------------------------
	// The following is our loop that runs until the connection is terminated.
	// This code responds to all the events coming in.
	// -------------------------------------------------------------------------
	for {
		select {
		// ---------------------------------------------------------------------
		// These will be events sent from RTAP via the messageQ and HMITask
		// ---------------------------------------------------------------------
			case hmiMsg := <- hmiChan:
        	fmt.Println("received msqgq", hmiMsg)
      

		// ---------------------------------------------------------------------
		// These will be events sent from the client browser
		// ---------------------------------------------------------------------
			case clientMsg := <- clientChan:

			var data interface{}
			err := json.Unmarshal(clientMsg, &data)
			if err != nil {
				log.Fatal(err)
			}
		
			fmt.Printf("received client %v\n", data)	
		}			
	}
}

