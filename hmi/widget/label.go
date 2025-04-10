
// -----------------------------------------------------------------------------
// Label widget. For the digial clock we will utilize the 'output' type.
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


type Label struct {
	DisplayID	string 
	Req 		required // parameters required for all widgets
	Title		string
	// lastValue is used to decide if we need to send again. Handy for dates
	// and time without seconds, etc.
	lastValue	string 
}



// -----------------------------------------------------------------------------
// Lua support stuff
// -----------------------------------------------------------------------------
const luaLabelTypeName = "label"

// -----------------------------------------------------------------------------
// Registers 'label' type to a given L.
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




/*
func (lbl * Label) Init(display_id string, parent string, top int, left int, width int, height int, 
						zIndex int, content string,
						options map[string]string, styles map[string]string  ) error {

/*
		err := lbl.req.parse(parent, top, left, width, height, zIndex, content)
		if err != nil {
			return err
		}
*/
		// TODO: Check params esp id and parent
		/*
		lbl.DisplayID = display_id
		lbl.req.parent = parent
		lbl.req.top = top
		lbl.req.left = left
		lbl.req.width = width
		lbl.req.height = height
		lbl.req.zIndex = zIndex
		lbl.req.content = content
		*/
/*
		return nil
}*/
	


func (lbl * Label) Show(conn *websocket.Conn) error {

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

	attributes["top"]		= fmt.Sprintf("%dpx", lbl.Req.top)
	attributes["left"]		= fmt.Sprintf("%dpx", lbl.Req.left)
	if lbl.Req.height != 0 {
		attributes["height"]= fmt.Sprintf("%dpx", lbl.Req.height)
	}
	if lbl.Req.width != 0 {
		attributes["width"]	= fmt.Sprintf("%dpx", lbl.Req.width)
	}
	domterm.SetStyle(conn, lbl.DisplayID, attributes)

	domterm.SetValue(conn, lbl.DisplayID, lbl.Req.content)

	return nil
}


// -----------------------------------------------------------------------------
// Update() does nothing since it is not tied to a real time value
// -----------------------------------------------------------------------------
func (lbl * Label)UpdateRealtime(conn *websocket.Conn) error{
	return nil
}


func (lbl * Label)UpdateConfig( conn *websocket.Conn ) error {
	return nil
}



func (lbl * Label) ClientEvent(conn *websocket.Conn, data any) error {

	fmt.Printf("received label event %v\n", data)	
	return nil

}

