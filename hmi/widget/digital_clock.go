
// -----------------------------------------------------------------------------
// DigitalClock widget. For the digial clock we will utilize the 'input' type
// -----------------------------------------------------------------------------

package widget





/*
[[--
	dc1 = display:NewDigitalClock("body", 100, 50, 200, 0, 1, "Local" )
	
    dc2 = display:NewDigitalClock("body", 100, 300, 200, 0, 1, "New York" ) 
	
    dc3 = display:NewDigitalClock("body", 100, 550, 200, 0, 1, "UTC" )  ]]

*/
import (
	_ "embed"
	"fmt"
	"rtap/hmi/domterm"
	"time"
//	"gorm.io/gorm"
	"github.com/yuin/gopher-lua"
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
	tzLocation	*time.Location
	title		string
	// lastValue is used to decide if we need to send again. Handy for dates
	// and time without seconds, etc.
	lastValue	string 
}

// -----------------------------------------------------------------------------
// Lua support stuff
// -----------------------------------------------------------------------------
const luaDigitalClockTypeName = "digitalClock"

// -----------------------------------------------------------------------------
// Registers my person type to given L.
// -----------------------------------------------------------------------------
func RegisterDigitalClockType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaDigitalClockTypeName)
	L.SetGlobal("digitalClock", mt)
	// static attributes
	 //   L.SetField(mt, "new", L.NewFunction(newDisplay))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), digitalClockMethods))
}

// -----------------------------------------------------------------------------
// DigitalClockMethods table
// -----------------------------------------------------------------------------
var digitalClockMethods = map[string]lua.LGFunction{
 //   "newLabel": luaNewLabel,
 //   "show" : luaShow,
}


// ----------------------------------------------------------------------------
// dcoL get the following defauilt format from the database
// ----------------------------------------------------------------------------
var defaultTimeFormat string



func init() {
	defaultTimeFormat = "2006-01-02 15:04:05"
}


func (dc * DigitalClock) Init(id string, parent string, top int, left int, width int, height int, 
							zIndex int, content string, 
							options map[string]string, styles map[string]string ) error {

	// TODO: Check params esp id and parent

	dc.DisplayID = id
	dc.Parent = parent
	dc.Top = top
	dc.Left = left
	dc.Width = width
	dc.Height = height
	dc.ZIndex = zIndex
	dc.Timezone = content

	if dc.Timezone == "" {
		dc.Timezone = "Local"
	}


	tzLocation, err := time.LoadLocation(dc.Timezone)
	if err != nil {
		dc.Timezone = "ERROR"
		// todo: add an error message to title
		dc.title = fmt.Sprintf("Timezone error: %v",  err)
		return nil
	}
	dc.title = fmt.Sprintf("Timezone: %s",dc.Timezone)
	dc.tzLocation = tzLocation
	return nil
}

func (dc * DigitalClock) Show(conn *websocket.Conn) error {

	// Append the basic element
	attributes := make(map[string]string)
	attributes["title"] = dc.title
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

	if dc.Timezone == "ERROR" {
		domterm.SetValue(conn, dc.DisplayID, "ERROR!")
	}

	return nil
}



func (dc * DigitalClock) Update(conn *websocket.Conn) error {

	var ft string


	// If the timezone could not be processed, then exit
	if dc.Timezone == "ERROR" {
		return nil  // THIS SHOULD PROBABLY BE AN ERROR, but we dont want to break later processing
	} 


	// get the current time for the specified location

	t := time.Now().In(dc.tzLocation)
//	t := time.Now()

	// format the time as speified, or using the default.
	if dc.Format == "" {
		ft = t.Format(defaultTimeFormat)
	} else {
		ft = t.Format(dc.Format)
	}

	// If there is no change from the 'lastValue' sent, then don't send again.
	if ft == dc.lastValue {
		return nil
	}

	// Update the value in the browser and copy to the 'lastValue'
	err := domterm.SetValue(conn, dc.DisplayID, ft)
	if err != nil {
		return err
	}

	dc.lastValue = ft

	return nil
}


func (dc * DigitalClock) ClientEvent(conn *websocket.Conn, data any) error {

	fmt.Printf("received clock event %v\n", data)	
	return nil

}
