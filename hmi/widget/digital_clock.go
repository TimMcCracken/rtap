
// -----------------------------------------------------------------------------
// DigitalClock widget. For the digial clock we will utilize the 'input' type
// -----------------------------------------------------------------------------

package widget


import (
	_ "embed"
	"fmt"
	"rtap/hmi/domterm"
	"time"
//	"gorm.io/gorm"
	"github.com/gorilla/websocket"
)


type DigitalClock struct {
	Global		domterm.GlobalAttributes
	DisplayID	string
	Parent		string
	Top			int
	Left		int
	Height		int
	Width		int
	ZIndex		int
	Format		string
	Timezone	string

	// unexported properties below
	TzLocation	*time.Location
	Title		string
	// lastValue is used to decide if we need to send again. Handy for dates
	// and time without seconds, etc.
	lastValue	string 
}

/*
func (dc *DigitalClock) Init(display_id string){

	dc.DisplayID = display_id

	//TODO: read frin database


	// Permanent code
	if dc.Timezone == "" {
		dc.Timezone = "Local"
	}

	tzLocation, err := time.LoadLocation(dc.Timezone)
	dc.tzLocation = tzLocation
	if err != nil {
		dc.Timezone = "Local"
		dc.title = fmt.Sprintf("%s: %v", dc.Timezone, err)
	}

} */


// ----------------------------------------------------------------------------
// dcoL get the following defauilt format from the database
// ----------------------------------------------------------------------------
var defaultTimeFormat string



func init() {
	defaultTimeFormat = "2006-01-02 15:04:05"
}


func (dc *DigitalClock) Init(id string, parent string, top int, left int, width int, height int, zIndex int, content string) error {

	// TODO: Check params esp id and parent

	dc.DisplayID = id
	dc.Parent = parent
	dc.Top = top
	dc.Left = left
	dc.Width = width
	dc.Height = height
	dc.ZIndex = zIndex

	if dc.Timezone == "" {
		dc.Timezone = "Local"
	}

	tzLocation, err := time.LoadLocation(dc.Timezone)
	dc.TzLocation = tzLocation
	if err != nil {
		dc.Timezone = "Local"
		dc.Title = fmt.Sprintf("%s: %v", dc.Timezone, err)
	}

	return nil
}

func (dc *DigitalClock) Show(conn *websocket.Conn){

	// Append the basic element
	attributes := make(map[string]string)
	attributes["title"] = fmt.Sprintf("Timezone: %s",dc.Timezone)
	attributes["tag"] 	= "output"
	attributes["id"] 	= dc.DisplayID
	attributes["style"] = "position: absolute;"

	attributes["onclick"] = "sendMouseEvent(event)"


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

func (dc *DigitalClock) Update(conn *websocket.Conn){

	var ft string

	// get the current time for the specified location
	t := time.Now().In(dc.TzLocation)

	// format the time as speified, or using the default.
	if dc.Format == "" {
		ft = t.Format(defaultTimeFormat)
	} else {
		ft = t.Format(dc.Format)
	}

	// If there is no change from the 'lastValue' sent, then don't send again.
	if ft == dc.lastValue {
		return
	}

	// Update the value in the browser and copy to the 'lastValue'
	domterm.SetValue(conn, dc.DisplayID, ft)
	dc.lastValue = ft

}

