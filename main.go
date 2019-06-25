package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/juju/errors"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	clean = flag.Bool("clean", false, "if true, the rubbish folder is removed")
	show  = flag.Bool("show", false, "if true, outputs the current rubbish folders")
	ver   = flag.Bool("ver", false, "if true, the rubbish version will be shown")
	add   = flag.String("add", "", "folder name to add")
	del   = flag.String("del", "", "folder number to delete")
	root  = flag.String("root", "/tmp", "temporary location for the rubbish folder")
	use   = flag.String("use", "", "folder number to use")
)

func init() {
	flag.Usage = usage
	if err := Initialize(*root); err != nil {
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Print("Release information:\n")
	fmt.Printf("  - Version:\t%s\n", version)
	fmt.Printf("  - Commit: \t%s\n", commit)
	fmt.Printf("  - Date:   \t%s\n", date)
	fmt.Print("\nOptions:\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if *ver {
		fmt.Printf("%s\n", version)
		return nil
	}

	config, err := Load()
	if err != nil {
		return err
	}

	if *show {
		if err := config.Show(); err != nil {
			return err
		}
		return nil
	}

	if *use != "" {
		fn, err := strconv.Atoi(*use)
		if err != nil {
			return errors.Errorf("failed to parse folder number to use: %+v", err)
		}
		if err := config.Use(fn); err != nil {
			return err
		}
		config.Save()
		return nil
	}

	if *clean {
		config.Clean()
		Init(*root)
		return nil
	}

	if *add != "" {
		if err := config.AddFolder(*add); err != nil {
			return err
		}
	} else {
		timestamp := time.Now().Format("20060102")
		if err := config.AddFolder(timestamp); err != nil {
			return err
		}
	}

	if *del != "" {
		fn, err := strconv.Atoi(*del)
		if err != nil {
			return errors.Errorf("failed to parse folder number to delete: %+v", err)
		}
		if err := config.RemoveFolder(fn); err != nil {
			return err
		}
		config.Save()
		return nil
	}

	// There possibily was a restart and the `tmp` folder was removed.
	if err := config.Flush(); err != nil {
		return err
	}

	// Try to change to the target directory
	if err := os.Chdir(config.Latest.Path); err != nil {
		return errors.Errorf("failed to change to directory: %+v", err)
	}

	// Dump all the changes in the done configuration
	config.Save()

	// Print path as we cannot change the shell working directory
	// from an external binary
	fmt.Printf(config.Latest.Path)

	return nil
}
