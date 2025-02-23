package hmi

import (
	"fmt"
	"net/http"
	"time"
)

type tabularDisplay struct {
	title            string
	color            string
	background_color string
	font_family      string
	table            table
}

func (d *tabularDisplay) show(w http.ResponseWriter, req *http.Request) {
	fmt.Println("tabular test handler")

	// Get the current time to use in the headers.
	t := time.Now()

	// Set the headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache") // TODO: Extend this time later to something reasonable
	w.Header().Set("Date", fmt.Sprint(t.Format("Mon, 2 Jan 2006 15:04:05 MST")))

	// Override the table location since this is a pre-defined display type.
	d.table.top = 100
	d.table.left = 50

	fmt.Fprintf(w, "<head>\n")
	fmt.Fprintf(w, "<title>%s</title>\n", d.title)
	// do the table border style as it needs to be internal, rather than inline for efficiency
	fmt.Fprintf(w, "<style>\n")
	if len(d.table.border_style) > 0 {
		fmt.Fprintf(w, "table, th, td{border: %dpx %s %s;",
			d.table.border_size, d.table.border_style, d.table.border_color)
		if d.table.border_collapse == true {
			fmt.Fprintf(w, "border-collapse: collapse;}")
		}
		fmt.Fprintf(w, "\n")
	}

	// set the background colors for the odd & even rows.
	if len(d.table.odd_color) > 0 {
		fmt.Fprintf(w, "tr:nth-child(odd) {background-color: %s;}\n", d.table.odd_color)
	}
	if len(d.table.even_color) > 0 {
		fmt.Fprintf(w, "tr:nth-child(even) {background-color: %s;}\n", d.table.even_color)
	}

	for i := 0; i < len(d.table.columns); i++ {
		if len(d.table.columns[i].text_align) > 0 {
			fmt.Fprintf(w, ".col%d{text-align:%s;}\n", i, d.table.columns[i].text_align)
		}

	}
	fmt.Fprintf(w, "</style>\n")
	fmt.Fprintf(w, "</head>\n")

	fmt.Fprintf(w, "<body style='margin:0px;padding:0px;")
	fmt.Fprintf(w, "color:%s;", d.color)
	fmt.Fprintf(w, "background-color:%s;", d.background_color)
	fmt.Fprintf(w, "font-family:%s;", d.font_family)
	fmt.Fprintf(w, "'>\n")

	fmt.Fprintf(w, "<label style='position:absolute;top:25px;left:25px;")
	fmt.Fprintf(w, "text-align:center;width:%dpx;font-size:30px;'>%s</label>\n", d.table.width, d.title)
	d.table.write(w, req)
	fmt.Fprintf(w, "</body>\n</html>")
}
