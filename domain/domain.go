// -----------------------------------------------------------------------------
// display.go
//
// This file contatins the main display handler worker. IAs each display can
// only display data from a single domain, it operates as a task at the domain
// level.
// All displays execute a Utelety Command Language (UCL) based script, whihc is
// stored in the display (sqlite) file. Displays can optionally have declared
// elements as well, which will be loaded and created prior to calling the UCL
// code.
//
// Evenutally displays will also support the 'D3' vosualization libraries.  
//
// Note that this module only contains the high level display functionality. 
// The primiteives are all contained in the HMI package.
//
// Rev By  Date     Description
// --- --- -------- ------------------------------------------------------------
//
//
//
//
// -----------------------------------------------------------------------------



package domain

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"log"
	"time"
	"encoding/json"

	"github.com/glebarez/sqlite"
//	"github.com/google/uuid"

	"github.com/yuin/gopher-lua"
	ws "github.com/gorilla/websocket"
	

	bp "rtap/buffer_pool"
	"rtap/hmi/domterm"
	"rtap/dacc"
	"rtap/hmi"
	"rtap/metronome"
	mq "rtap/message_q"
	"rtap/rtdsms"
	"rtap/hmi/widget"
)



type Domain struct {
	Descriptor		DomainDescriptor
	messageQueue	mq.MessageQ
	Datastores		[]*rtdsms.Datastore
	Datastores_map 	map[string]int

	dacc			dacc.DACC
	metronome		metronome.Metronome
	bufferPool		bp.BufferPool


	// Variables used by the JMI subsystem
	handlerMessageQueue	mq.MessageQ
	mqChannel			chan mq.Message	 			// for messages rec'd from the message_q 
	handlerChannel		chan mq.Message				// for messages rec'd from displayHandler workers 
//	handlerTxChannels	map[uuid.UUID]mq.Message	// for messages to transmit to displayHandler workers
}


// TODO: Try to eliminate this code later. It is used by REALM
func (domain * Domain)  MessageQueue() (* mq.MessageQ) {

	return &domain.messageQueue
}



func (domain * Domain) Start() {
	
	domain.bufferPool.Start()
	domain.messageQueue.Start(&domain.bufferPool)
	domain.handlerMessageQueue.Start(&domain.bufferPool)

	domain.metronome.Start(&domain.bufferPool, &domain.messageQueue)
 
	// temporarily not started
	//	domain.dacc.Start(&domain.bufferPool, &domain.messageQueue)

	go domain.HMITask() 

}



// -----------------------------------------------------------------------------
// construct   ()  TODO: I had a different name for this, but it escapes me.
// try to remember waht it was
// -----------------------------------------------------------------------------
func (domain * Domain) Construct(filename string) error  {

	//	fmt.Printf("Constructing domain [%s]\n", domain.descriptor.domain_name)				

	if domain.Datastores_map == nil {
		domain.Datastores_map = make(map[string]int)
	}

	// Open the sqlite file
	db, err := gorm.Open(sqlite.Open(filename + "?_pragma=foreign_keys(1)" ), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database using path [%s]", filename)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Internal error: failed to obtain database.DB in CreateDomain()")
	}
	defer sqlDB.Close()


	// ------------------------------------------------------------------------
	// Get the domain_id
	// ------------------------------------------------------------------------
	var domain_id	int64
	result := db.Raw("SELECT domain_id FROM domains where domain_name = ?", domain.Descriptor.DomainName).Scan(&domain_id)
	if result.Error != nil {
		return result.Error
	}


	// ------------------------------------------------------------------------
	// Select the rows
	// ------------------------------------------------------------------------
	rows, err := sqlDB.Query("SELECT datastore_name FROM datastores where domain_id = ? order by datastore_id", domain_id)
    if err != nil {
        return err
    }
    defer rows.Close()

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
		var name string
        err := rows.Scan(&name)
		if err != nil {
            return fmt.Errorf("%v\n", err)
        } else {
			ds, err := rtdsms.NewDatastore(domain.Descriptor.RealmName, 
										domain.Descriptor.DomainName, name )
			if err != nil {
				return fmt.Errorf("%v\n", err)
			}
			/*
			ds := rtdsms.Datastore {
			}
			ds.Descriptor.Datastore_name = name
			ds.Descriptor.Domain_name = domain.Descriptor.Domain_name
			ds.Descriptor.Realm_name = domain.Descriptor.Realm_name
			*/
			domain.Datastores_map[name] = len(domain.Datastores_map)
			domain.Datastores = append(domain.Datastores, ds)
		}
    }
    if err = rows.Err(); err != nil {
		return fmt.Errorf("%v\n", err)
    }

	for _, ds := range domain.Datastores {
		err = ds.Construct(filename)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}

	return nil
}



