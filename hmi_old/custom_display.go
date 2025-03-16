package hmi

import (
	"fmt"
	"net/http"
	"time"
)

type event struct {
	event_type FormElementEventType
	script     string
}

type element struct {
	element_type   FormElementType
	input_type     InputType
	form           string //id of the form generted automatically?
	formaction     string //overides 'action' in the form
	formenctype    string // overrides 'enctype' in the form
	formmethod     uint16 // overrides 'method' in the form
	formtarget     string // overrides 'target' in the form
	formnovalidate bool   // overrides 'novalidate' in the form
	novalidate     bool
	// required styles
	top     uint32
	left    uint32
	height  uint32
	width   uint32
	content string
	zindex  uint32
	// optional styles
	style string
	//The global attributes
	//NOTE: Access keys are not advised for multiple reasons
	class      string
	editable   bool
	dir        bool
	draggable  bool // may only be used in editor?
	hidden     bool
	id         string // appended with a dynamic value differentiator?
	lang       string
	spellcheck bool
	tabindex   uint16
	title      string
	translate  bool
	// The 'input' attributes
	value       string
	readonly    bool
	disabled    bool
	size        uint16
	maxlength   uint16
	min         string
	max         string
	multiple    bool
	pattern     string
	placeholder string
	required    bool
	step        string
	autofocus   bool
	//list TODO: Figure this out. probably an href string?
	autocomplete bool

	events []event
}

type form struct {
	accept_charset string
	action         string
	elements       []element
	enctype        string
	name           string
	rel            int16
	target         int16
	method         int16
	autocomplete   bool
	novalidate     bool
}

type customDisplay struct {
	// required styles
	background StandardColor
	color      StandardColor
	width      uint32
	height     uint32
	style      string
	elements   []element
	forms      []form
}

func (f *customDisplay) show(w http.ResponseWriter, req *http.Request) {

	// Get the current time to use in the headers.
	t := time.Now()

	// Set the headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache") // TODO: Extend this time later to something reasonable
	w.Header().Set("Date", fmt.Sprint(t.Format("Mon, 2 Jan 2006 15:04:05 MST")))

	fmt.Fprintf(w, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">\n")
	fmt.Fprintf(w, "<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")

	fmt.Fprintf(w, "<head>\n")
	fmt.Fprintf(w, "<title>PLACEHOLDER</title>\n") //TODO FIX THIS
	fmt.Fprintf(w, "</head>\n")

	fmt.Fprintf(w, "<body style='margin:0px; padding:0px'>\n")

	fmt.Fprintf(w, "<div style='width:%dpx; height:%dpx; background-color:%s;color:%s; %s'>\n",
		f.width, f.height, colors[f.background], colors[f.color], f.style)

	// write each of the elements
	for i := 0; i < len(f.elements); i++ {
		f.elements[i].write(w, req)
	}

	fmt.Fprintf(w, "</div>\n")
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")
}
