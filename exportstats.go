package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func exportJson() {
	beautifulJSON, err := json.MarshalIndent(LobbyHistory, "", "    ") // Use four spaces for indentation
	if err != nil {
		panic(err)
	}
	timestamp := time.Now()
	err = os.Mkdir(".\\vtolvrdata", os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(".\\vtolvrdata\\vtolvr-"+fmt.Sprint(timestamp.Unix())+".json", beautifulJSON, 0644)
	if err != nil {
		panic(err)

	}

}
