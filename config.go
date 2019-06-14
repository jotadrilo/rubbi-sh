package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	rubbishFolder = ""
	configFile    = ""
	configFolder  = ""
)

const (
	configFolderName = ".rubbish"
	configFileName   = "config.json"
)

type Folder struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type Config struct {
	Folders []Folder `json:"folders"`
	Latest  Folder   `json:"latest"`
	Root    string   `json:"root"`
}

func Initialize(root string) error {
	// Pre-configure folders
	if home, err := os.UserHomeDir(); err != nil {
		os.Exit(1)
	} else {
		configFolder = filepath.Join(home, configFolderName)
		if err := os.MkdirAll(configFolder, 0755); err != nil {
			return fmt.Errorf("unable to create directory tree: %+v", err)
		}
		configFile = filepath.Join(home, configFolderName, configFileName)
	}

	if _, err := os.Stat(configFile); err == nil {
		// Reload configuration
		config, err := load()
		if err != nil {
			return err
		}
		rubbishFolder = config.Root
		return nil
	} else if os.IsNotExist(err) {
		// Create configuration from scratch
		if err := Init(root); err != nil {
			return err
		}
	}

	return nil
}

func Init(root string) error {
	rubbishFolder = filepath.Join(root, "rubbish")
	if err := initConifg(); err != nil {
		return err
	}
	return nil
}

func initConifg() error {
	config := &Config{
		Folders: []Folder{},
		Latest:  Folder{},
		Root:    rubbishFolder,
	}
	if _, err := os.Create(configFile); err != nil {
		return fmt.Errorf("failed to create configuration file: %+v", err)
	}
	if err := save(config); err != nil {
		return err
	}
	return nil
}

func Clean() error {
	config, err := load()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(config.Root); err != nil {
		return fmt.Errorf("failed to remove the folder: %+v", err)
	}
	if err := initConifg(); err != nil {
		return err
	}
	return nil
}

func load() (*Config, error) {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %+v", err)
	}

	config := &Config{}
	if err := json.Unmarshal([]byte(b), config); err != nil {
		return nil, fmt.Errorf("failed to decode data: %+v", err)
	}

	return config, nil
}

func save(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode data: %+v", err)
	}
	if err := ioutil.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %+v", err)
	}
	return nil
}

func createFolder(folder Folder) error {
	if err := os.MkdirAll(folder.Path, 0755); err != nil {
		return fmt.Errorf("failed to create directory tree: %+v", err)
	}
	return nil
}

func doesFolderExist(folder Folder) bool {
	config, err := load()
	if err != nil {
		return false
	}

	for _, fol := range config.Folders {
		if fol.Name == folder.Name {
			return true
		}
	}
	return false
}

func GetFolder(name string) Folder {
	return Folder{
		Name: name,
		Path: filepath.Join(rubbishFolder, name),
	}
}

func AddFolder(name string) error {
	config, err := load()
	if err != nil {
		return err
	}

	folder := GetFolder(name)
	if doesFolderExist(folder) {
		return nil
	}

	if err := createFolder(folder); err != nil {
		return err
	}

	folders := append(config.Folders, folder)
	config.Folders = folders

	if err := save(config); err != nil {
		return err
	}

	if err := updateLatest(folder); err != nil {
		return err
	}

	return nil
}

func GetLatest() (*Folder, error) {
	config, err := load()
	if err != nil {
		return nil, err
	}

	return &config.Latest, nil
}

func Show() error {
	config, err := load()
	if err != nil {
		return err
	}

	for n, folder := range config.Folders {
		fmt.Printf("[%d] %s\t%s\n", n, folder.Name, folder.Path)
	}
	return nil
}

func Use(fn int) error {
	config, err := load()
	if err != nil {
		return err
	}
	if fn > len(config.Folders)-1 {
		return fmt.Errorf("the provided folder number does not match any existing folder")
	}
	config.Latest = config.Folders[fn]
	if err := save(config); err != nil {
		return err
	}
	return nil
}

func RemoveFolder(fn int) error {
	config, err := load()
	if err != nil {
		return err
	}

	targetFolder := config.Folders[fn]
	if err := os.RemoveAll(targetFolder.Path); err != nil {
		return fmt.Errorf("failed to remove the folder: %+v", err)
	}
	config.Folders = remove(config.Folders, fn)

	if err := save(config); err != nil {
		return err
	}

	// If we are removing the latest folder, point to the last folder in the list
	if config.Latest.Path == targetFolder.Path {
		updateLatest(config.Folders[len(config.Folders)-1])
	}
	return nil
}

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang/37335777#37335777
func remove(s []Folder, i int) []Folder {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func updateLatest(folder Folder) error {
	config, err := load()
	if err != nil {
		return err
	}

	config.Latest = folder
	if err := save(config); err != nil {
		return err
	}

	latestFolder := GetFolder("latest")

	// Remove existing symlink
	if _, err := os.Lstat(latestFolder.Path); err == nil {
		if err := os.Remove(latestFolder.Path); err != nil {
			return fmt.Errorf("failed to remove the current latest symlink: %+v", err)
		}
	}

	if err := os.Symlink(config.Latest.Name, latestFolder.Path); err != nil {
		return fmt.Errorf("failed to create the latest symlink: %+v", err)
	}

	return nil
}
