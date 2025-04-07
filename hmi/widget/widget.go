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
//	"gorm.io/gorm"
//	"github.com/glebarez/sqlite"
	"github.com/gorilla/websocket"
)


// -----------------------------------------------------------------------------
// 'widget' is the interface that all widgets must support. The interface 
// consists of very simple functions that could be considered "triggers" to
// tell the widget 'what' to do. The widget itself as to know how to do it and
// may call upon the RTDSMS or sqlite or other external systems. 
// -----------------------------------------------------------------------------


type Widget interface {


	Init(id string, parent string, top int, left int, width int, height int, zIndex int, content string, 
		 options map[string]string, styles map[string]string ) error


	// Put() sends all the DOM commands to create the widgets in the display
	Show( conn *websocket.Conn ) error	
	
	// Update() causes the widget to read all the data sources and update 
	// display elements as required
	Update( conn *websocket.Conn ) error

	ClientEvent(conn *websocket.Conn, data any) error

}



