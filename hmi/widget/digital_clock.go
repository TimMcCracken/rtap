
// -----------------------------------------------------------------------------
// DigitalClock widget. For the digial clock we will utilize the 'input' type
// -----------------------------------------------------------------------------

package widget

import (
//	_ "embed"
	"fmt"
	"rtap/hmi/domterm"
	"time"
	"github.com/yuin/gopher-lua"
	"github.com/gorilla/websocket"
)


type DigitalClock struct {
	DisplayID	string
	Req			required
	Format		string
	Title		string

	// uneported values
	timezone	string
	tzLocation	*time.Location


	// lastValue is used to decide if we need to send again. Handy for dates
	// and time without seconds, etc.
	lastValue	string 

	Options 	* map[string]string
	Styles 		* map[string]string


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
    "test" : luaTest,
}

func luaTest(L *lua.LState) int {
	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	return 0
}

// ----------------------------------------------------------------------------
// dcoL get the following defauilt format from the database
// ----------------------------------------------------------------------------
var defaultTimeFormat string



func init() {
	defaultTimeFormat = "2006-01-02 15:04:05"
}




func (dc * DigitalClock) Show(conn *websocket.Conn) error {

	if dc.Req.content == "" {
		dc.timezone = "Local"
	} else {
		dc.timezone = dc.Req.content
									
		tzLocation, err := time.LoadLocation(dc.timezone)
		if err != nil {
			dc.timezone = "ERROR"
			dc.Title = fmt.Sprintf("Timezone error: %v",  err)
			fmt.Printf("Timezone error: %v\n", err)
		} else {
			dc.Title = fmt.Sprintf("Timezone: %s",dc.timezone)
			dc.tzLocation = tzLocation
		}
	}

	// Append the basic element
	attributes := make(map[string]string)
	attributes["title"] = dc.Title
	attributes["tag"] 	= "output"
	attributes["id"] 	= dc.DisplayID
	attributes["style"] = "position: absolute;"

	attributes["onclick"] = "sendMouseEvent(event)"

	domterm.AppendElement(conn, dc.Req.parent,"input", attributes)
	clear(attributes)

	// set the styles.
	attributes["text-align"]	= "center"
	attributes["font"]			= "Consolas"
	attributes["font-weight"]	= "bold"
	attributes["font-size"]		= "32px"
	attributes["top"]		= fmt.Sprintf("%dpx", dc.Req.top)
	attributes["left"]		= fmt.Sprintf("%dpx", dc.Req.left)

	if dc.Req.height != 0 {
		attributes["height"]	= fmt.Sprintf("%dpx", dc.Req.height)
	}
	if dc.Req.width != 0 {
		attributes["width"]		= fmt.Sprintf("%dpx", dc.Req.width)
	}
	domterm.SetStyle(conn, dc.DisplayID, attributes)

	if dc.timezone == "ERROR" {
		domterm.SetValue(conn, dc.DisplayID, "ERROR!")
	}

	return nil
}




func (dc * DigitalClock)  UpdateRealtime(conn *websocket.Conn) error {

	var ft string


	// If the timezone could not be processed, then exit
	if dc.timezone == "ERROR" {
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

func (dc * DigitalClock) UpdateConfig( conn *websocket.Conn ) error {
	return nil
}



func (dc * DigitalClock) ClientEvent(conn *websocket.Conn, data any) error {

	fmt.Printf("received clock event %v\n", data)	
	return nil

}
