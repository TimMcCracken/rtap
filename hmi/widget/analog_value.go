
// -----------------------------------------------------------------------------
// AnalogValue widget. For the digial clock we will utilize the 'input' type
// -----------------------------------------------------------------------------

package widget


import (
	_ "embed"
	"fmt"
	"rtap/hmi/domterm"
//	"time"
//	"gorm.io/gorm"
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
	Content		string
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



// ----------------------------------------------------------------------------
// 
// ----------------------------------------------------------------------------


func (av *AnalogValue) Init(display_id string, parent string, top int, left int, width int, height int, 
	zIndex int, content string,
	options map[string]string, styles map[string]string  ) error {

	// TODO: Check params esp id and parent
	av.DisplayID = display_id
	av.Parent = parent
	av.Top = top
	av.Left = left
	av.Width = width
	av.Height = height
	av.ZIndex = zIndex
	av.Content = content

	return nil
}

func (av *AnalogValue) Show(conn *websocket.Conn) error {

	// Append the basic element
	attributes := make(map[string]string)
	attributes["title"] = fmt.Sprintf("Analog Value:")
	attributes["tag"] 	= "output"
	attributes["id"] 	= av.DisplayID
	attributes["style"] = "position: absolute;"

	domterm.AppendElement(conn, "body","input", attributes)
	clear(attributes)

	// set the styles.
	attributes["text-align"]	= "center"
	attributes["font"]			= "Consolas"
	attributes["font-weight"]	= "bold"
	attributes["font-size"]		= "32px"
	attributes["top"]		= fmt.Sprintf("%dpx", av.Top)
	attributes["left"]		= fmt.Sprintf("%dpx", av.Left)

	if av.Height != 0 {
		attributes["height"]	= fmt.Sprintf("%dpx", av.Height)
	}
	if av.Width != 0 {
		attributes["width"]		= fmt.Sprintf("%dpx", av.Width)
	}
	domterm.SetStyle(conn, av.DisplayID, attributes)
	domterm.SetValue(conn, av.DisplayID, "0.0")

	return nil
}

func (av *AnalogValue) Update(conn *websocket.Conn) error {

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

	return nil
}
 

func (av *AnalogValue) ClientEvent(conn *websocket.Conn, data any) error {

	fmt.Printf("received clock event %v\n", data)	

	return nil
}
