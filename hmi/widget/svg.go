
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


type SVG struct {
	DisplayID	string
	Req			required
	Format		string
	Title		string

	Widgets				[]Widget
	WidgetMap 			map[string]int

	// unexported values

	// lastValue is used to decide if we need to send again. Handy for dates
	// and time without seconds, etc.
	// lastValue	string 
}


// ----------------------------------------------------------------------------
// GetAsString tries to get a Go string from an LValue.
// Returns (string, true) if successful, ("", false) otherwise.
// ----------------------------------------------------------------------------
func GetAsString(lv lua.LValue) (string, bool) {
	switch v := lv.(type) {
	case lua.LString:
		return string(v), true
	case lua.LNumber:
		return v.String(), true
	case *lua.LNilType:
		return "", true	
	default:
		return "", false
	}
}

// ----------------------------------------------------------------------------
// helper function to convert a Lua table to a Go map[string]string
// ----------------------------------------------------------------------------
func luaTableToStringMap(tbl *lua.LTable) map[string]string {
	result := make(map[string]string)
	tbl.ForEach(func(key, value lua.LValue) {

		k, ok1 := key.(lua.LString)
		v, ok2 := GetAsString(value)

		if ok1 && ok2 {
			result[string(k)] = string(v)
		} 
	})
	return result
}





// *****************************************************************************
// Widget Interface Methods
// *****************************************************************************
func (svg * SVG) Show(conn *websocket.Conn) error {

	fmt.Printf("----------- SVG.SHOW() %d widgets ------------\n", len(svg.Widgets))

	// Append the basic element
	attributes := make(map[string]string)
	attributes["tag"] 	= "svg"
	attributes["id"] 	= svg.DisplayID
	attributes["style"] = "position: absolute;"
	attributes["onclick"] = "sendMouseEvent(event)"

	domterm.AppendElementNS(conn, svg.Req.parent, "http://www.w3.org/2000/svg", "svg", attributes)
	clear(attributes)

	attributes["top"]		= fmt.Sprintf("%dpx", svg.Req.top)
	attributes["left"]		= fmt.Sprintf("%dpx", svg.Req.left)
	attributes["height"]	= fmt.Sprintf("%dpx", svg.Req.height)
	attributes["width"]		= fmt.Sprintf("%dpx", svg.Req.width)

	domterm.SetStyle(conn, svg.DisplayID, attributes)

	for _, widget := range svg.Widgets {
		widget.Show(conn)
	}
	return nil
}

func (svg * SVG)  UpdateRealtime(conn *websocket.Conn) error {
	return nil
}

func (svg * SVG) UpdateConfig( conn *websocket.Conn ) error {
	return nil
}

func (svg * SVG) ClientEvent(conn *websocket.Conn, data any) error {

	fmt.Printf("received svg event %v\n", data)	
	return nil

}

// *****************************************************************************
// Factory methods
// *****************************************************************************

// -----------------------------------------------------------------------------
// NewSVGCircle()
// -----------------------------------------------------------------------------
func (svg * SVG) NewSVGCircle(parent string, centerX int, centerY int, 
	radius int, zIndex int, fill string, stroke string, strokeWidth int, 
	options * map[string]string, styles * map[string]string ) (*SVGCircle, error) {

	// -------------------------------------------------------------------------
	// Create the new object
	// -------------------------------------------------------------------------
	circle := new( SVGCircle)

	// -------------------------------------------------------------------------
	// Create a unique ID for the object that will be used as the javascript ID
	// -------------------------------------------------------------------------
	displayID := fmt.Sprintf("%s_circle_%d", svg.DisplayID, len(svg.WidgetMap))
	circle.DisplayID = displayID

	// -------------------------------------------------------------------------
	// Parse the required parameters
	// -------------------------------------------------------------------------
	circle.Parent	= parent
	circle.CenterX	= centerX
	circle.CenterY 	= centerY
	circle.Radius 	= radius
	circle.ZIndex	= zIndex
	circle.Fill		= fill
	circle.Stroke 	= stroke
	circle.StrokeWidth = strokeWidth

	circle.Options = options
	circle.Styles = styles

	// -------------------------------------------------------------------------
	// add the object to the svg's widget map.
	// -------------------------------------------------------------------------
	svg.WidgetMap[circle.DisplayID] = len(svg.Widgets)
	svg.Widgets = append(svg.Widgets, circle)

	// -------------------------------------------------------------------------
	// return the new circle and nil error code
	// -------------------------------------------------------------------------
	return circle, nil
}


