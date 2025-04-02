
// -----------------------------------------------------------------------------
// Label widget. For the digial clock we will utilize the 'output' type.
// -----------------------------------------------------------------------------

package widget


import (
	_ "embed"
	"fmt"
	"rtap/hmi/domterm"
	"time"
//	"gorm.io/gorm"
	"github.com/yuin/gopher-lua"
	"github.com/gorilla/websocket"

)


type Label struct {
	Global		domterm.GlobalAttributes
	DisplayID	string 
	Parent		string
	Content		string
	Top			int
	Left		int
	Height		int
	Width		int
	ZIndex		int

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
const luaLabelTypeName = "label"

// -----------------------------------------------------------------------------
// Registers my person type to given L.
// -----------------------------------------------------------------------------
func RegisterLabelType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaLabelTypeName)
	L.SetGlobal("label", mt)
	// static attributes
	 //   L.SetField(mt, "new", L.NewFunction(newDisplay))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), labelMethods))
}

// -----------------------------------------------------------------------------
// labelMethods table
// -----------------------------------------------------------------------------
var labelMethods = map[string]lua.LGFunction{
 //   "newLabel": luaNewLabel,
 //   "show" : luaShow,
}





func (lbl * Label) Init(display_id string, parent string, top int, left int, width int, height int, zIndex int, content string ) error {

		// TODO: Check params esp id and parent
		lbl.DisplayID = display_id
		lbl.Parent = parent
		lbl.Top = top
		lbl.Left = left
		lbl.Width = width
		lbl.Height = height
		lbl.ZIndex = zIndex
		lbl.Content = content
		
		return nil
}
	


func (lbl * Label) Show(conn *websocket.Conn){

	// Append the basic element
	attributes := make(map[string]string)
	attributes["tag"] 	= "output"
	attributes["id"] 	= lbl.DisplayID
	attributes["style"] = "position: absolute;"
	domterm.AppendElement(conn, "body","input", attributes)

	clear(attributes)

	// set the styles.
	attributes["text-align"]	= "center"
	attributes["font"]			= "Consolas"
	attributes["font-weight"]	= "bold"
	attributes["font-size"]		= "24px"

	attributes["top"]		= fmt.Sprintf("%dpx", lbl.Top)
	attributes["left"]		= fmt.Sprintf("%dpx", lbl.Left)
	if lbl.Height != 0 {
		attributes["height"]= fmt.Sprintf("%dpx", lbl.Height)
	}
	if lbl.Width != 0 {
		attributes["width"]	= fmt.Sprintf("%dpx", lbl.Width)
	}
	domterm.SetStyle(conn, lbl.DisplayID, attributes)

	domterm.SetValue(conn, lbl.DisplayID, lbl.Content)


}


// -----------------------------------------------------------------------------
// Update() does nothing since it is not tied to a real time value
// -----------------------------------------------------------------------------
func (lbl * Label) Update(conn *websocket.Conn) error{
	return nil
}


func (lbl * Label) ClientEvent(data any) {

	fmt.Printf("received label event %v\n", data)	
}

