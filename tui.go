package main

import (
	"github.com/juju/errors"
	"github.com/manifoldco/promptui"
)

// SelectFolder prompts a selection menu and returns the choosen folder number
func SelectFolder(config *Config) (int, error) {
	folders := make(map[string]int)
	items := make([]string, 0, len(folders))

	for fn, folder := range config.Folders {
		folders[folder.Path] = fn
		items = append(items, folder.Path)
	}

	prompt := promptui.Select{
		Label: "Select Folder",
		Items: items,
	}

	_, fn, err := prompt.Run()
	if err != nil {
		return -1, errors.Errorf("failed to prompt folders selection menu: %+v", err)
	}
	return folders[fn], nil
}
