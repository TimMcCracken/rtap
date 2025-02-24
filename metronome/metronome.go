/* -----------------------------------------------------------------------------
metronome

The metronome is the "heartbeat of the system. Once each second it sends a 
broadcast message to the messageQueue with a structure that indicates the 
seconds of offset since the most recent minute, hour and day. 

The values are based on either UTC 

Rev Date     By  Reason
--- -------- --- ---------------------------------------------------------------


------------------------------------------------------------------------------*/
package metronome
 
import (
	"fmt"
	"time"
)


func init() {
	fmt.Println("initing metronome")
} 




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






func Metronome() {

	var tic Tic

	now := time.Now()


	soon := now.Add(1 * time.Second)
	soon_rounded := soon.Round(time.Second)
	wait_time := time.Until(soon_rounded)
	time.Sleep(wait_time)

	for {

		now := time.Now()
		fmt.Printf("%v\n", now)

		second := now.Second()
		minute := now.Minute()
		hour := now.Hour()

		tic.Local = now
		tic.UTC = now.UTC()
		tic.IsDST = now.IsDST()
		tic.Second = second
		tic.Minute = minute
		tic.Hour = hour
	

		// calculate the offsetts
		tic.Offsets.Seconda_1 = 0
		tic.Offsets.Seconda_2 = second % 2
		tic.Offsets.Seconda_3 = second % 3
		tic.Offsets.Seconda_4 = second % 4
		tic.Offsets.Seconda_5 = second % 5
		tic.Offsets.Seconda_6 = second % 6
		tic.Offsets.Seconda_10 = second % 10
		tic.Offsets.Seconda_12 = second % 12
		tic.Offsets.Seconda_15 = second % 15
		tic.Offsets.Seconda_20 = second % 20
		tic.Offsets.Seconda_30 = second % 30
		
		tic.Offsets.Minutes_1 = ((minute * 60) + second) % 60
		tic.Offsets.Minutes_2 = ((minute * 60) + second) % 120
		tic.Offsets.Minutes_3 = ((minute * 60) + second) % 180
		tic.Offsets.Minutes_4 = ((minute * 60) + second) % 240
		tic.Offsets.Minutes_5 = ((minute * 60) + second) % 300
	
		tic.Offsets.Minutes_6 = ((minute * 60) + second) % 360
		tic.Offsets.Minutes_10 = ((minute * 60) + second) % 600
		tic.Offsets.Minutes_12 = ((minute * 60) + second) % 720
		tic.Offsets.Minutes_15 = ((minute * 60) + second) % 900
		tic.Offsets.Minutes_20 = ((minute * 60) + second) % 1200
		tic.Offsets.Minutes_30 = ((minute * 60) + second) % 1800
		
        fmt.Printf("%4d %4d %4d %4d %4d\n", tic.Offsets.Minutes_1, tic.Offsets.Minutes_2,  tic.Offsets.Minutes_3, tic.Offsets.Minutes_4, tic.Offsets.Minutes_5,)
		
		soon := now.Add(1 * time.Second)
		soon_rounded := soon.Round(time.Second)
		wait_time := time.Until(soon_rounded)
        time.Sleep(wait_time)
    }
}