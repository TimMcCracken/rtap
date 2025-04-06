/*******************************************************************************

The HMI subsysterm consts of 3 major components:

	1.	hmiServerTask() (runs at the 'system' ;evel and belongs to the RTAP 
		package)
	2.	hmiDisplayTask() (runs at the domain level - with one instance for each
		domain in each realm)
	3.	hmiDisplayWorker() runs at the hmi level, with one instance for each
		open display. Note that a display can also be a subset of the "desktop"


*******************************************************************************/



package rtap


import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	hmi "rtap/hmi"
	"github.com/gorilla/mux"
)


func homeStub(w http.ResponseWriter, r *http.Request) {
		wsTemplate.Execute(w, "ws://" +r.Host+ "/ws/authenticate")
	}
	
func chooserStub(w http.ResponseWriter, r *http.Request) {
	wsTemplate.Execute(w, "ws://" +r.Host+ "/ws/chooser")
}
	
func displayStub(w http.ResponseWriter, r *http.Request) {
	wsTemplate.Execute(w, "ws://" +r.Host+ "/ws" + r.URL.String())
}
	
	
func domtermHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(domterm_js))
}
	
	
// -----------------------------------------------------------------------------
// HMIServerTask starts the web server. It runs at the rtap_system scope/
// -----------------------------------------------------------------------------
func (RTAP * rtapSystem ) HMIServerTask(serverAddress string) {
		

	r := mux.NewRouter()

	// Route with regex constraint: only numbers for ID
	r.HandleFunc("/", homeStub)
	r.HandleFunc("/ws/authenticate", hmi.AuthenticateHandler)
	
	r.HandleFunc("/chooser", chooserStub)
	r.HandleFunc("/ws/chooser", hmi.ChooserHandler)
	
	r.HandleFunc("/display/r/{realm:[0-9a-zA-Z]+}/d/{domain:[0-9a-zA-Z]+}/n/{name:[0-9a-zA-Z]+}", displayStub)

	
	realm, _ := RTAP.Realm("devscada")
	if realm == nil {
		fmt.Printf("****************************** REALM not exist")
	}
	domain, _ := realm.Domain("scada")
	if realm == nil {
		fmt.Printf("****************************** Domain not exist")
	}



	r.HandleFunc("/ws/display/r/{realm:[0-9a-zA-Z]+}/d/{domain:[0-9a-zA-Z]+}/n/{name:[0-9a-zA-Z]+}", domain.DisplayHandlerStub)
	
	r.HandleFunc("/js/domterm.js", domtermHandler)


	// Optional: a 404 handler
//	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		http.Error(w, "Custom 404 - Not Found", http.StatusNotFound)
//	})

	// Start the HMI Server loop.
	fmt.Println("WebSocket server started on", serverAddress)

	log.Fatal(http.ListenAndServe(serverAddress, r))

}
		

// -----------------------------------------------------------------------------
// The following commented lines are go directives, not comments!
// We are going to embed two files which are called by all the stubs above.
// The template_html file is a barebones by valid HTML file, which has a line
// in the <head> to import the domterm.js file, which includes all the 
// javascript required for the 'dom terminal'
// -----------------------------------------------------------------------------

//go:embed template.html
var template_html 	string

//go:embed hmi/domterm.js
var  domterm_js			string
	
var wsTemplate 	= template.Must(template.New("").Parse(template_html))
