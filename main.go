package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const usage string = ` Usage: clrd [options]
 - clrd: Move the current content of the Downloads directory to a tmp path (defaults to $HOME/.clrd)
 - clrd -d: Remove the contents saved in the tmp path (defaults to $HOME/.clrd)
`

func main() {
	flag.Usage = func() {
		fmt.Printf(usage)
	}
	del := flag.Bool("d", false, "specify to delete all saved download entries")
	flag.Parse()

	// ensure relevant paths
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("couldn't get user's home directory. set $HOME env variable: %s\n", err.Error())
	}

	clrdPath := os.Getenv("CLRD_PATH")
	if clrdPath == "" {
		clrdPath = filepath.Join(home, ".clrd")
	}

	// purge if desired
	if *del {
		err := doPurge(clrdPath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if _, err := os.Stat(clrdPath); os.IsNotExist(err) {
		err := os.MkdirAll(clrdPath, os.ModePerm)
		if err != nil {
			log.Fatalf("couldn't create the directory %s: %s\n", clrdPath, err.Error())
		}
	}

	downloadPath := os.Getenv("Downloads")
	var files []os.DirEntry
	if downloadPath == "" {
		downloadPath = filepath.Join(home, "Downloads")
		if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
			log.Fatalf("directory %s not found. set the $Downloads env variable\n", downloadPath)
		}
		files, err = os.ReadDir(downloadPath)
		if err != nil {
			log.Fatalf("error when listing %s: %s\n", downloadPath, err.Error())
		}
		if len(files) == 0 {
			fmt.Println("nothing to clear")
			return
		}

	}

	// mv the files
	t := time.Now()
	tStr := t.Format("2006-01-02 15:04:05")
	savePath := filepath.Join(clrdPath, tStr)
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		log.Fatalf("couldn't create the directory %s: %s\n", savePath, err.Error())
	}

	for _, file := range files {
		oldPath := filepath.Join(downloadPath, file.Name())
		newPath := filepath.Join(savePath, file.Name())
		err := os.Rename(oldPath, newPath)
		if err != nil {
			log.Fatalf("couldn't move %s to %s: %s\n", oldPath, newPath, err.Error())
		}
	}
	fmt.Printf("moved data to %s\n", savePath)
}

func doPurge(clrdPath string) error {
	files, err := os.ReadDir(clrdPath)
	if err != nil {
		return fmt.Errorf("couldn't read %s: %s\n", clrdPath, err.Error())
	}

	if len(files) == 0 {
		return nil
	}

	for _, f := range files {
		if err := os.RemoveAll(filepath.Join(clrdPath, f.Name())); err != nil {
			return fmt.Errorf("couldn't remove %s: %s\n", filepath.Join(clrdPath, f.Name()), err.Error())
		}

	}
	fmt.Printf("cleared the content of %s\n", clrdPath)
	return nil
}
