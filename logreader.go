package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func tickTime() {
	for {
		tick <- true
		time.Sleep(0 * time.Millisecond)
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
	readlog := true

	// Reads log every tick.
	for {
		<-tick

		// Reads log line by line
		logFile := readLogFile()

		// Appends each log line to a variable.
		logLinesTmp := strings.Split(logFile, "\r\n")

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
	if strings.Contains(newline, "Attempting to join lobby") || strings.Contains(newline, "Creating a lobby for 16 players") {
		var host bool
		if strings.Contains(newline, "Creating a lobby for 16 players") {
			host = true
		}
		go onLobbyJoin(host)
		InLobby = true
		Message <- newline
		<-done

	} else if InLobby {
		Message <- newline
		<-done
	}

	if strings.Contains(newline, "Set current pilot to ") {
		currentPilot = strings.TrimPrefix(newline, "Set current pilot to ")
	}
	return true

}

var InLobby bool

func UpdateStatus(statusType int) {
	fmt.Println("")
}

// reads the log file.
func readLogFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	file, err := os.ReadFile(home + "\\AppData\\LocalLow\\Boundless Dynamics, LLC\\VTOLVR\\Player.log")
	if err != nil {
		log.Panic(err)
	}

	return string(file)
}
