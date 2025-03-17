/*******************************************************************************
widget.go


Widgets are the standard display objects for the RTAP system. Widgets mangage
all the low levll DOM functions through domterm. Widgets support the floowing
types which may be expanded in the future.


*******************************************************************************/


package widget


import (
	_ "embed"
//	"encoding/json"
//	"flag"
//	"fmt"
//	"html/template"
//	"log"
//	"maps"
//	"net/http"
//	"time"
//	"rtap/hmi/domterm"
//	"strconv"
//	"time"
	"gorm.io/gorm"
//	"github.com/glebarez/sqlite"
	"github.com/gorilla/websocket"
)


// -----------------------------------------------------------------------------
// 'widget' is the interface that all widgets must support. The interface 
// consists of very simple functions that could be considered "triggers" to
// tell the widget 'what' to do. The widget itself as to know how to do it and
// may call upon the RTDSMS or sqlite or other external systems. 
// -----------------------------------------------------------------------------


type widget interface {

	// Init() initializes all the data in the widget. 
	// conn is a connection to a sqlite file that contains the descriptor(s)
	// descriptor is the name (key) of the descriptor(s)
	// rtdata is used to locate the data in the database
	Init(db *gorm.DB, descriptor string, rtdata string, display_id string)	


	// Put() sends all the DOM commands to create the widgets in the display
	Put( conn *websocket.Conn )	
	
	// Update() causes the widget to read all the data sources and update 
	// display elements as required
	Update( conn *websocket.Conn )

}



