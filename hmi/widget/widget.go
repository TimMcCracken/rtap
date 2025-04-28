/*******************************************************************************
widget.go


Widgets are the standard display objects for the RTAP system. Widgets mangage
all the low level DOM functions through domterm. In OO languages, widgets would
be the base class for various derived clasees. In GO, widget provides a common
interface and some shared helper function. 

Widgets support the following "derived" types;

	label - generally used for static infomration on a display
	digital clock - updates the date and/or time once each second
	analog value - real time
	status value - real time
	counter value - real time
	string value - real time
	
Once initialize widget generate DOM TERMINAL commands to show and update the
display.


Rev By  Date     Description
--- --- -------- --------------------------------------------------------------
002 TFM 20250301 Orinal develop started.

*******************************************************************************/


package widget


import (
	"github.com/gorilla/websocket"
)


// -----------------------------------------------------------------------------
// 'widget' is the interface that all widgets must support. The interface 
// consists of very simple functions that could be considered "triggers" to
// tell the widget 'what' to do. The widget itself as to know how to do it and
// may call upon the RTDSMS or sqlite or other external systems. 
// -----------------------------------------------------------------------------

type Widget interface {

	// Show() sends all the DOM commands to create the widgets in the display
	Show( conn *websocket.Conn ) error	
	
	// Update() causes the widget to read all the data sources and update 
	// display elements as required
	UpdateRealtime( conn *websocket.Conn ) error
	UpdateConfig( conn *websocket.Conn ) error

	ClientEvent(conn *websocket.Conn, data any) error

}

// ----------------------------------------------------------------------------
// The following structure includes the variables that are required for all
// widgets and are set wihthe display.New[widget] factory method.
// ----------------------------------------------------------------------------
type required struct {
	parent		string
	content		string
	top			int
	left		int
	height		int
	width		int
	zIndex		int
}


func( req * required ) Parse( parent string, top int, left int, width int, height int, 
	zIndex int, content string,) error {

		// TODO: Check params esp id and parent
		//	lbl.DisplayID = display_id
		req.parent = parent
		req.top = top
		req.left = left
		req.width = width
		req.height = height
		req.zIndex = zIndex
		req.content = content

		return nil
}