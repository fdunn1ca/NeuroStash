// Program to stash neurofeedback client files with a date and time in the file name.
// Author: Frank Dunn
// Date started: 2018-03-26

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/beevik/etree"
	times "gopkg.in/djherbis/times.v1"
)

const sleepMinutes = 1

type lastSessionT struct {
	sessionNum int
	fileName   string
}
type lastSessionsT map[string]lastSessionT

func main() {
	sessionDir := filepath.Join("/", "home", "frank", "go", "src", "github.com", "fdunn1ca", "neuroClientStash", "Session")
	// sessionDir := filepath.Join("D:", "Session")
	lastSessions := make(lastSessionsT)

	// loop through the archived files and get the highest session number for each client
	files, err := ioutil.ReadDir(sessionDir)
	checkErr("", err)
	for _, file := range files {
		if file.IsDir() {
			client := file.Name()
			sessions, err := filepath.Glob(filepath.Join(sessionDir, client, "archive", "*.xml"))
			checkErr("", err)

			var maxSession lastSessionT
			for _, archiveFile := range sessions {
				session := getSession(archiveFile)
				checkErr("failed to convert to integer", err)
				if session > maxSession.sessionNum {
					maxSession.sessionNum = session
					maxSession.fileName = archiveFile
				}
			}
			if maxSession.sessionNum != 0 {
				lastSessions[client] = maxSession
			}
		}
	}
	for client, session := range lastSessions {
		fmt.Printf("%10s %2d %s\n", client, session.sessionNum, session.fileName)
	}

	// loop through the files in the Session directory
	// forever
	files, err = ioutil.ReadDir(sessionDir)
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

				// stash a date stamped copyFile in the archive directory
				stash(clientId, sessionDir, lastSessions)
			}
		}

		// wait for a while then  check again
		time.Sleep(sleepMinutes * time.Minute)
	}
}
//ToDo comment
// more of the
//not part of the
func checkErr(message string, err error) {
	if err != nil {
		log.Fatal(message+" ", err)
	}

}

func stash(clientId, sessionDir string, lastSessions lastSessionsT) {
	oldLocation := filepath.Join(sessionDir, clientId, clientId+".xml")
	lastSession := getSession(oldLocation)
	// fmt.Printf("%7s %5s %16s %3d %11s %s\n", "client:", clientId, "session number:", lastSession, "file name:", oldLocation)
	if lastSessions[clientId].sessionNum == lastSession {
		fmt.Printf("Removing: %s\n", lastSessions[clientId].fileName)
		err := os.Remove(lastSessions[clientId].fileName)
		checkErr(lastSessions[clientId].fileName, err)
	}

	// get the file change date-time
	// to use for the archived file name
	fileTimes, err := times.Stat(oldLocation)
	checkErr("", err)
	changeTime := fileTimes.ModTime()
	timeStamp := changeTime.Format("2006-01-02_15.04.05")

	// create a time stamped file name
	newName := timeStamp + "_client_" + clientId + ".xml"
	newLocation := filepath.Join(sessionDir, clientId, "archive", newName)

	// If the archived file does not already exist
	// copyFile the client file to the archive directory with the new name
	if _, err := os.Stat(newLocation); os.IsNotExist(err) {
		err = copyFile(oldLocation, newLocation)
		checkErr("", err)
	}
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copyFile file attributes.
func copyFile(src, dst string) error {
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

func getSession(fileName string) int {
	doc := etree.NewDocument()
	err := doc.ReadFromFile(fileName)
	checkErr("Can't parse xml", err)
	result, err := strconv.Atoi(doc.SelectElement("Configuration").SelectElement("sessionnumber").Text())
	checkErr("failed to convert session number to integer", err)
	return result
}
