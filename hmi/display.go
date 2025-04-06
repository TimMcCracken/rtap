package hmi

import (
	_ "embed"
//	"encoding/json"
	"fmt"
//	"log"
	"net/http"
//	"time"
//	"github.com/yuin/gopher-lua"
//	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
//	"rtap/hmi/domterm"
//	mq "rtap/message_q"
	"rtap/hmi/widget"
//	"rtap/domain"
)

//go:embed display_test.lua
var  DisplayTest string


type Display struct {
	WidgetMap map[string]widget.Widget
	ClockMap  map[string]widget.Widget
}


var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

func NewDisplay() * Display {

	display := new(Display)

	display.WidgetMap = make(map[string]widget.Widget)
	display.ClockMap  = make(map[string]widget.Widget)

	return display
}



func (d * Display) NewLabel(parent string, top int, left int, width int, height int, zIndex int, content string ) (*widget.Label, error) {

	lbl := new( widget.Label)

	// make up an ID
	displayID := fmt.Sprintf("lbl_%d", len(d.WidgetMap))

	err := lbl.Init(displayID, parent, top, left, width, height, zIndex, content)
	if err != nil {
		return lbl, err
	}

	// add to the map
	d.WidgetMap[lbl.DisplayID] = lbl

	// return the new label and error code
	return lbl, nil
}

func (d * Display) NewDigitalClock(parent string, top int, left int, width int, height int, zIndex int, content string ) (*widget.DigitalClock, error) {

	dc := new( widget.DigitalClock)

	// make up an ID
	displayID := fmt.Sprintf("dc_%d", len(d.WidgetMap))

	err := dc.Init(displayID, parent, top, left, width, height, zIndex, content)
	if err != nil {
		return dc, err
	}

	// add to the maps
	d.WidgetMap[dc.DisplayID] = dc
	d.ClockMap[dc.DisplayID] = dc

	// return the new label and error code
	return dc, nil
}




// -----------------------------------------------------------------------------
// Show() should be called after creating all the display objects
// -----------------------------------------------------------------------------
func (d * Display) Show(conn *ws.Conn) {

	for _, widget := range d.WidgetMap {
		widget.Show(conn)
	}

}


