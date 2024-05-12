package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type ConfigStruct struct {
	SFX    bool `json:"enablesfx"`
	Config int  `json:"configversion"`
}

var defaultConfig = ConfigStruct{
	SFX:    true,
	Config: 1,
}

var config = ConfigStruct{}

func readOrCreateConfig(filename string) (ConfigStruct, error) {
	var config ConfigStruct

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		config = defaultConfig
	} else if err != nil {
		return ConfigStruct{}, err
	} else {
		defer file.Close()

		bytes, err := io.ReadAll(file)
		if errors.Is(err, os.ErrNotExist) {
			config.SFX = defaultConfig.SFX
			config.Config = defaultConfig.Config

			bytes, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return ConfigStruct{}, err
			}

			if err := os.WriteFile(filename, bytes, 0644); err != nil {
				return ConfigStruct{}, err
			}
		}

		if err := json.Unmarshal(bytes, &config); err != nil {
			config.SFX = defaultConfig.SFX

			bytes, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return ConfigStruct{}, err
			}

			if err := os.WriteFile(filename, bytes, 0644); err != nil {
				return ConfigStruct{}, err
			}
		}

	}

	return config, nil
}
