

package rtdsms

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
//	"os"
//	"time"
)



type domain struct {
	descriptor		domainDescriptor
	datastores		[]*datastore
	datastores_map 	map[string]int
}




// -----------------------------------------------------------------------------
// construct   () 
// -----------------------------------------------------------------------------
func (domain * domain) construct(filename string) error  {

//	fmt.Printf("Constructing domain [%s]\n", domain.descriptor.domain_name)

	if domain.datastores_map == nil {
		domain.datastores_map = make(map[string]int)
	}

	// Open the sqlite file
	db, err := gorm.Open(sqlite.Open(filename + "?_pragma=foreign_keys(1)" ), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database using path [%s]", filename)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Internal error: failed to obtain database.DB in CreateDomain()")
	}
	defer sqlDB.Close()


	// ------------------------------------------------------------------------
	// Get the domain_id
	// ------------------------------------------------------------------------
	var domain_id	int64
	result := db.Raw("SELECT domain_id FROM domains where domain_name = ?", domain.descriptor.domain_name).Scan(&domain_id)
	if result.Error != nil {
		return result.Error
	}


	// ------------------------------------------------------------------------
	// Select the rows
	// ------------------------------------------------------------------------
	rows, err := sqlDB.Query("SELECT datastore_name FROM datastores where domain_id = ? order by datastore_id", domain_id)
    if err != nil {
        return err
    }
    defer rows.Close()

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
		var name string
        err := rows.Scan(&name)
		if err != nil {
            return fmt.Errorf("%v\n", err)
        } else {
			ds := datastore {
			}
			ds.descriptor.datastore_name = name
			ds.descriptor.domain_name = domain.descriptor.domain_name
			ds.descriptor.realm_name = domain.descriptor.realm_name
			domain.datastores_map[name] = len(domain.datastores_map)
			domain.datastores = append(domain.datastores, &ds)
		}
    }
    if err = rows.Err(); err != nil {
		return fmt.Errorf("%v\n", err)
    }

	for _, ds := range domain.datastores {
		err = ds.construct(filename)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}



//	fmt.Printf("domain done.\n")

	return nil
}