// -----------------------------------------------------------------------------
// DisplayHandler() is a goroutine that is spawned by the HMI listener to 
// handle create, execute and close displays.
// -----------------------------------------------------------------------------
/*
func DisplayHandler(domain domain.Domain, w http.ResponseWriter, r *http.Request) {

	// Parse the URL variables
//	vars := mux.Vars(r)
//	realm := vars["realm"]
//	domain := vars["domain"]
//	name := vars["name"] //display name

	// Validate the URL variables


//	fmt.Printf("Realm : %s\n", realm)
//	fmt.Printf("Domain: %s\n", domain)
//	fmt.Printf("D Name: %s\n", name)




	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()



	fmt.Println("Starting display handler.")

	// authenticate the connection

	// Get the realm, if any. make smart enought to default if only 1

	// Get the domain, if any. make smart enought tp default if only 1



	// -------------------------------------------------------------------------
	// Create a display object
	// -------------------------------------------------------------------------
	display := Display{
		widgetMap : make(map[string]widget.Widget),
		clockMap  : make(map[string]widget.Widget),
	}


	// -------------------------------------------------------------------------
	// Create a lua state
	// -------------------------------------------------------------------------
	L := lua.NewState()
	defer L.Close()
	
	registerDisplayType(L)
	widget.RegisterLabelType(L)
	widget.RegisterDigitalClockType(L)

	// Add it to Lua
	ud := L.NewUserData()
	ud.Value = &display
	L.SetMetatable(ud, L.GetTypeMetatable("display"))
	L.Push(ud)
	L.SetGlobal("display", ud)

	
	
	
	// -------------------------------------------------------------------------
	// hmiChan is a channel to receive messages from the hmiTask that processes 
	// messages from the messageQ
	// -------------------------------------------------------------------------
	hmiChan := make(chan mq.Message)
	
	// -------------------------------------------------------------------------
	// clientChan is a channel to receive messages from the client via the 
	// web socket connetion..
	// -------------------------------------------------------------------------
	clientChan 	:= make(chan []byte)

	// -------------------------------------------------------------------------
	// tickChan is used to keep digital clocks updated. without it, we would
	// get concurrent access errors on the web socket connection.
	// -------------------------------------------------------------------------
	tickChan 	:= make(chan int)

	// -------------------------------------------------------------------------
	// The following is an anonymous go routine that recives messages from the
	// web socket and foreads them to the client channel.
	// -------------------------------------------------------------------------
	go func(){
		for {
			msgtype, payload, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			} else {
				fmt.Printf("Message Type %v\n", msgtype)
				
				payload = payload
				msgtype = msgtype
	
				clientChan <- payload
			}
		}
	}()


	fmt.Println("Starting lua.")
	if err := L.DoString(display_test); err != nil {

		fmt.Printf("Lua Error: %v\n", err)

		//panic(err)
	}
	fmt.Println("Finished Lua DoString().")


	// Get the "main" function from Lua
	mainFunc := L.GetGlobal("main")
	if mainFunc.Type() == lua.LTFunction {
		// Call main() in Lua
		if err := L.CallByParam(lua.P{
			Fn:      mainFunc,
			NRet:    0, // Number of return values expected
			Protect: true,
		}); err != nil {
			fmt.Println("Error calling Lua function:", err)
		}
	} else {
		fmt.Println("Error: 'main' function not found in Lua state.")
	}

	fmt.Println("Finished Lua main.")



	// -------------------------------------------------------------------------
	// The following is test code while we more fully develop the dHMI sytem
	// -------------------------------------------------------------------------
	// set the page style
	attributes := make(map[string]string)
	attributes["background-color"] = "Black"
	attributes["color"] = "White"
	domterm.SetStyle( conn, "html", attributes)
	clear(attributes)


	// -------------------------------------------------------------------------
	// Send all the display objects to the client browser
	// -------------------------------------------------------------------------

	fmt.Println("starting show()")

	display.Show(conn)

	fmt.Println("finished show()")


	// -------------------------------------------------------------------------
	// Create a goroutine to update the digital clocks if any
	// -------------------------------------------------------------------------
	go func(){

		for {
			tickChan <- 0
			time.Sleep(1 * time.Second)
		}

	}()

	
	// -------------------------------------------------------------------------
	// The following is our loop that runs until the connection is terminated.
	// This code responds to all the events coming in.
	// -------------------------------------------------------------------------
	for {
		select {

			// ----------------------------------------------------------------
			// These will be events sent from RTAP via the messageQ and HMITask
			// ----------------------------------------------------------------
			case _ = <- tickChan:
        		fmt.Printf("received tick %v\n", time.Now().Unix())

				for _, dc := range display.clockMap {
					err = dc.Update(conn)
					if err != nil {
						fmt.Printf("display error: %v\n", err)
						return
					}
				}


			// ----------------------------------------------------------------
			// These will be events sent from RTAP via the messageQ and HMITask
			// ----------------------------------------------------------------
			case hmiMsg := <- hmiChan:
        	fmt.Println("received msqgq", hmiMsg)
      

			// ----------------------------------------------------------------
			// These will be events sent from the client browser
			// ----------------------------------------------------------------
			case clientMsg := <- clientChan:

			var data map[string]any
			err := json.Unmarshal(clientMsg, &data)
			if err != nil {
				fmt.Printf("!!!!!!!!!!!!!!!!  doing the fatal error\n")
				log.Fatal(err)
			}
		
			target_id, ok := data["id"]
			if ok == false {
				fmt.Printf("Error parsing target id\n")
			} else {
			
				widget, ok := display.widgetMap[target_id.(string)]
				if ok == false {
					fmt.Printf("Error parsing target id\n")
				} else {
					widget.ClientEvent(data)
				}
			}

		}			
	}
}
*/
