// Program to stash neurofeedback client files with a date and time in the file name.
// Author: Frank Dunn
// Date started: 2018-03-26

package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/djherbis/times.v1"
)

const sleepMinutes = 1

func main() {
	// sessionDir := filepath.Join("/", "home", "frank", "go", "src", "github.com", "fdunn1ca", "neuroClientStash", "Session")
	sessionDir := filepath.Join("D:", "Session")

	// loop through the files in the Session directory
	// forever
	files, err := ioutil.ReadDir(sessionDir)
	checkErr("", err)
	for {
		for _, file := range files {
			if file.IsDir() {
				clientId := file.Name()
				archiveDir := filepath.Join(sessionDir, clientId, "archive")

				// make the archive directory if it doesn't exist
				if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
					os.Mkdir(archiveDir, 0777)
				}

				// stash a date stamped copy in the archive directory
				stash(clientId, sessionDir)
			}
		}

		// wait for a while then  check again
		time.Sleep(sleepMinutes * time.Minute)
	}
}

func checkErr(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func stash(clientId, sessionDir string) {
	oldLocation := filepath.Join(sessionDir, clientId, clientId+".xml")

	// get the file change date-time
	// to use for the archived file name
	fileTimes, err := times.Stat(oldLocation)
	checkErr("", err)
	changeTime := fileTimes.ChangeTime()
	timeStamp := changeTime.Format("2006-02-01_15.04.05")

	// create a time stamped file name
	newName := timeStamp + "_client_" + clientId + ".xml"
	newLocation := filepath.Join(sessionDir, clientId, "archive", newName)

	// If the archived file does not already exist
	// copy the client file to the archive directory with the new name
	if _, err := os.Stat(newLocation); os.IsNotExist(err) {
		err = copy(oldLocation, newLocation)
		checkErr("", err)
	}
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return err
	}
	return dest.Close()
}
