package hmi

import (
	_ "embed"
	"fmt"
	"net/http"
	ws "github.com/gorilla/websocket"
	"rtap/hmi/widget"
)



// Begin test code

//go:embed display_test.lua
var  DisplayTest string


var testRegions widget.AnalogRegions 



func init() {
	testRegions.Region[0].Name 	= "Low Reasonability"
	testRegions.Region[0].Color = "Purple"

	testRegions.Region[1].Name 	= "Low Critical"
	testRegions.Region[1].Color = "Red"

	testRegions.Region[2].Name 	= "Low Operating"
	testRegions.Region[2].Color = "Orange"

	testRegions.Region[3].Name 	= "Low Warning"
	testRegions.Region[3].Color = "Yellow"

	testRegions.Region[4].Name 	= "Normal"
	testRegions.Region[4].Color = "White"

	testRegions.Region[5].Name 	= "High Warning"
	testRegions.Region[5].Color = "Yellow"

	testRegions.Region[6].Name 	= "High Operating"
	testRegions.Region[6].Color = "Orange"

	testRegions.Region[7].Name 	= "High Critical"
	testRegions.Region[7].Color = "Red"

	testRegions.Region[8].Name 	= "High Reasonability"
	testRegions.Region[8].Color = "Purple"
}

// end test code





type Display struct {
	// Do these need to be capitalized?
	Widgets				[]widget.Widget
	WidgetMap 			map[string]int
	ClockMap  			map[string]int
	AnalogRegionsMap 	map[string]*widget.AnalogRegions
}



// -----------------------------------------------------------------------------
// websocket helperfunction
// -----------------------------------------------------------------------------
var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

// -----------------------------------------------------------------------------
// NewDisplay
// -----------------------------------------------------------------------------
func NewDisplay() (Display) {

	display := Display{
		Widgets				: make([]widget.Widget, 0, 100),
		WidgetMap 			: make(map[string]int),
		ClockMap  			: make(map[string]int),
		AnalogRegionsMap  	: make(map[string]*widget.AnalogRegions),
	}

	return display
}



// -----------------------------------------------------------------------------
// NewLabel()
// -----------------------------------------------------------------------------
func (d * Display) NewLabel(parent string, top int, left int, width int, height int, zIndex int, content string, 
							options map[string]string, styles map[string]string) (*widget.Label, error) {

	// -------------------------------------------------------------------------
	// Create the new object
	// -------------------------------------------------------------------------
	lbl := new( widget.Label)

	// -------------------------------------------------------------------------
	// Create a unique ID for the object that will be used as the javascript ID
	// -------------------------------------------------------------------------
	displayID := fmt.Sprintf("lbl_%d", len(d.WidgetMap))
	lbl.DisplayID = displayID

	// -------------------------------------------------------------------------
	// Parse the required parameters
	// -------------------------------------------------------------------------
	err := lbl.Req.Parse(parent, top, left, width, height, zIndex, content)
		if err != nil {
			return nil, err
		}

	// -------------------------------------------------------------------------
	// add the object to with display's widgets slice and the widget map
	// -------------------------------------------------------------------------
	d.WidgetMap[lbl.DisplayID] = len(d.Widgets)
	d.Widgets = append(d.Widgets, lbl)

	// -------------------------------------------------------------------------
	// return the new label and error code
	// -------------------------------------------------------------------------
	return lbl, nil
}

// -----------------------------------------------------------------------------
// NewDigitalClock()
// -----------------------------------------------------------------------------
func (d * Display) NewDigitalClock(parent string, top int, left int, width int, height int, zIndex int, content string,  
									options map[string]string, styles map[string]string ) (*widget.DigitalClock, error) {

	// -------------------------------------------------------------------------
	// Create the new object
	// -------------------------------------------------------------------------
	dc := new( widget.DigitalClock)

	// -------------------------------------------------------------------------
	// Create a unique ID for the object that will be used as the javascript ID
	// -------------------------------------------------------------------------
	displayID := fmt.Sprintf("dc_%d", len(d.WidgetMap))
	dc.DisplayID = displayID

	// -------------------------------------------------------------------------
	// Parse the required parameters
	// -------------------------------------------------------------------------
	err := dc.Req.Parse(parent, top, left, width, height, zIndex, content)
		if err != nil {
			return nil, err
		}

	// -------------------------------------------------------------------------
	// add the object to with display's widget map.
	// Clocks also have to be added to a clock map for now
	// -------------------------------------------------------------------------
	d.WidgetMap[dc.DisplayID] = len(d.Widgets)
	d.ClockMap[dc.DisplayID] = len(d.Widgets)
	d.Widgets = append(d.Widgets, dc)

	// -------------------------------------------------------------------------
	// return the new label and error code
	// -------------------------------------------------------------------------
	return dc, nil
}



// -----------------------------------------------------------------------------
// NewAnalogValue()
// -----------------------------------------------------------------------------
func (d * Display) NewAnalogValue(parent string, top int, left int, width int, height int, zIndex int, content string,  
	options map[string]string, styles map[string]string ) (*widget.AnalogValue, error) {

	// test code, unitl the database is ready
	d.AnalogRegionsMap["test"] = &testRegions
	// end test code

	// -------------------------------------------------------------------------
	// Create the new object
	// -------------------------------------------------------------------------
	av := new( widget.AnalogValue)

	// test code, unitl the database is ready
	d.AnalogRegionsMap["test"] = &testRegions
	av.Regions = d.AnalogRegionsMap["test"]
	// end test code



	// -------------------------------------------------------------------------
	// Create a unique ID for the object that will be used as the javascript ID
	// -------------------------------------------------------------------------
	displayID := fmt.Sprintf("dc_%d", len(d.WidgetMap))
	av.DisplayID = displayID

	// -------------------------------------------------------------------------
	// Parse the required parameters
	// -------------------------------------------------------------------------
	err := av.Req.Parse(parent, top, left, width, height, zIndex, content)
	if err != nil {
	return nil, err
	}

	// -------------------------------------------------------------------------
	// add the object to with display's widget map.
	// -------------------------------------------------------------------------
	d.WidgetMap[av.DisplayID] = len(d.Widgets)
	d.Widgets = append(d.Widgets, av)

	// -------------------------------------------------------------------------
	// return the new label and error code
	// -------------------------------------------------------------------------
	return av, nil
}




// -----------------------------------------------------------------------------
// Show() should be called after creating all the display objects
// -----------------------------------------------------------------------------
func (d * Display) Show(conn *ws.Conn) {

	for _, widget := range d.Widgets {
		widget.Show(conn)
	}

}

