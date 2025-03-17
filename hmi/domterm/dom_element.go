
package domterm



type event struct {
	event_type string
	script     string
}


type GlobalAttributes struct {
	//NOTE: Access keys are not advised for multiple reasons
	Class     	string
	Editable  	bool
	Dir       	bool
	Draggable 	bool // may only be used in editor?
	EnterKeyHint	string
	Hidden    	bool
	ID        	string // appended with a dynamic value differentiator?
	Inert 		bool
	InputMode	string
	Lang      	string
	Spellcheck	bool
	Style		string
	Tabindex  	uint16
	Title     	string
	Translate 	bool
}


type InputAttributes struct {
	accept			string
	alt 			string
	autocomplete	string
	autofocus		string
	checked			bool
	dirname			bool
	disabled		bool
	form 			string
	formaction 		string
	formenctype 	string
	formmethod 		string
	formnovalidate	string
	formtarget 		string
	height			int 	// only for images
	list 			string
	max         	string
	maxlength   	uint16
	min         	string
	minlength 		string
	multiple    	bool
	name 			string
	pattern     	string
	placeholder 	string
	popovertarget	string
	popovertargetaction	string
	readonly    bool
	required    bool
	size        uint16
	src 		string
	step        string
	type_		string
	value       string
	width		int 		//only for images
}

