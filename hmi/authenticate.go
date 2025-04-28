package hmi

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"
	dt "rtap/hmi/domterm"
)


// -----------------------------------------------------------------------------
// testPageHandler()
// -----------------------------------------------------------------------------
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {






	fmt.Printf("**************** Starting testpage handler. *************\n")
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	


	// set the page style
	attributes := make(map[string]string)
	attributes["background-color"] = "White"
	attributes["color"] = "Black"
	dt.SetStyle( conn, "html", attributes)
	clear(attributes)


	dt.SetDocumentTitle( conn, "Login to the Utelety RTAP Server")
/*
	lbl1 := widget.Label  {
		Top : 50,
		Left : 50,
		Width : 200,
		Parent : "body",
		Content : "Username",
	}
	lbl1.Init("lbl_101")
	lbl1.Put(conn)

	

	lbl2 := widget.Label  {
		Top : 50,
		Left : 300,
		Width : 200,
		Parent : "body",
		Content : "Password:",
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
//		dc1.Update(conn)
//		dc2.Update(conn)
//		dc3.Update(conn)
	
		time.Sleep(1 * time.Second)
	}

}

