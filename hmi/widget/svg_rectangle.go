
// -----------------------------------------------------------------------------
// DigitalClock widget. For the digial clock we will utilize the 'input' type
// -----------------------------------------------------------------------------

package widget

import (
	"fmt"
	"rtap/hmi/domterm"
	"github.com/yuin/gopher-lua"
	"github.com/gorilla/websocket"
)


type SVGRectangle struct {
	Parent		string
	DisplayID	string
	X			int
	Y 			int
	Width		int
	Height 		int
	RX			int
	RY 			int
	ZIndex		int
	Fill		string
	Stroke 		string
	StrokeWidth int

	Options 	* map[string]string
	Styles 		* map[string]string


}

// *****************************************************************************
// Widget interface methods
// *****************************************************************************

func (rect * SVGRectangle) Show(conn *websocket.Conn) error {

	// Append the basic element
	attributes := make(map[string]string)
	//attributes["title"] = rectangle.Title
	attributes["id"] 		= rect.DisplayID
	attributes["x"] 		= fmt.Sprintf("%dpx", rect.X)
	attributes["y"]			= fmt.Sprintf("%dpx", rect.Y)
	attributes["width"]		= fmt.Sprintf("%dpx", rect.Width)
	attributes["height"]	= fmt.Sprintf("%dpx", rect.Height)
	attributes["fill"]		= rect.Fill

	if rect.StrokeWidth > 0 {
		attributes["stroke-width"] = fmt.Sprintf("%dpx", rect.StrokeWidth)
		attributes["stroke"] = rect.Stroke
	}

	if rect.RX > 0 {
		attributes["rx"] = fmt.Sprintf("%dpx", rect.RX)
	}
	if rect.RY > 0 {
		attributes["ry"] = fmt.Sprintf("%dpx", rect.RY)
	}

	if rect.Options != nil {
		for key, value := range *rect.Options {
			attributes[key] = value
		}
	}

	domterm.AppendElementNS(conn, rect.Parent, "http://www.w3.org/2000/svg", "rect", attributes)


	if rect.Styles != nil {
		fmt.Printf("Styles length %d\n", len(*rect.Styles))
		domterm.SetStyle(conn, rect.DisplayID, *rect.Styles)
	}

	return nil
}




func (rect * SVGRectangle)  UpdateRealtime(conn *websocket.Conn) error {
	return nil
}

func (rect * SVGRectangle) UpdateConfig( conn *websocket.Conn ) error {
	return nil
}

func (rect * SVGRectangle) ClientEvent(conn *websocket.Conn, data any) error {

	fmt.Printf("received svg rectangle event %v\n", data)	
	return nil

}


// *****************************************************************************
// Lua support stuff
// *****************************************************************************
const luaSVGrectangleTypeName = "svgrectangle"

// -----------------------------------------------------------------------------
// Registers my person type to given L.
// -----------------------------------------------------------------------------
func RegisterSVGRectangleType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaSVGrectangleTypeName)
	L.SetGlobal(luaSVGrectangleTypeName, mt)
	// static attributes
	 //   L.SetField(mt, "new", L.NewFunction(newDisplay))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), svgRectangleMethods))
}

// -----------------------------------------------------------------------------
// rectangleMethods table
// -----------------------------------------------------------------------------
var svgRectangleMethods = map[string]lua.LGFunction{
 //   "newLabel": luaNewLabel,
 //   "show" : luaShow,
}