// -----------------------------------------------------------------------------
// HMILoop listens for messages from the MessageQ as well as from the HMIWorker
// goroutines that are spawned by HMI loop. When it receives a message from the
// message_q, it sends it to the worker go routines. When it receives a msg from
// a worker, it .....
// -----------------------------------------------------------------------------
func (domain * Domain) HMITask() {


	fmt.Printf("Starting HMI Loop [%s] [%s]\n", domain.Descriptor.RealmName, domain.Descriptor.DomainName)

	// TODO FIX THIS ERROR HANDLER

	// -------------------------------------------------------------------------
	// Register with the messageQ
	// -------------------------------------------------------------------------
	ch, err := domain.messageQueue.Register("HMI") 
	if err != nil {
		fmt.Printf("ERROR! Invalid object name [HMI]: %v\n", err)
		return
	}
	domain.mqChannel = ch
	
	// -------------------------------------------------------------------------
	// Create a channel for displayWorkers to send messages
	// -------------------------------------------------------------------------
	domain.handlerChannel = make(chan mq.Message) 
	if err != nil {
		fmt.Printf("ERROR! Invalid object name [HMI]: %v\n", err)
		return
	}
	
	// -------------------------------------------------------------------------
	// Listen for messages from both the message queue and all the displays that
	// are active on this domain.
	// -------------------------------------------------------------------------
	for {

		select {
		
			// For now, and message received on the HMI channel will be sent to
			// all the displayHandler workers. Later, this will be upgraded to
			// use the publish/subscribe model
			case mqMsg := <- domain.mqChannel:
//				fmt.Println("Received MQ Msg:", mqMsg)
				
				// Change the destinations to *
				mqMsg.Destinations = mqMsg.Destinations[:0]
				mqMsg.Destinations = append(mqMsg.Destinations, "*")
				
				// send the message on the handler Queue
				domain.handlerMessageQueue.Send(mqMsg.Destinations, mqMsg.Data)
				
				// return the buffer to the pool
				domain.bufferPool.Put(mqMsg.Data)
			
				// For now, just silently discard everything. This is where 
				// our publish.subscribe messages and operator actions will
				// come through
			case handlerMsg := <- domain.handlerChannel:
				fmt.Println("Received Handler MSG:", handlerMsg)
				domain.bufferPool.Put(handlerMsg.Data)
		}
	}
}




// -----------------------------------------------------------------------------
// upgrader is used by the displayHandler method
// -----------------------------------------------------------------------------

var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}


// -----------------------------------------------------------------------------
// Domain.DisplayHandler() is created for each 'display' request. A display may
// be a complete display, but may also be a component used int he desktop.
// -----------------------------------------------------------------------------
func (domain * Domain)DisplayHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Starting display handler.")

	// todo: authenticate the connection and get the username associated 
	// with the connection.



	// -------------------------------------------------------------------------
	// Upgrade the connection to a wenspcket
	// -------------------------------------------------------------------------
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()



	// -------------------------------------------------------------------------
	// Create a display object
	// -------------------------------------------------------------------------
	display := hmi.NewDisplay()
	

	// -------------------------------------------------------------------------
	// Create a lua state
	// -------------------------------------------------------------------------
	L := lua.NewState()
	defer L.Close()

	hmi.RegisterDisplayType(L)
	widget.RegisterLabelType(L)
	widget.RegisterDigitalClockType(L)
	widget.RegisterAnalogValueType(L)

	// -------------------------------------------------------------------------
	// Add the display to Lua as a global variable
	// -------------------------------------------------------------------------
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
	if err := L.DoString(hmi.DisplayTest); err != nil {

		// TODO Log an error
		eMsg := fmt.Sprintf("Lua parsing error: %v", err)
		fmt.Println(eMsg)

		// send the message to the browser, and quit!
		domterm.ShowResponseError(conn, 500, eMsg)
		return
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
			eMsg := fmt.Sprintf("Lua error in main(): %s", err)
			fmt.Println(eMsg)
			// send the message to the browser, and quit!
			domterm.ShowResponseError(conn, 500, eMsg)
			return
		}
	} else {
		eMsg := fmt.Sprintf("Error: 'main()' function not found in Lua state.")
		fmt.Println(eMsg)
		// send the message to the browser, and quit!
		domterm.ShowResponseError(conn, 500, eMsg)
		return
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
	// The following is the loop that runs until the connection is terminated.
	// This code responds to all the events coming in.
	// -------------------------------------------------------------------------
	for {
		select {

			// ----------------------------------------------------------------
			// These will be events sent from RTAP via the messageQ and HMITask
			// ----------------------------------------------------------------
			case _ = <- tickChan:
				fmt.Printf("received tick %v\n", time.Now().Unix())

				for _, ndx := range display.ClockMap {
					err = dc.Widgets[ndx].UpdateRealtime(conn)
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
		
				widget, ok := display.WidgetMap[target_id.(string)]
				if ok == false {
					fmt.Printf("Error parsing target id\n")
				} else {
					widget.ClientEvent(conn, data)
				}
			}
		}			
	}
}
