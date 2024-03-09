package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

func tickTime() {
	for {
		tick <- true
		time.Sleep(100 * time.Millisecond)
	}
}

func readLog() {

	// Start log reading ticker
	go tickTime()

	/*
	*Set readlog to false,
	*prevents unwanted media inputs caused by log being read for first time
	*when player is already in game and has pressed buttons.
	 */
	readlog := false

	// Reads log every tick.
	for {
		<-tick

		// Reads log line by line
		logFile := readLogFile()
		scanLog := bufio.NewScanner(strings.NewReader(logFile))
		var logLinesTmp []string

		// Appends each log line to a variable.
		for scanLog.Scan() {
			logLinesTmp = append(logLinesTmp, scanLog.Text())
		}

		// Only takes action on log lines that are new, and ignores old ones.
		for y, x := range logLinesTmp {
			if y > (len(logLines) - 1) {
				if readlog {
					// Handles new log lines.
					logHandler(x)
				}

			}
		}
		logLines = logLinesTmp

		// Set to true after the first for loop passes
		readlog = true
	}
}

// Runs the appropriate function depending on the log line contents.
func logHandler(newline string) bool {
	if strings.Contains(newline, "Playing song:") {
		splitline := strings.SplitAfter(newline, "clip length")
		if strings.Contains(splitline[0], "0") {
			fmt.Println("Spike bearing 0, angels " + fmt.Sprint(currentTrack) + ".")
			track0(currentTrack)
		} else if strings.Contains(splitline[0], "1") {
			fmt.Println("Spike bearing 1, angels " + fmt.Sprint(currentTrack) + ".")
			track1(currentTrack)
		} else if strings.Contains(splitline[0], "2") {
			fmt.Println("Spike bearing 2, angels " + fmt.Sprint(currentTrack) + ".")
			track2(currentTrack)
		}
		return true
	} else if strings.Contains(newline, "FlightLogger:") && strings.Contains(newline, "has spawned.") {
		currentTrack = 0
		fmt.Println("Splash 1, bearing 0")
	}
	return false
}

// Selects the proper actions depending on current context
func track2(track int) {

	switch track {
	case 0:
		Track2{}.RW()

	case 1:
		Track2{}.FF()

	default:
		Track2{}.Play()
	}

}

// Selects the proper actions depending on current context
func track1(track int) {

	switch track {
	case 0:
		Track1{}.FF()
	case 2:
		Track1{}.RW()

	default:
		Track1{}.Play()
	}

}

// Selects the proper actions depending on current context
func track0(track int) {

	switch track {
	case 1:
		Track0{}.RW()

	case 2:
		Track0{}.FF()

	default:
		Track0{}.Play()
	}

}
