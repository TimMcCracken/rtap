/* -----------------------------------------------------------------------------
metronome

The metronome is the "heartbeat of the system. Once each second it sends a 
broadcast message to the messageQueue with a structure that indicates the 
seconds of offset since the most recent minute, hour and day. 

The metronome only calculates offset for seconds and minutes. Time zone 
dependencies could vary with different systems, and even different goroutines
and so should be calculated as needed.

Rev Date     By  Reason
--- -------- --- ---------------------------------------------------------------


------------------------------------------------------------------------------*/
package metronome
 
import (
	"fmt"
	"rtap/message_q"
	pb 		"rtap/metronome/metronome.pb"
	proto 	"github.com/golang/protobuf/proto"
	"time"
)


func init() {
	fmt.Println("initing metronome")
} 



/*
type Offsets struct{
	Seconda_1	int 
	Seconda_2	int 
	Seconda_3	int
	Seconda_4	int
	Seconda_5	int
	Seconda_6	int
	Seconda_10	int
	Seconda_12	int
	Seconda_15	int
	Seconda_20	int
	Seconda_30	int

	Minutes_1	int
	Minutes_2	int
	Minutes_3	int
	Minutes_4	int
	Minutes_5	int
	Minutes_6	int
	Minutes_10	int
	Minutes_12	int
	Minutes_15	int
	Minutes_20	int
	Minutes_30	int

	Hours_0		int
	Hours_1		int
	Hours_2		int
	Hours_3		int
	Hours_4		int
	Hours_5		int
	Hours_6		int

}

type Tic struct{

	Second	int
	Minute  int
	Hour 	int

	UTC		time.Time
	Local	time.Time
	IsDST	bool

	// todoL add variables for NASA.Celestial time


	Offsets Offsets
}
*/





func Metronome(mq * message_q.MessageQ) {

	// Declare teh structure for the data to be marshalled.
	var pb_tick pb.Tick

	// fill the structure with data.
	for {

		fmt.Printf(".")

		// get the current time.
		now := time.Now()
		
		// calculate a sleep time to get us close to the top of the second.
		soon 			:= now.Add(1 * time.Second)
		soon_rounded 	:= soon.Round(time.Second)
		wait_time 		:= time.Until(soon_rounded)
        time.Sleep(wait_time)

		// get the current second and minute.
		second	:= int32(now.Second())
		minute	:= int32(now.Minute())
	
		pb_tick.Second = int32(second)
		pb_tick.Minute = int32(minute)
	
		// calculate the offsetts
		pb_tick.Seconda_2 = second % 2
		pb_tick.Seconda_3 = second % 3
		pb_tick.Seconda_4 = second % 4
		pb_tick.Seconda_5 = second % 5
		pb_tick.Seconda_6 = second % 6
		pb_tick.Seconda_10 = second % 10
		pb_tick.Seconda_12 = second % 12
		pb_tick.Seconda_15 = second % 15
		pb_tick.Seconda_20 = second % 20
		pb_tick.Seconda_30 = second % 30
		
		pb_tick.Minutes_1 = ((minute * 60) + second) % 60
		pb_tick.Minutes_2 = ((minute * 60) + second) % 120
		pb_tick.Minutes_3 = ((minute * 60) + second) % 180
		pb_tick.Minutes_4 = ((minute * 60) + second) % 240
		pb_tick.Minutes_5 = ((minute * 60) + second) % 300
	
		pb_tick.Minutes_6 = ((minute * 60) + second) % 360
		pb_tick.Minutes_10 = ((minute * 60) + second) % 600
		pb_tick.Minutes_12 = ((minute * 60) + second) % 720
		pb_tick.Minutes_15 = ((minute * 60) + second) % 900
		pb_tick.Minutes_20 = ((minute * 60) + second) % 1200
		pb_tick.Minutes_30 = ((minute * 60) + second) % 1800

		// -------------------------------------------------------------------------
		// Get a buffer from the pool
		// -------------------------------------------------------------------------
		buf_ptr := mq.GetBuffer().(*[]byte)
		
		// -------------------------------------------------------------------------
		// Marshall the data and copy it into the buffer
		// -------------------------------------------------------------------------

		data, err := proto.Marshal(&pb_tick)
		if err != nil {
			// TODO: fix this
    		fmt.Printf("Failed to encode address book:", err)
		}

		// -------------------------------------------------------------------------
		// Ensure buffer is large enough
		// -------------------------------------------------------------------------
		if len(data) > cap(*buf_ptr) {
			fmt.Printf("Buffer is too small.")
		}

		// -------------------------------------------------------------------------
		// Copy marshaled data into the buffer to minimize allocations
		// -------------------------------------------------------------------------
		copy(*buf_ptr, data)		

		// -------------------------------------------------------------------------
		// Send the data to the message_q.
		// -------------------------------------------------------------------------
		destinations := []string {"*"}
		mq.Send(destinations, buf_ptr)

	}
}