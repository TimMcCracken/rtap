package hmi

import (
	_ "embed"
//	"encoding/json"
//	"flag"
	"fmt"
//	"html/template"
	"log"
//	"maps"
	"net/http"
	"time"
//	"github.com/gorilla/websocket"
	"rtap/hmi/domterm"
//	"rtap/message_q"
	"rtap/hmi/widget"
)


// -----------------------------------------------------------------------------
// testPageHandler()
// -----------------------------------------------------------------------------
func testPageHandler(w http.ResponseWriter, r *http.Request) {
	
	fmt.Printf("Starting testpage handler.")
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()
	
	// set the page style
	attributes := make(map[string]string)
	attributes["background-color"] = "Black"
	attributes["color"] = "White"
	domterm.SetStyle( conn, "html", attributes)
	clear(attributes)


	lbl1 := widget.Label  {
		Top : 50,
		Left : 50,
		Width : 200,
		Parent : "body",
		Content : "Local",
	}
	lbl1.Init(nil, "", "", "lbl_101")
	lbl1.Put(conn)

	dc1 := widget.DigitalClock  {
		Top : 100,
		Left : 50,
		Width : 200,
		Timezone : "Local",
		Parent : "body",
	}

	dc1.Init(nil, "", "", "dc_101")
	dc1.Put(conn)


	lbl2 := widget.Label  {
		Top : 50,
		Left : 300,
		Width : 200,
		Parent : "body",
		Content : "New York",
	}
	lbl2.Init(nil, "", "", "lbl_102")
	lbl2.Put(conn)
	
	dc2 := widget.DigitalClock  {
		Top : 100,
		Left : 300,
		Width : 200,
		Timezone : "America/New_York",
		Parent : "body",
	}

	dc2.Init(nil, "", "", "dc_102")
	dc2.Put(conn)


	lbl3 := widget.Label  {
		Top : 50,
		Left : 550,
		Width : 200,
		Parent : "body",
		Content : "UTC",
	}
	lbl3.Init(nil, "", "", "lbl_103")
	lbl3.Put(conn)

	dc3 := widget.DigitalClock  {
		Top : 100,
		Left : 550,
		Width : 200,
		Timezone : "UTC",
		Parent : "body",
	}

	dc3.Init(nil, "", "", "dc_103")
	dc3.Put(conn)



/*
	select {
			case recv1 := <-chan1:
		// Code to execute if a value is received from chan1
			fmt.Println("Received from chan1:", recv1)
		case chan2 <- send2:
			// Code to execute if a value can be sent to chan2
			fmt.Println("Sent to chan2:", send2)
		default:
			// Code to execute if none of the other cases are ready
			fmt.Println("No channel operation ready")
	}
*/

	for {
		dc1.Update(conn)
		dc2.Update(conn)
		dc3.Update(conn)
	
		time.Sleep(1 * time.Second)
	}

}

