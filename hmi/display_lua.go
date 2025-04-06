


package hmi

import (
//	_ "embed"
//	"encoding/json"
//	"flag"
	"fmt"
//	"html/template"
//	"log"
//	"maps"
//	"net/http"
//	"time"
	 
	"github.com/yuin/gopher-lua"
//	"github.com/gorilla/websocket"
//	"rtap/hmi/domterm"
//	mq "rtap/message_q"
//	"rtap/hmi/widget"
)





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
    
    lbl, err := d.NewLabel(parent, top, left, width, height, zIndex, content) 
    if err != nil {
        fmt.Printf("We gots a problem\n")
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
    
    dc, err := d.NewDigitalClock(parent, top, left, width, height, zIndex, content) 
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
