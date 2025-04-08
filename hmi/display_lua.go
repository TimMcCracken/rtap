


package hmi

import (
	"fmt"
	"github.com/yuin/gopher-lua"
    "rtap/hmi/widget"
)


// ----------------------------------------------------------------------------
// help function to convert a Lua table to a Go map[string]string
// ----------------------------------------------------------------------------
func luaTableToStringMap(tbl *lua.LTable) map[string]string {
	result := make(map[string]string)

	tbl.ForEach(func(key, value lua.LValue) {
		// Only handle string keys and string values
		k, ok1 := key.(lua.LString)
		v, ok2 := value.(lua.LString)
		if ok1 && ok2 {
			result[string(k)] = string(v)
		}
	})

	return result
}


const luaDisplayTypeName = "display"

    // Registers my person type to given L.
    func RegisterDisplayType(L *lua.LState) {
        mt := L.NewTypeMetatable(luaDisplayTypeName)
        L.SetGlobal("display", mt)
        // static attributes
     //   L.SetField(mt, "new", L.NewFunction(newDisplay))
        // methods
        L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), displayMethods))
    }

    // Constructor 
    /*
    func newDisplay(L *lua.LState) int {
        display := &Display{L.CheckString(1)}
        ud := L.NewUserData()
        ud.Value = display
        L.SetMetatable(ud, L.GetTypeMetatable(luaDisplayTypeName))
        L.Push(ud)
        return 1
    }*/

    // Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
    func checkDisplay(L *lua.LState) *Display {
        ud := L.CheckUserData(1)
        if v, ok := ud.Value.(*Display); ok {
            return v
        }
        L.ArgError(1, "display expected")
        return nil
    }

// -----------------------------------------------------------------------------
// displayMethods table
// -----------------------------------------------------------------------------
var displayMethods = map[string]lua.LGFunction{
    "newLabel": luaNewLabel,
    "newDigitalClock": luaNewDigitalClock,

    "show" : luaShow,
}


func luaShow(L *lua.LState) int {
    d := checkDisplay(L)

    d = d
    return 0
}




// -----------------------------------------------------------------------------
// luaNewLabel
// -----------------------------------------------------------------------------
func luaNewLabel(L *lua.LState) int {

    d := checkDisplay(L)
    lua_top := L.GetTop()
    var lbl     * widget.Label
    var err     error
    var optionsMap map[string]string
    var stylesMap map[string]string

    // Check that we got the correct number of arguments.
//    if L.GetTop() != 8 {
 //       L.ArgError(1, "8 arguments expected including object")
 //       return 0
 //   }

    parent  := L.CheckString(2)
    top     := L.CheckInt(3)
    left    := L.CheckInt(4)
    width   := L.CheckInt(5)
    height  := L.CheckInt(6)
    zIndex  := L.CheckInt(7)
    content := L.CheckString(8)

    switch lua_top {

        case 8:
            lbl, err = d.NewLabel(parent, top, left, width, height, zIndex, content, nil, nil) 
            if err != nil {
                fmt.Printf("We gots a problem\n")
            }

        case 9:
            arg9 := L.Get(9) // Gets the first argument, even if it's nil
	        if arg9 == lua.LNil {
                optionsMap = nil
	        } else {
                optionsTable := L.CheckTable(9)
        	    optionsMap = luaTableToStringMap(optionsTable)
            }
            lbl, err = d.NewLabel(parent, top, left, width, height, zIndex, content, optionsMap, nil) 
    
  
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



            lbl, err = d.NewLabel(parent, top, left, width, height, zIndex, content, optionsMap, stylesMap) 
            if err != nil {
                fmt.Printf("We gots a problem\n")
            }

        default:

            L.ArgError(1, "8, 9, or 10 arguments expected including object")
            return 0
            
    }

    ud := L.NewUserData()
	ud.Value = &lbl

    L.SetMetatable(ud, L.GetTypeMetatable("label"))
    L.Push(ud)
    return 1
}


// -----------------------------------------------------------------------------
// luaNewDigitalClock
// -----------------------------------------------------------------------------
func luaNewDigitalClock(L *lua.LState) int {
    d := checkDisplay(L)

    // Check that we got the correct number of arguments.
    if L.GetTop() != 8 {
        L.ArgError(1, "8 arguments expected including object")
        return 0
    }

    parent  := L.CheckString(2)
    top     := L.CheckInt(3)
    left    := L.CheckInt(4)
    width   := L.CheckInt(5)
    height  := L.CheckInt(6)
    zIndex  := L.CheckInt(7)
    content := L.CheckString(8)

    // TODO: options and styles
    

    dc, err := d.NewDigitalClock(parent, top, left, width, height, zIndex, content, nil, nil) 
    if err != nil {
        fmt.Printf("We gots a problem\n")
    }

    ud := L.NewUserData()
	ud.Value = &dc

    L.SetMetatable(ud, L.GetTypeMetatable("digitalClock"))
    L.Push(ud)
    return 1
}







    // Getter and setter for the Person#Name
    func displayGetSetName(L *lua.LState) int {
        d := checkDisplay(L)

        d=d

        if L.GetTop() == 2 {
          //  d.Name = L.CheckString(2)
            return 0
        }
        L.Push(lua.LString("THIS IS THE DISPLAY"))
//        L.Push(lua.LString(d.Name))
        return 1
    }
