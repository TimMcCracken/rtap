


package buffer_pool


import (
	"fmt"
	"runtime"
	"sync"
)


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




func getCallingPackage() (string) {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}
    
    functionName := f.Name()
    
    for i := len(functionName) - 1; i >= 0; i-- {
        if functionName[i] == '.' {
            packageName := functionName[:i]
            for j := len(packageName) - 1; j >= 0; j-- {
                if packageName[j] == '/' {
                    return packageName[j+1:i]
                }
            }
            return packageName + ":" + functionName
        }
    }
	return ""
}





// -----------------------------------------------------------------------------
// GeBuffer() returns a buffer from the buffer pool.
// -----------------------------------------------------------------------------
func (bp * BufferPool ) Get(size int) ( *[]byte) {

//	fmt.Println("Get Calling package:", getCallingPackage())

bp.totalGets++
	fmt.Printf("Total Buffers, Gets, Puts, scrubs: %d  %d  %d  %d\n", bp.totalBuffers, bp.totalGets, bp.totalPuts, bp.totalScrubs)
		return bp.pool.Get().(*[]byte)
}

// -----------------------------------------------------------------------------
// PutBuffer() sends a buffer to the scrubber
// -----------------------------------------------------------------------------
func (bp * BufferPool)Put(buffer *[]byte){

//	fmt.Println("Put Calling package:", getCallingPackage())

	bp.totalPuts++
	bp.scrubber <- buffer
}


//------------------------------------------------------------------------------
// scrubberTask zeroes out a buffer before returning it to the pool.
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
	