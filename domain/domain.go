

package domain

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
	"rtap/dacc"
	"rtap/rtdsms"
	bp "rtap/buffer_pool"
	"rtap/message_q"
	"rtap/metronome"
//	"rtap/hmi"
//	"os"
//	"time"
)



type Domain struct {
	Descriptor		DomainDescriptor
	messageQueue	message_q.MessageQ
	Datastores		[]*rtdsms.Datastore
	Datastores_map 	map[string]int

	dacc			dacc.DACC
	metronome		metronome.Metronome
	bufferPool		bp.BufferPool
}



func (domain * Domain) Start() {
	
	domain.bufferPool.Start()
	domain.messageQueue.Start(&domain.bufferPool)
	domain.metronome.Start(&domain.bufferPool, &domain.messageQueue)

	domain.dacc.Start(&domain.bufferPool, &domain.messageQueue)
}

func (domain * Domain)  MessageQueue() (* message_q.MessageQ) {

	return &domain.messageQueue
}






// -----------------------------------------------------------------------------
// construct   () 
// -----------------------------------------------------------------------------
func (domain * Domain) Construct(filename string) error  {

//	fmt.Printf("Constructing domain [%s]\n", domain.descriptor.domain_name)

	if domain.Datastores_map == nil {
		domain.Datastores_map = make(map[string]int)
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
	result := db.Raw("SELECT domain_id FROM domains where domain_name = ?", domain.Descriptor.Domain_name).Scan(&domain_id)
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
			ds, err := rtdsms.NewDatastore(domain.Descriptor.Realm_name, 
										domain.Descriptor.Domain_name, name )
			if err != nil {
				return fmt.Errorf("%v\n", err)
			}
			/*
			ds := rtdsms.Datastore {
			}
			ds.Descriptor.Datastore_name = name
			ds.Descriptor.Domain_name = domain.Descriptor.Domain_name
			ds.Descriptor.Realm_name = domain.Descriptor.Realm_name
			*/
			domain.Datastores_map[name] = len(domain.Datastores_map)
			domain.Datastores = append(domain.Datastores, ds)
		}
    }
    if err = rows.Err(); err != nil {
		return fmt.Errorf("%v\n", err)
    }

	for _, ds := range domain.Datastores {
		err = ds.Construct(filename)
		if err != nil {
			return fmt.Errorf("%v\n", err)
		}
	}



//	fmt.Printf("domain done.\n")

	return nil
}

