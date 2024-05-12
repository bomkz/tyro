package main

import (
	"errors"
	"fmt"
	"os"
)

func ensureSFXExists() {

	fmt.Println("Naking sure soundAssets exist!")
	if _, err := os.Stat("soundAssets"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("soundAssets folder not found, creating!")
		os.Mkdir("soundAssets", os.ModePerm)
	}
	if _, err := os.Stat("soundAssets/humiliation.wav"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Humiliation SFX does not exist, creating!")
		humiliation, err := embeds.ReadFile("soundAssets/humiliation.wav")
		if err != nil{
			fmt.Println(err)
		}
		err = os.WriteFile("soundAssets/humiliation.wav", humiliation, 0777)
		if err != nil {
			fmt.Println(err)
		}
	}
}
