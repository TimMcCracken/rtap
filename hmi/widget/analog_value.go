
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
	"github.com/yuin/gopher-lua"

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


type analogRegion struct {
	Name	string
	Color	string
	Enabled bool
}

type AnalogRegions struct {
	Region 	[9] analogRegion
}


type AnalogValue struct {
	DisplayID	string
	Req 		required // parameters required for all widgets

	Content		string
	lastValue	string 

	// configuration values
	Format		string
	Regions * 	AnalogRegions

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





// -----------------------------------------------------------------------------
// Lua support stuff
// -----------------------------------------------------------------------------
const luaAnalogValueTypeName = "analogValue"

// -----------------------------------------------------------------------------
// Registers my person type to given L.
// -----------------------------------------------------------------------------
func RegisterAnalogValueType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaDigitalClockTypeName)
	L.SetGlobal("analogValue", mt)
	// static attributes
	 //   L.SetField(mt, "new", L.NewFunction(newDisplay))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), analogValueMethods))
}

// -----------------------------------------------------------------------------
// DigitalClockMethods table
// -----------------------------------------------------------------------------
var analogValueMethods = map[string]lua.LGFunction{
 //   "newLabel": luaNewLabel,
 //   "show" : luaShow,
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
	attributes["top"]		= fmt.Sprintf("%dpx", av.Req.top)
	attributes["left"]		= fmt.Sprintf("%dpx", av.Req.left)

	if av.Req.height != 0 {
		attributes["height"]	= fmt.Sprintf("%dpx", av.Req.height)
	}
	if av.Req.width != 0 {
		attributes["width"]		= fmt.Sprintf("%dpx", av.Req.width)
	}
	domterm.SetStyle(conn, av.DisplayID, attributes)
	domterm.SetValue(conn, av.DisplayID, "-.-")

	return nil
}


func (av *AnalogValue)UpdateRealtime( conn *websocket.Conn ) error {
	return nil
}



func (av *AnalogValue)UpdateConfig( conn *websocket.Conn ) error {
	return nil
}


 

func (av *AnalogValue) ClientEvent(conn *websocket.Conn, data any) error {

	fmt.Printf("received clock event %v\n", data)	

	return nil
}
