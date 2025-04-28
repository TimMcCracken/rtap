
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


type SVGCircle struct {
	Parent		string
	DisplayID	string
	CenterX		int
	CenterY 	int
	Radius 		int
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

func (circle * SVGCircle) Show(conn *websocket.Conn) error {

	// Append the basic element
	attributes := make(map[string]string)
	//attributes["title"] = circle.Title
	attributes["id"] 	= circle.DisplayID
	attributes["cx"]	= fmt.Sprintf("%dpx", circle.CenterX)
	attributes["cy"]	= fmt.Sprintf("%dpx", circle.CenterY)
	attributes["r"]		= fmt.Sprintf("%dpx", circle.Radius)
	attributes["fill"]	= circle.Fill

	if circle.StrokeWidth > 0 {
		attributes["stroke-width"] = fmt.Sprintf("%dpx", circle.StrokeWidth)
		attributes["stroke"] = circle.Stroke
	}

	domterm.AppendElementNS(conn, circle.Parent, "http://www.w3.org/2000/svg", "circle", attributes)
	return nil
}




func (circle * SVGCircle)  UpdateRealtime(conn *websocket.Conn) error {
	return nil
}

func (circle * SVGCircle) UpdateConfig( conn *websocket.Conn ) error {
	return nil
}

func (circle * SVGCircle) ClientEvent(conn *websocket.Conn, data any) error {

	fmt.Printf("received svg Circle event %v\n", data)	
	return nil

}


// *****************************************************************************
// Lua support stuff
// *****************************************************************************
const luaSVGCircleTypeName = "svgCircle"

// -----------------------------------------------------------------------------
// Registers my person type to given L.
// -----------------------------------------------------------------------------
func RegisterSVGCircleType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaSVGCircleTypeName)
	L.SetGlobal("svgCircle", mt)
	// static attributes
	 //   L.SetField(mt, "new", L.NewFunction(newDisplay))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), svgCircleMethods))
}

// -----------------------------------------------------------------------------
// CircleMethods table
// -----------------------------------------------------------------------------
var svgCircleMethods = map[string]lua.LGFunction{
 //   "newLabel": luaNewLabel,
 //   "show" : luaShow,
}
