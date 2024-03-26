package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nxadm/tail"
)

func tailLogFile() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	t, err := tail.TailFile(home+"\\AppData\\LocalLow\\Boundless Dynamics, LLC\\VTOLVR\\Player-prev.log", tail.Config{Follow: true, ReOpen: true, Poll: true})
	if err != nil {
		log.Fatal(err)
	}

	for x := range t.Lines {

		logHandler(strings.TrimSuffix(x.Text, "\r"))
	}

}

// Runs the appropriate function depending on the log line contents.
func logHandler(newline string) bool {
	if strings.Contains(newline, "Attempting to join lobby") || strings.Contains(newline, "Creating a lobby for ") {
		var host bool

		if strings.Contains(newline, "Creating a lobby for ") {
			host = true
		}
		go onLobbyJoin(host, false, LobbyStruct{})
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
