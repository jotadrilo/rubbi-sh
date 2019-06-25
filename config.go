package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/juju/errors"
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

// Folder represents a folder object
type Folder struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

// Config represents a configuration file object
type Config struct {
	Folders []Folder `json:"folders"`
	Latest  Folder   `json:"latest"`
	Root    string   `json:"root"`
}

// Initialize creates the required folder structure and an initial configuration file
func Initialize(root string) error {
	// Pre-configure folders
	if home, err := os.UserHomeDir(); err != nil {
		os.Exit(1)
	} else {
		configFolder = filepath.Join(home, configFolderName)
		if err := os.MkdirAll(configFolder, 0755); err != nil {
			return errors.Errorf("unable to create directory tree: %+v", err)
		}
		configFile = filepath.Join(home, configFolderName, configFileName)
	}

	if _, err := os.Stat(configFile); err == nil {
		// Reload configuration
		config, err := Load()
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

// Init initializes the configuration file
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
		return errors.Errorf("failed to create configuration file: %+v", err)
	}
	if err := config.Save(); err != nil {
		return err
	}
	return nil
}

// Clean will remove the root rubbish folder and will recreate the configuration
func (config *Config) Clean() error {
	if err := os.RemoveAll(config.Root); err != nil {
		return errors.Errorf("failed to remove the folder: %+v", err)
	}
	if err := initConifg(); err != nil {
		return err
	}
	return nil
}

// Load returns the configuration parsed from the configuration file
func Load() (*Config, error) {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.Errorf("failed to read file: %+v", err)
	}

	config := &Config{}
	if err := json.Unmarshal([]byte(b), config); err != nil {
		return nil, errors.Errorf("failed to decode data: %+v", err)
	}

	return config, nil
}

// Save dumps the current status of the configuration to disk
func (config *Config) Save() error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.Errorf("failed to encode data: %+v", err)
	}
	if err := ioutil.WriteFile(configFile, data, 0644); err != nil {
		return errors.Errorf("failed to write file: %+v", err)
	}
	return nil
}

func createFolder(folder Folder) error {
	if err := os.MkdirAll(folder.Path, 0755); err != nil {
		return errors.Errorf("failed to create directory tree: %+v", err)
	}
	return nil
}

func (config *Config) doesFolderExist(folder Folder) bool {
	for _, fol := range config.Folders {
		if fol.Name == folder.Name {
			return true
		}
	}
	return false
}

// GetFolder returns a cannonical folder for the given name
func GetFolder(name string) Folder {
	return Folder{
		Name: name,
		Path: filepath.Join(rubbishFolder, name),
	}
}

// AddFolder adds a folder to the configuration and becomes it the latest folder
func (config *Config) AddFolder(name string) error {
	folder := GetFolder(name)
	if config.doesFolderExist(folder) {
		return nil
	}

	if err := createFolder(folder); err != nil {
		return err
	}

	config.Folders = append(config.Folders, folder)
	sortFolders(config.Folders)

	if err := config.updateLatest(folder); err != nil {
		return err
	}

	return nil
}

// Show prints all the folder entries
func (config *Config) Show() error {
	for n, folder := range config.Folders {
		fmt.Printf("[%d] %s\t%s\n", n, folder.Name, folder.Path)
	}
	return nil
}

// Use changes the latest folder to the provided folder number
func (config *Config) Use(fn int) error {
	if fn > len(config.Folders)-1 {
		return errors.Errorf("the provided folder number does not match any existing folder")
	}
	if err := config.updateLatest(config.Folders[fn]); err != nil {
		return err
	}
	return nil
}

// RemoveFolder removes a folder from the configuration and the filesystem
// If the folder is marked as latest, the last folder entry becomes latest.
func (config *Config) RemoveFolder(fn int) error {
	targetFolder := config.Folders[fn]
	if err := os.RemoveAll(targetFolder.Path); err != nil {
		return errors.Errorf("failed to remove the folder: %+v", err)
	}
	config.Folders = remove(config.Folders, fn)
	sortFolders(config.Folders)

	// If we are removing the latest folder, point to the last folder in the list
	if config.Latest.Path == targetFolder.Path {
		config.updateLatest(config.Folders[len(config.Folders)-1])
	}
	return nil
}

// Flush iterates over the folder entries and remove them from the configuration if they are not present anymore.
func (config *Config) Flush() (errs error) {
	var folders = []Folder{}
	for _, fol := range config.Folders {
		if _, err := os.Stat(fol.Path); err == nil {
			folders = append(folders, fol)
		}
	}

	config.Folders = folders
	return errors.Trace(errs)
}

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang/37335777#37335777
func remove(s []Folder, i int) []Folder {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (config *Config) updateLatest(folder Folder) error {
	config.Latest = folder
	latestFolder := GetFolder("latest")

	// Remove existing symlink
	if _, err := os.Lstat(latestFolder.Path); err == nil {
		if err := os.Remove(latestFolder.Path); err != nil {
			return errors.Errorf("failed to remove the current latest symlink: %+v", err)
		}
	}

	if err := os.Symlink(config.Latest.Name, latestFolder.Path); err != nil {
		return errors.Errorf("failed to create the latest symlink: %+v", err)
	}

	return nil
}

func sortFolders(folders []Folder) {
	sort.SliceStable(folders, func(i int, j int) bool {
		return folders[i].Name < folders[j].Name
	})
}
