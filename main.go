package main

import (
	"context"
	"fmt"
	"time"
)

var Version = "dev"

func main() {

	fmt.Print(
		`JAMCAT-MACH ` + Version + `
		
Jet Audio and Music Control Access Terminal-Media Access and Control Hub
		
JAMCAT-MACH made by @f45a - discord.
GIthub repo: https://github.com/angelfluffyookami/jamcat-mach
Please be aware, this utility will append ".bkp" to any mp3 files currently in the VTOL VR folder.
If closed properly without crashing or using TaskMan, will remove ".bkp" from said files and delete 
any temporary files placed in this folder.

I have not experienced any data loss using this, and there shouldn't be any to expect, however, please don't be stupid:
		- Don't touch the Player.log file while program is in use.
		- Do NOT modify, add, or remove any file in VTOLVR's RadioMusic folder while in use. 
			(I did add a way to handle this, however, I cannot guarantee it will work 100 percent of the time, so just d o n t)
			
------------------------------------------------------------------------------------------------------
Licensed under MPLv2
------------------------------------------------------------------------------------------------------

Log meaning:
Bearing X, the song VTOL VR thinks it's playing, where X equals the name of the mp3. E.g. 1.mp3 = Bearing 1

Angels X, the song VTOL VR thought it was playing before the current bearing, where X equals the name of the mp3. E.g. 0.mp3 = Angels 0

Interpreting this means if spike bearing goes backwards in regards to angels, then song will be rewinded. E.g. 0<-1<-2<-0
Vice versa is also true. E.g. 0->1->2->0

Splash 1, angels 0. Means player died, and the program is now expecting VTOL VR to start playing 0.mp3 again if player presses play button.

------------------------------------------------------------------------------------------------------
! ! JAMCAT-MACH is starting up, please make sure you are currently not spawned in an aircraft if VTOL VR is already running ! !
------------------------------------------------------------------------------------------------------


`)

	bkpPlayerMp3()
	InitMP3()

	go readLog()
	fmt.Println("JAMCAT-MACH is now listening to log events.")

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"revertBkp": func(ctx context.Context) error {
			revertBkp()
			return nil
		},
	})

	<-wait
}