// -----------------------------------------------------------------------------
// NewSVGRectangle()
// -----------------------------------------------------------------------------
func (svg * SVG) NewSVGRectangle(parent string, x int, y int, width int, 
	height int, zIndex int, rx int, ry int, fill string,
	stroke string, strokeWidth int, options * map[string]string, 
	styles * map[string]string ) (*SVGRectangle, error) {

	// -------------------------------------------------------------------------
	// Create the new object
	// -------------------------------------------------------------------------
	rect := new( SVGRectangle)

	// -------------------------------------------------------------------------
	// Create a unique ID for the object that will be used as the javascript ID
	// -------------------------------------------------------------------------
	displayID := fmt.Sprintf("%s_rect_%d", svg.DisplayID, len(svg.WidgetMap))
	rect.DisplayID = displayID

	// -------------------------------------------------------------------------
	// Parse the required parameters
	// -------------------------------------------------------------------------
	rect.Parent		= parent
	rect.X			= x
	rect.Y 			= y
	rect.Width 		= width
	rect.Height		= height
	rect.RX			= rx
	rect.RY 		= ry
	rect.ZIndex		= zIndex
	rect.Fill		= fill
	rect.Stroke 	= stroke
	rect.StrokeWidth= strokeWidth
	
	rect.Options = options
	rect.Styles = styles

	// -------------------------------------------------------------------------
	// add the object to the svg's widget map.
	// -------------------------------------------------------------------------
	svg.WidgetMap[rect.DisplayID] = len(svg.Widgets)
	svg.Widgets = append(svg.Widgets, rect)

	// -------------------------------------------------------------------------
	// return the new circle and nil error code
	// -------------------------------------------------------------------------
	return rect, nil
}











// *****************************************************************************
// Lua support stuff
// *****************************************************************************
const luaSVGTypeName = "svg"

// -----------------------------------------------------------------------------
// Registers circle type to given L.
// -----------------------------------------------------------------------------
func RegisterSVGType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaSVGTypeName)
	L.SetGlobal("svg", mt)
	// static attributes
	 //   L.SetField(mt, "new", L.NewFunction(newDisplay))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), svgMethods))
}


// -----------------------------------------------------------------------------
// svgMethods table
// -----------------------------------------------------------------------------
var svgMethods = map[string]lua.LGFunction{
    "newCircle": luaNewSVGCircle,
    "newRectangle": luaNewSVGRectangle,

 //   "show" : luaShow,
}




func getMetatableName(L *lua.LState, ud *lua.LUserData) string {
	udMt := L.GetMetatable(ud)
	if udMt == nil {
		return ""
	}

	// Access the _MT table from the registry
	mtTable := L.G.Registry.RawGetString("_MT")
	if tbl, ok := mtTable.(*lua.LTable); ok {
		var mtName string
		tbl.ForEach(func(key, value lua.LValue) {
			if value == udMt {
				mtName = key.String()
			}
		})
		return mtName
	}
	return ""
}


// -----------------------------------------------------------------------------
// Checks whether the first lua argument is a *LUserData with * SVG and 
// returns * svg.
// -----------------------------------------------------------------------------
func checkSVG(L *lua.LState) *SVG {
	ud := L.CheckUserData(1)
	if ud == nil {
		errstr := fmt.Sprintf("svg expected as the first parameter but got [nil].")
		L.ArgError(1, errstr)	
		return nil
	}
	if v, ok := ud.Value.(*SVG); ok {
		return v
	} else {
		name := getMetatableName(L, ud)
		errstr := fmt.Sprintf("svg expected as the first parameter but got [%s].", name)
		L.ArgError(1, errstr)
		return nil	
	}
}






