package hmi

import (
	"fmt"
	"net/http"
)

// element.write() is a helper function to write an element.
func (e *element) write(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<%s ", formElements[e.element_type])

	// TODO: input type
	if e.element_type == Input {
		fmt.Fprintf(w, "type='%s' ", inputTypes[e.input_type])
	}

	// TODO: generate ID

	// write the required and optional style attribute
	fmt.Fprintf(w, "style='")
	fmt.Fprintf(w, "position:absolute;")
	fmt.Fprintf(w, "top:%dpx;", e.top)
	fmt.Fprintf(w, "left:%dpx;", e.left)
	fmt.Fprintf(w, "width:%dpx;", e.width)
	fmt.Fprintf(w, "height:%dpx;", e.height)
	fmt.Fprintf(w, "%s", e.style)
	fmt.Fprintf(w, "' ")

	// write the remaining global attributes. Some values may be inherited,
	// so we all always explicitly set them.
	fmt.Fprintf(w, " contenteditable='%t'", e.editable) // may be inherited
	fmt.Fprintf(w, " draggable='%t'", e.draggable)      // 3 state-may be inherited or default

	if len(e.class) > 0 {
		fmt.Fprintf(w, " class='%s'", e.class)
	}

	// Dir actually has 3 states, so we may want to change this.
	if e.dir {
		fmt.Fprintf(w, " dir='rtl'")
	}
	if e.hidden {
		fmt.Fprintf(w, " hidden")
	}

	if len(e.lang) > 0 {
		fmt.Fprintf(w, " lang='%s'", e.lang)
	}
	if e.editable && e.spellcheck {
		fmt.Fprintf(w, " spellcheck='true'")
	}
	if len(e.title) > 0 {
		fmt.Fprintf(w, " title='%s'", e.title)
	}
	if e.tabindex > 0 {
		fmt.Fprintf(w, " tabindex='%d'", e.tabindex)
	}
	if e.translate {
		fmt.Fprintf(w, " translate='yes'")
	} else {
		fmt.Fprintf(w, " translate='no'")
	}

	// TODO: Output the event handlers
	for i := 0; i < len(e.events); i++ {
		fmt.Fprintf(w, " %s='%s'", FormElementEvents[e.events[i].event_type],
			e.events[i].script)
	}

	// ------------------------------------------------------------------------
	// 'Input types do not have content, are self closing and have additional
	// attributes.
	// ------------------------------------------------------------------------
	if e.element_type == Input {
		if len(e.value) > 0 {
			fmt.Fprintf(w, " value='%s'", e.value)
		}
		if len(e.placeholder) > 0 {
			fmt.Fprintf(w, " placeholder='%s'", e.placeholder)
		}
		if len(e.pattern) > 0 {
			fmt.Fprintf(w, " pattern='%s'", e.pattern)
		}

		if e.readonly {
			fmt.Fprintf(w, " readonly")
		}
		if e.disabled {
			fmt.Fprintf(w, " disabled")
		}
		if e.required {
			fmt.Fprintf(w, " required")
		}
		if e.multiple {
			fmt.Fprintf(w, " multiple")
		}
		if e.autofocus {
			fmt.Fprintf(w, " autofocus")
		}
		if e.autocomplete {
			fmt.Fprintf(w, " autocomplete='yes'")
		} else {
			fmt.Fprintf(w, " autocomplete='no'")
		}

		if e.size > 0 {
			fmt.Fprintf(w, " size='%d'", e.size)
		}
		if e.maxlength > 0 {
			fmt.Fprintf(w, " maxlength='%d'", e.maxlength)
		}

		if len(e.min) > 0 {
			fmt.Fprintf(w, " min='%s'", e.min)
		}
		if len(e.max) > 0 {
			fmt.Fprintf(w, " max='%s'", e.max)
		}
		if len(e.step) > 0 {
			fmt.Fprintf(w, " step='%s'", e.step)
		}

		fmt.Fprintf(w, "/>\n")
	} else {

		// Write the content, if any
		fmt.Fprintf(w, ">%s", e.content)
		//write the closing tag.
		fmt.Fprintf(w, "</%s>\n", formElements[e.element_type])

	}
}
