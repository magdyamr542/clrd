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
 - clrd -purge: Remove the contents saved in the tmp path (defaults to $HOME/.clrd)
`

func main() {
	flag.Usage = func() {
		fmt.Printf(usage)
	}
	purge := flag.Bool("purge", false, "specify to clear all saved download entries")
	flag.Parse()

	// ensure relevant paths
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("couldn't get user's home directory. set $HOME env variable")
	}

	clrdPath := os.Getenv("CLRD_PATH")
	if clrdPath == "" {
		clrdPath = filepath.Join(home, ".clrd")
	}

	// purge if desired
	if *purge {
		if err := os.RemoveAll(clrdPath); err != nil {
			log.Fatalf("couldn't clear from %s\n", clrdPath)
		}
		fmt.Printf("cleared the content of %s\n", clrdPath)
		return
	}

	if _, err := os.Stat(clrdPath); os.IsNotExist(err) {
		err := os.MkdirAll(clrdPath, os.ModePerm)
		if err != nil {
			log.Fatalf("couldn't create the directory %s\n", clrdPath)
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
			log.Fatalf("error when listing %s\n", downloadPath)
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
		log.Fatalf("couldn't create the directory %s\n", savePath)
	}

	for _, file := range files {
		oldPath := filepath.Join(downloadPath, file.Name())
		newPath := filepath.Join(savePath, file.Name())
		err := os.Rename(oldPath, newPath)
		if err != nil {
			log.Fatalf("couldn't move %s to %s", oldPath, newPath)
		}
	}
	fmt.Printf("moved data to %s\n", savePath)
}
