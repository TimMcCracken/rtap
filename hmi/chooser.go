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
//	"rtap/hmi/widget"
)


// -----------------------------------------------------------------------------
// testPageHandler()
// -----------------------------------------------------------------------------
func ChooserHandler(w http.ResponseWriter, r *http.Request) {

	
	// authenticate

	// Get the realm, if any. make smart enought to default if only 1

	// Get the domain, if any. make smart enought tp default if only 1

	// is the user authorized?

	// Do we need to determine the type of display?

	// Show the display



	fmt.Printf("Starting chooser handler.")
	
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

/*
	lbl1 := widget.Label  {
		Top : 50,
		Left : 50,
		Width : 200,
		Parent : "body",
		Content : "Realm:",
	}
	lbl1.Init("lbl_101")
	lbl1.Put(conn)


	lbl2 := widget.Label  {
		Top : 50,
		Left : 300,
		Width : 200,
		Parent : "body",
		Content : "Domain",
	}
	lbl2.Init("lbl_102")
	lbl2.Put(conn)
*/

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
	
		time.Sleep(1 * time.Second)
	}


}

