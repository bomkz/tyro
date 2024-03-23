package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func exportJson() {
	if saveoutput {
		beautifulJSON, err := json.MarshalIndent(LobbyHistory, "", "    ") // Use four spaces for indentation
		if err != nil {
			panic(err)
		}
		timestamp := time.Now()
		err = os.WriteFile(".\\vtolvr-"+fmt.Sprint(timestamp.Unix())+".json", beautifulJSON, 0644)
		if err != nil {
			panic(err)
		}
	}

}
