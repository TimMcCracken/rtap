
// -----------------------------------------------------------------------------
// Label widget. For the digial clock we will utilize the 'output' type.
// -----------------------------------------------------------------------------

package widget


import (
	_ "embed"
	"fmt"
	"rtap/hmi/domterm"
	"time"
	"gorm.io/gorm"
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

func (dc *Label) Init(db *gorm.DB, descriptor string, rtdata string, display_id string){

	dc.DisplayID = display_id

	//TODO: read from database
}

func (lbl *Label) Put(conn *websocket.Conn){

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
func (dc *Label) Update(conn *websocket.Conn){
}

