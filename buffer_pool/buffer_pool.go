


package buffer_pool


import (
	"fmt"
	"sync"
//	"gorm.io/gorm"
//	"github.com/glebarez/sqlite"
//	"rtap/rtdsms"
//	"rtap/message_q"
//	"os"
//	"time"
)


/*
TODO: make a buffer pool a separate packe instantiated in the domain. Add a calling parameter for future use as the size of the buffer
Then all the packages can request buffers from a common pool, pass them arround and not worry about
memory leaks
*/

type BufferPool struct {

	pool			sync.Pool

	scrubber		chan *[]byte // returned buffers to be scrubbed.
	totalBuffers	int
	totalGets		int
	totalPuts 		int
	totalScrubs		int

}


func (bp * BufferPool )Start() {
	
	if bp.scrubber == nil {
		bp.scrubber		= make(chan *[]byte, 256 )
	}

//	if bp.pool == nil {
		bp.pool = sync.Pool{
			New: func() interface{} {
				buf := make([]byte, 1024) // Allocate a 1KB buffer (adjust size as needed)
				bp.totalBuffers++
			return &buf
			},
		}
//	}
	go scrubberTask(bp)

}

// -----------------------------------------------------------------------------
// GeBuffer() returns a buffer from the buffer pool.
// -----------------------------------------------------------------------------
 	func (bp * BufferPool ) Get(size int) ( interface{}) {
	bp.totalGets++

	fmt.Printf("Total Buffers, Gets, Puts, scrubs: %d  %d  %d  %d\n", bp.totalBuffers, bp.totalGets, bp.totalPuts, bp.totalScrubs)
	
	return bp.pool.Get().(*[]byte)
}

// -----------------------------------------------------------------------------
// PutBuffer() sends a buffer to the scrubber
// -----------------------------------------------------------------------------
func (bp * BufferPool)Put(buffer *[]byte){
	bp.totalPuts++
	bp.scrubber <- buffer
}

//------------------------------------------------------------------------------
// scrubber loop zeroes out a buffer before returning it to the pool.
//------------------------------------------------------------------------------
func scrubberTask(bp * BufferPool){

	//	fmt.Println("scrubber loop started")
		for {
			buffer := <- bp.scrubber
	//		fmt.Printf("Scrubbing %d bytes\n", len(*data))
			for  i := range *buffer {
				(*buffer)[i] = 0
			}
			bp.pool.Put(buffer)
			bp.totalScrubs++
		}
	}
	