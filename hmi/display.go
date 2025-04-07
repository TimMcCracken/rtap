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

func NewDisplay() (Display) {

	display := Display{

		WidgetMap : make(map[string]widget.Widget),
		ClockMap  : make(map[string]widget.Widget),
	}

	return display
}



func (d * Display) NewLabel(parent string, top int, left int, width int, height int, zIndex int, content string, 
							options map[string]string, styles map[string]string) (*widget.Label, error) {

	lbl := new( widget.Label)

	// make up an ID
	displayID := fmt.Sprintf("lbl_%d", len(d.WidgetMap))

	err := lbl.Init(displayID, parent, top, left, width, height, zIndex, content, options, styles)
	if err != nil {
		return lbl, err
	}

	// add to the map
	d.WidgetMap[lbl.DisplayID] = lbl

	// return the new label and error code
	return lbl, nil
}

func (d * Display) NewDigitalClock(parent string, top int, left int, width int, height int, zIndex int, content string,  
									options map[string]string, styles map[string]string ) (*widget.DigitalClock, error) {

	dc := new( widget.DigitalClock)

	// make up an ID
	displayID := fmt.Sprintf("dc_%d", len(d.WidgetMap))

	err := dc.Init(displayID, parent, top, left, width, height, zIndex, content, options, styles)
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

