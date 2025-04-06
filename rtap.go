


package rtap

import(
	"fmt"
//	"gorm.io/gorm"
// 	"github.com/glebarez/sqlite"
	"log"
	"os"
//	"rtap/hmi"
	"rtap/message_q"
	"rtap/realm"
)



type rtapSystem struct {
	realms			[]*realm.Realm
	realms_map 		map[string]int
//	hmi				hmi.HMI
}

// ----------------------------------------------------------------------------
// "RTAP" is a singleton and is the root object of the RTAP system.
// ----------------------------------------------------------------------------
var RTAP rtapSystem




// Start starts up the entire system. 
func (rtap * rtapSystem) Start(serverAddress string) {

	// Start all the realms, which will start all the domains.


	for _, realm := range rtap.realms {
		realm.Start()
	}


	// Start the web server
	go rtap.HMIServerTask(serverAddress)

}




func (rtap * rtapSystem)Realm(realm_name string) (* realm.Realm, error) {

	realm_index, ok := RTAP.realms_map[realm_name]
	if ! ok {
		return nil, fmt.Errorf("Realm [%s] not found.", realm_name)
	} 

	return RTAP.realms[realm_index], nil
}


func (rtap * rtapSystem) MessageQueue(realm_name string, domain_name string) (* message_q.MessageQ, error) {

	realm_index, ok := RTAP.realms_map[realm_name]
	if ! ok {
		return nil, fmt.Errorf("Realm [%s] not found.", realm_name)
	} 

	return RTAP.realms[realm_index].MessageQueue(domain_name)
}




// -----------------------------------------------------------------------------
//Function to initialize the "RTAP" structure at statup. It functions by 
// looping through all the folders within the "base" folder that is declared
// as the environment variable.
// -----------------------------------------------------------------------------
func init(){

	RTAP.realms_map = make(map[string]int)
	base := os.Getenv("rtdsms_base")

	// Open the directory
	d, err := os.Open(base)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	// Read the directory contents
	files, err := d.Readdir(-1) // -1 means no limit on number of files
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the files in the directory
	for _, file := range files {
		// Check if it's a directory and not the current or parent directory
		if file.IsDir() && file.Name() != "." && file.Name() != ".." {
			
			// Print the name of the subdirectory
			//fmt.Println(file.Name())

			// create a realm in 'System'
			/*
			realm := realm.Realm{	
				path: base + "/" + file.Name(),
				filename: base + "/" + file.Name() + "/" + file.Name() + ".realm",
				} 
			realm.descriptor.realm_name = file.Name()
			*/
			realm, err := realm.NewRealm(file.Name())
			if err != nil {
				log.Fatal(err)
			}
		
			RTAP.realms_map[file.Name()]= len(RTAP.realms_map)
			RTAP.realms = append(RTAP.realms, realm)

		}
	}

	serverAddress := ":8080" // TODO: GET FROM THE SQLITE DATABASE
	RTAP.Start(serverAddress)
}
