package hmi

type StandardColor int8

// should we add dark gray and light gray? others?
const (
	None StandardColor = iota
	Aqua
	Black
	Blue
	Cyan
	Fuscia
	Gray
	Green
	Indigo
	Lime
	LightBlue
	LightGreen
	Magenta
	Maroon
	Navy
	Olive
	Orange
	Pink
	Purple
	Red
	Silver
	Teal
	Violet
	White
	Yellow
)

var colors [24]string = [24]string{"Aqua", "Black", "Blue", "Cyan", "Fuscia", "Gray",
	"Green", "Indigo", "Lime", "LightBlue", "LightGreen", "Magenta", "Maroon", "Navy",
	"Olive", "Orange", "Pink", "Purple", "Red", "Silver", "Teal", "Violet",
	"White", "Yellow"}

// declarations for FormElementType
// These are required to match domain tables in the database.
type FormElementType int16

const (
	Label FormElementType = iota
	Input
	Select
	TextArea
	//Button
	FieldSet
	Legend
	DataList
	Output
	Option
	OptGroup
	Table
)

var formElements [11]string = [11]string{"label", "input", "select", "textarea",
	"fieldset", "legend", "datalist", "output", "option", "optgroup", "table"}

// declarations for InputType
// These are required to match domain tables in the database.
type InputType int16

const (
	Button InputType = iota
	Checkbox
	Color
	Date
	DateTimeLocal
	Email
	File
	Hidden
	Image
	Month
	Number
	Password
	Radio
	Range
	Reset
	Search
	Submit
	Tel
	Text
	Time
	Url
	Week
)

var inputTypes [22]string = [22]string{"Button", "Checkbox", "Color", "Date",
	"DateTimeLocal", "Email", "File", "Hidden", "Image", "Month", "Number",
	"Password", "Radio", "Range", "Reset", "Search", "Submit", "Tel", "Text",
	"Time", "Url", "Week"}

type FormElementEventType int16

const (
	// Element Events
	OnBlur FormElementEventType = iota
	OnChange
	OnContextMenu
	OnFocus
	OnInput
	OnInvalid
	OnReset
	OnSearch
	OnSelect
	OnSubmit

	// key events
	OnKeyDown
	OnKeyPress
	OnKeyUp

	// Mouse events
	OnClick
	OnDblClick
	OnMouseDown
	OnMouseMove
	OnMouseOut
	OnMouseOver
	OnMouseUp
	OnWheel

	//Drag events
	OnDrag
	OnDragEnd
	OnDragCenter
	OnDragLeave
	OnDragOver
	OnDragStart
	OnDrop
	Onscroll

	// Clipboard events
	OnCopy
	OnCut
	OnPaste
)

var FormElementEvents [32]string = [32]string{
	"onblur", "onchange", "oncontextmenu", "onfocus", "oninput", "oninvalid",
	"onreset", "onsearch", "onselect", "onsubmit",
	"onkeydown", "onkeypress", "onkeyup",
	"onclick", "ondblclick", "onmousedown", "onmousemove", "onmouseout",
	"onmouseover", "onmouseup", "onwheel",
	"ondrag", "ondragend", "ondragcenter", "ondragleave", "ondragover",
	"ondragstart", "ondrop", "onscroll", "oncopy", "oncut", "onpaste"}