// -----------------------------------------------------------------------------
// luaNewSVGCircle
// -----------------------------------------------------------------------------
func luaNewSVGCircle(L *lua.LState) int {

    svg := checkSVG(L)
    lua_top := L.GetTop()
    var circle  * SVGCircle
    var err     error
    var optionsMap map[string]string
    var stylesMap map[string]string

    centerX		:= L.CheckInt(2)
    centerY		:= L.CheckInt(3)
    radius		:= L.CheckInt(4)
    zIndex		:= L.CheckInt(5)
    fill		:= L.CheckString(6)
	stroke		:= L.CheckString(7)
	strokeWidth := L.CheckInt(8)

    switch lua_top {

        case 8:
            circle, err = svg.NewSVGCircle(svg.DisplayID, centerX, centerY, 
											radius, zIndex, fill, 
											stroke, strokeWidth, nil, nil) 
            if err != nil {
                L.ArgError(1, err.Error())
            }

        case 9:
            arg9 := L.Get(9) // Gets the first argument, even if it's nil
	        if arg9 == lua.LNil {
                optionsMap = nil
	        } else {
                optionsTable := L.CheckTable(9)
        	    optionsMap = luaTableToStringMap(optionsTable)
            }
            circle, err = svg.NewSVGCircle(svg.DisplayID, centerX, centerY, 
											radius, zIndex, fill, 
											stroke, strokeWidth, &optionsMap, 
											nil)   
        case 10:
            arg9 := L.Get(9) // Gets the first argument, even if it's nil
	        if arg9 == lua.LNil {
                optionsMap = nil
	        } else {
                optionsTable := L.CheckTable(9)
        	    optionsMap = luaTableToStringMap(optionsTable)
            }

            arg10 := L.Get(10) // Gets the first argument, even if it's nil
	        if arg10 == lua.LNil {
                stylesMap = nil
	        } else {
                stylesTable := L.CheckTable(10)
        	    stylesMap = luaTableToStringMap(stylesTable)
            }

            circle, err = svg.NewSVGCircle(svg.DisplayID, centerX, centerY, 
											radius, zIndex, fill, 
											stroke, strokeWidth, &optionsMap, 
											&stylesMap) 
            if err != nil {
                L.ArgError(1, err.Error())
            }

        default:
            msg := fmt.Sprintf("8, 9, or 10 arguments expected including object. Got %d.", lua_top)
            L.ArgError(1, msg)
            return 0
    }

	// Create the user data
    ud := L.NewUserData()
	ud.Value = circle

	// push it back to Lua
    L.SetMetatable(ud, L.GetTypeMetatable("svgCircle"))
    L.Push(ud)
    return 1
}


// -----------------------------------------------------------------------------
// luaNewSVGRectangle
// -----------------------------------------------------------------------------
func luaNewSVGRectangle(L *lua.LState) int {

    svg := checkSVG(L)
    lua_top := L.GetTop()
    var rect  * SVGRectangle
    var err     error
    var optionsMap map[string]string
    var stylesMap map[string]string

    x			:= L.CheckInt(2)
    y			:= L.CheckInt(3)
    width		:= L.CheckInt(4)
    height		:= L.CheckInt(5)
    zIndex		:= L.CheckInt(6)
    rx			:= L.CheckInt(7)
    ry			:= L.CheckInt(8)
    fill		:= L.CheckString(9)
	stroke		:= L.CheckString(10)
	strokeWidth := L.CheckInt(11)

	fmt.Printf("Rectangle %d parameters\n", lua_top)

    switch lua_top {

        case 11:
            rect, err = svg.NewSVGRectangle(svg.DisplayID, x, y, 
											width, height, zIndex, rx, ry, 
											fill, stroke, strokeWidth, 
											nil, nil) 
            if err != nil {
                L.ArgError(1, err.Error())
            }

        case 12:
            arg12 := L.Get(12) // Gets the first argument, even if it's nil
	        if arg12 == lua.LNil {
                optionsMap = nil
	        } else {
                optionsTable := L.CheckTable(12)
        	    optionsMap = luaTableToStringMap(optionsTable)
            }
            rect, err = svg.NewSVGRectangle(svg.DisplayID, x, y, 
										width, height, zIndex, rx, ry, 
										fill, stroke, strokeWidth, 
										&optionsMap, nil)   
        case 13:
            arg12 := L.Get(12) // Gets the first argument, even if it's nil
	        if arg12 == lua.LNil {
                optionsMap = nil
	        } else {
                optionsTable := L.CheckTable(12)
        	    optionsMap = luaTableToStringMap(optionsTable)
            }

            arg13 := L.Get(13) // Gets the first argument, even if it's nil
	        if arg13 == lua.LNil {
                stylesMap = nil
	        } else {
                stylesTable := L.CheckTable(13)
        	    stylesMap = luaTableToStringMap(stylesTable)
            }

			fmt.Printf("Lua styles length %d\n", len(stylesMap))

            rect, err = svg.NewSVGRectangle(svg.DisplayID, x, y, 
											width, height, zIndex, rx, ry, 
											fill, stroke, strokeWidth, &optionsMap, 
											&stylesMap) 
            if err != nil {
                L.ArgError(1, err.Error())
            }

        default:
            msg := fmt.Sprintf("11, 12, or 13 arguments expected including object. Got %d.", lua_top)
            L.ArgError(1, msg)
            return 0
    }

	// Create the user data
    ud := L.NewUserData()
	ud.Value = rect

	// push it back to Lua
    L.SetMetatable(ud, L.GetTypeMetatable("svgRectangle"))
    L.Push(ud)
    return 1
}

