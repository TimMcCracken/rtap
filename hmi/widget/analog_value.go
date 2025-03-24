
// -----------------------------------------------------------------------------
// AnalogValue widget. For the digial clock we will utilize the 'input' type
// -----------------------------------------------------------------------------

package widget


import (
	_ "embed"
	"fmt"
	"rtap/hmi/domterm"
//	"time"
	"gorm.io/gorm"
	"github.com/gorilla/websocket"
)


var	AnalogUpdateMsg struct {
	value			float64
	valueColor		string
	valueBackground	string
	stateName		string
	stateColor		string
	stateBackground string
	unackAlarm		bool	// 'true' sets flashing mode
}


type AnalogValue struct {
	Global		domterm.GlobalAttributes
	DisplayID	string
	Parent		string
	Top			int
	Left		int
	Height		int
	Width		int
	ZIndex		int
	Format		string

	lastValue	string 

	// configuration values
	valueFormat		string
	valueEU			string

	// The following values are local storage for information sent in the update
	// msg from DACC

	value			float64
	valueColor		string
	valueBackground	string
	stateName		string
	stateColor		string
	stateBackground string
	unackAlarm		bool	// 'true' sets flashing mode

}

func (anlg *AnalogValue) Init(db *gorm.DB, descriptor string, rtdata string, display_id string){

	anlg.DisplayID = display_id

	//TODO: read from display database file 


	// Permanent code
//	if dc.Timezone == "" {
//		dc.Timezone = "Local"
//	}

//	tzLocation, err := time.LoadLocation(dc.Timezone)
//	dc.tzLocation = tzLocation
//	if err != nil {
//		dc.Timezone = "Local"
//		dc.title = fmt.Sprintf("%s: %v", dc.Timezone, err)
//	}
}


// ----------------------------------------------------------------------------
// 
// ----------------------------------------------------------------------------

func init() {
	
}



func (dc *AnalogValue) Put(conn *websocket.Conn){

	// Append the basic element
	attributes := make(map[string]string)
	attributes["title"] = fmt.Sprintf("Analog Value:")
	attributes["tag"] 	= "output"
	attributes["id"] 	= dc.DisplayID
	attributes["style"] = "position: absolute;"
	domterm.AppendElement(conn, "body","input", attributes)
	clear(attributes)

	// set the styles.
	attributes["text-align"]	= "center"
	attributes["font"]			= "Consolas"
	attributes["font-weight"]	= "bold"
	attributes["font-size"]		= "32px"
	attributes["top"]		= fmt.Sprintf("%dpx", dc.Top)
	attributes["left"]		= fmt.Sprintf("%dpx", dc.Left)

	if dc.Height != 0 {
		attributes["height"]	= fmt.Sprintf("%dpx", dc.Height)
	}
	if dc.Width != 0 {
		attributes["width"]		= fmt.Sprintf("%dpx", dc.Width)
	}
	domterm.SetStyle(conn, dc.DisplayID, attributes)

}

func (dc *AnalogValue) Update(conn *websocket.Conn){

//	var ft string

	// get the current time for the specified location
//	t := time.Now().In(dc.tzLocation)

	// format the time as speified, or using the default.
//	if dc.Format == "" {
//		ft = t.Format(defaultTimeFormat)
//	} else {
//		ft = t.Format(dc.Format)
//	}

	// If there is no change from the 'lastValue' sent, then don't send again.
//	if ft == dc.lastValue {
//		return
//	}

	// Update the value in the browser and copy to the 'lastValue'
//	domterm.SetValue(conn, dc.DisplayID, ft)
//	dc.lastValue = ft

}
 
