package main

import (
	"encoding/json"
	"os"
)

func exportJson() {
	beautifulJSON, err := json.MarshalIndent(LobbyHistory, "", "    ") // Use four spaces for indentation
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(".\\output.json", beautifulJSON, 0644)
	if err != nil {
		panic(err)
	}
}
