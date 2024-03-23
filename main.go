package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var Version = "dev"

func main() {

	fmt.Print(
		`TYRO ` + Version + `
Telemetry Yield Real-time Observations
Utility to display discord rich presence, and dump your statistics afterwards.
Credits to:

@joespeed52 - F/A-26B Photo on the default rich presence token.
@toast2812 - EF-24G and F-45A Photo on the default rich presence token.
@kentuckyfrieda10wallsimper - AH-94 Photo on the default rich presence token.
@romanian_wallet_inspector - A/V-42C Photo on the default rich presence token. 
@dubyaaa - T-55 Photo on the default rich presence token.

------------------------------------------------------------------------------------------------------
Licensed under MIT
------------------------------------------------------------------------------------------------------

Upon closure, the program will save your entire gameplay session statistics to a timestamp tagged json file,
You can delete or keeep it, however, in the future, I have some projects planned where you can use these files
and output graphs of various varieties as you choose, and also view mission summaries in an easily readable way.
This file is filled with a lot of useful information.
------------------------------------------------------------------------------------------------------
! ! JAMCAT-MACH is starting up, please make sure you are currently not spawned in an aircraft if VTOL VR is already running ! !
------------------------------------------------------------------------------------------------------


`)

	richPresence()
	time.Sleep(2 * time.Second)
	idle()

	go readLog()

	waiting := gracefulShutdown(context.Background(), 30*time.Second, map[string]operation{
		"writefile": func(ctx context.Context) error {
			exportJson()
			return nil
		},
	})

	<-waiting

	fmt.Println("JAMCAT-MACH is now listening to log events.")
}

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Println("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}

type operation func(ctx context.Context) error
