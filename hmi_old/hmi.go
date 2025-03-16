package hmi

import (
	"fmt"
	"net/http"
	"rtap/dbapi" // This will move when we move the HTTP server to its own package
	// This will move when we move the HTTP server to its own package
)

/*
func (d *display) write(w http.ResponseWriter, req *http.Request) {


}
*/

var testdisplay customDisplay

var testtabular tabularDisplay

func init() {

	// initialize the tabular display

	testtabular.title = "Test Tabular #1"

	testtabular.table.width = 1000
	testtabular.color = "white"
	testtabular.background_color = "black"
	testtabular.font_family = "arial,helvetica,san serif"

	testtabular.table.border_color = "white"
	testtabular.table.border_style = "solid" // make this use a list
	testtabular.table.border_size = 2
	testtabular.table.border_collapse = true

	testtabular.table.heading_color = "blue"
	testtabular.table.heading_background_color = "gray"

	testtabular.table.odd_color = "#C4A484"
	testtabular.table.even_color = "lightgreen" //"#90EE90"

	testtabular.table.color = "black"

	testtabular.table.columns = append(testtabular.table.columns, column{"myname", "Column 0", 100, "left", false})
	testtabular.table.columns = append(testtabular.table.columns, column{"myname2", "Column 1", 200, "right", false})
	testtabular.table.columns = append(testtabular.table.columns, column{"myname3", "Column 2", 400, "left", false})
	testtabular.table.columns = append(testtabular.table.columns, column{"myname4", "Column 3", 150, "left", false})
	testtabular.table.columns = append(testtabular.table.columns, column{"myname5", "Column 4", 150, "left", false})

	// Initialize the custom display
	testdisplay.background = None
	testdisplay.color = Cyan
	testdisplay.width = 900
	testdisplay.height = 600
	testdisplay.style = "font-family:sans-serif;font-size:12px;"

	x := 5
	testdisplay.elements = make([]element, x)
	testdisplay.elements[0] = element{element_type: Label, left: 100, top: 50, width: 100, height: 25, content: "Label1"}
	testdisplay.elements[1] = element{element_type: Input, input_type: Button, left: 100, top: 100, width: 100, height: 25, draggable: true, title: "Help Me!", content: "Label2", style: "color:red;", value: "Click Me!"}
	testdisplay.elements[2] = element{element_type: Input, input_type: Text, left: 100, top: 150, width: 100, height: 25, style: "color:blue;", placeholder: "Enter something here!", title: "real help!", autofocus: true}
	testdisplay.elements[3] = element{element_type: Label, left: 100, top: 200, width: 100, height: 25, content: "Label4", style: "color:blue;"}
	testdisplay.elements[4] = element{element_type: Label, left: 100, top: 250, width: 100, height: 25, content: "Label5", style: "background-color:gray;"}

	testdisplay.elements[4].events = make([]event, 1)
	testdisplay.elements[4].events[0].event_type = OnClick
	testdisplay.elements[4].events[0].script = "alert(\"I am an alert box!\");"
}

/*
func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("hello handler")
	fmt.Fprintf(w, "hello\n")
}
func headers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("headers handler")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}*/

// ----------------------------------------------------------------------------
// tabularHandler() parses the http.Request and selects the tabular display
// to show.
// ----------------------------------------------------------------------------
func tabularHandler(w http.ResponseWriter, req *http.Request) {

	testtabular.show(w, req)
}

// ----------------------------------------------------------------------------
// customeHandler() parses the http.Request and selects the custom display
// to show.
// ----------------------------------------------------------------------------
func customHandler(w http.ResponseWriter, req *http.Request) {

	testdisplay.show(w, req)

}

func Serve() {

	fmt.Println("starting serve")

	//http.HandleFunc("/hello", hello)
	//http.HandleFunc("/headers", headers)
	http.HandleFunc("/custom", customHandler)
	http.HandleFunc("/tabular", tabularHandler)
	http.HandleFunc("/dbapi/", dbapi.DBapiHandler)

	http.ListenAndServe(":8090", nil)

	fmt.Println("exiting serve")
}
