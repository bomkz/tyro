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

	"github.com/udhos/equalfile"
)

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

// function that runs the bkp restore functions for all possible VTOL VR paths
func revertBkp() {

	paths := getVTOLDir()

	if paths == nil {
		log.Panic(fmt.Errorf("no VTOL VR Directories found. Is game installed?"))
	}

	for _, x := range paths {
		files, err := os.ReadDir(x + "RadioMusic\\")
		if err != nil {
			log.Panic(err)
		}

		for _, y := range files {
			checkRemoval(x, y.Name())
		}
		for _, y := range files {
			checkRestore(x, y.Name())
		}

	}
}

// Checks if there are any .bkp files that need to be restored, and restores them
func checkRestore(path string, name string) {
	if strings.Contains(name, ".bkp") {
		fmt.Println("Restoring file: " + name + ".")
		splitStr := strings.Split(name, ".bkp")
		if splitStr[0] == name {
			return
		}
		os.Rename(path+"RadioMusic\\"+name, path+"RadioMusic\\"+splitStr[0])
	}
}

// checks if mp3 file is the one jamcat-mach itself embedded in there, and if so, deletes it.
func checkRemoval(path string, name string) {
	if validRemoval(name) && compareMp3File(path, name) {
		fmt.Println("Cleaning up blank mp3 files.")
		os.Remove(path + "RadioMusic\\" + name)
	}
}

// checks if files fit name filter
func validRemoval(name string) bool {
	return name == "0.mp3" || name == "1.mp3" || name == "2.mp3"
}

// Returns whether file contents match up with embedded file or not.
func compareMp3File(path string, name string) bool {
	blank, err := embeds.Open("blank.mp3")
	if err != nil {
		log.Panic(err)
	}
	defer blank.Close()

	compare := equalfile.New(nil, equalfile.Options{})

	file, err := os.Open(path + "RadioMusic\\" + name)
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	equal, err := compare.CompareReader(blank, file)
	if err != nil {
		log.Panic(err)
	}

	if !equal {
		fmt.Println("Found " + name + " in directory, file contents do not match up, marked unsafe for removal.")

		return false
	}
	fmt.Println("Found " + name + " in directory, file contents match up, and marked safe for removle.")

	return true
}
