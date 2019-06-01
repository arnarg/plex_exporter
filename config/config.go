package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"
)

type PlexConfig struct {
	Token string `json:"token"`
	UUID  string `json:"uuid"`
}

func Load(path string) (*PlexConfig, bool, error) {
	var plexConfig *PlexConfig

	// Get absolute path of config file
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, false, err
	}

	// Check if the file already exists
	_, err = os.Stat(absPath)
	if err != nil {
		// Config file doesn't exist so we generate new config
		if os.IsNotExist(err) {
			plexConfig = &PlexConfig{UUID: fmt.Sprintf("%s", uuid.NewV4())}
			return plexConfig, true, nil
		}

		return nil, false, err
	}

	// Read config file
	jsonString, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, false, err
	}

	// Parse JSON string
	plexConfig = &PlexConfig{}
	err = json.Unmarshal(jsonString, plexConfig)
	if err != nil {
		return nil, false, err
	}

	return plexConfig, false, nil
}

func Save(conf *PlexConfig, path string) error {
	targetDir := filepath.Dir(path)
	fileName := filepath.Base(path)

	// Check if base directory exists
	_, err := os.Stat(targetDir)
	if err != nil {
		// Attempt to create base directory
		if os.IsNotExist(err) {
			err = os.MkdirAll(targetDir, 0750)
			if err != nil {
				return err
			}
		}

		return err
	}

	// Create JSON string
	jsonString, err := json.Marshal(&conf)
	if err != nil {
		return err
	}

	// Write to file
	err = ioutil.WriteFile(filepath.Join(targetDir, fileName), jsonString, 0640)
	if err != nil {
		return err
	}

	return nil
}
