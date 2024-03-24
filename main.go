package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	updatechecker "github.com/Christian1984/go-update-checker"
)

var Version = "1.1.1"

func main() {
	fmt.Print(info)
	uc := updatechecker.New("angelfluffyookami", "tyro", "Tyro", "https://github.com/angelfluffyookami/tyro/releases/latest", 0, true)
	uc.CheckForUpdate(Version)
	uc.PrintMessage()
	// Warns user if they opened tyro.exe before opening VTOL VR
	ensureLogFileNew()

	// Starts the rich presence updater and sets status to idle
	go startRP()
	idle()

	// Starts the log reader.
	//go readLog()

	go tailLogFile()

	waiting := gracefulShutdown(context.Background(), 30*time.Second, map[string]operation{
		"writefile": func(ctx context.Context) error {

			fmt.Println("Gracefully flushing data to file... Please hold on.")
			Message <- "LeaveLobby()"
			time.Sleep(5 * time.Second)
			exportJson()
			return nil
		},
	})

	<-waiting

	fmt.Println("JAMCAT-MACH is now listening to log events.")
}
func taintFile() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}

	file, _ := os.OpenFile(home+"\\AppData\\LocalLow\\Boundless Dynamics, LLC\\VTOLVR\\Player.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	// Append some data
	if _, err := file.Write([]byte("This is a taint mark left by github.com/angelfluffyookami/tyro to ensure the program only reads new log files. Ignore this comment, as it does not modify this file, or the game, in any other way.")); err != nil {
		log.Fatal(err)
	}
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

func ensureLogFileNew() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}

	file, _ := os.ReadFile(home + "\\AppData\\LocalLow\\Boundless Dynamics, LLC\\VTOLVR\\Player.log")

	if strings.Contains(string(file), "This is a taint mark left by github.com/angelfluffyookami/tyro to ensure the program only reads new log files. Ignore this comment, as it does not modify this file, or the game, in any other way.") {
		fmt.Println(`
	
	
	
	
	
	


Warning: Old Player.log file has been detected, please make sure you run tyro.exe after you open VTOL VR, otherwise, weird behaviour may arise.
If you reopened tyro.exe after a crash or accidentally closing it, don't worry, it shouldn't experience any bugs.`)

	} else {
		taintFile()
	}
}
