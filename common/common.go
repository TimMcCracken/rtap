

package common

import(
	"fmt"
	"regexp"
	"rtap/options"
)


func ValidateObjectName(name string) error {

	// Validate the non-type dependent parameters in the descriptor
	if len(name) > options.MaxObjectNameLength {
		return	fmt.Errorf("Object name is too long: %s", name)
	}

	// Check if the name contains any disallowed characters
	regex := regexp.MustCompile(options.ObjectNameFormat)
	if ! regex.MatchString(name) {
		return 	fmt.Errorf("Object Name contains invalid characters: %s", name)
	}	 
	
	return nil
}

