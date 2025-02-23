package hmi

import (
	"fmt"
	"io"
	"net/http"
)

// descriptor for HTML column
type column struct {
	db_name    string
	heading    string
	width      int32
	text_align string // make an enum?
	read_only  bool
}

// descriptor for HTML table
type table struct {
	top                      uint32
	left                     uint32
	width                    uint32
	height                   uint32
	color                    string
	background_color         string
	heading_color            string
	heading_background_color string
	odd_color                string
	even_color               string
	border_color             string
	border_style             string // make this use a list
	border_size              uint16
	border_collapse          bool
	style                    string
	rows                     int16
	columns                  []column
}

func (t *table) write(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "<table style='position:absolute;")
	fmt.Fprintf(w, "top:%dpx;", t.top)
	fmt.Fprintf(w, "left:%dpx;", t.left)
	fmt.Fprintf(w, "width:%dpx;", t.width)
	fmt.Fprintf(w, "height:%dpx;", t.height)
	if len(t.color) > 0 {
		fmt.Fprintf(w, "color:%s;", t.color)
	}
	if len(t.background_color) > 0 {
		fmt.Fprintf(w, "background-color:%s;", t.background_color)
	}
	fmt.Fprintf(w, "%s'>\n", t.style)

	// Show the table headings
	fmt.Fprint(w, "<tr style='")
	if len(t.heading_color) > 0 {
		fmt.Fprintf(w, "color:%s;", t.heading_color)
	}
	if len(t.heading_background_color) > 0 {
		fmt.Fprintf(w, "background-color:%s;", t.heading_background_color)
	}
	fmt.Fprint(w, "'>")

	for i := 0; i < len(t.columns); i++ {
		fmt.Fprintf(w, "<th style='width:%dpx;'>",
			t.columns[i].width)

		fmt.Fprintf(w, t.columns[i].heading)
		fmt.Fprintf(w, "</th>")
	}
	io.WriteString(w, "</tr>\n")

	// todo:Need a loop to print all the field.
	fmt.Fprintf(w, "<tr>")
	for i := 0; i < len(t.columns); i++ {
		fmt.Fprintf(w, "<td class='col%d'>%s</td>", i, "datum")
	}
	fmt.Fprintf(w, "</tr>\n")

	fmt.Fprintf(w, "</table>\n")

}
