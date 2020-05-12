// Copyright 2020 codeM GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codeM-code/nodedb"
)

// Utility functions common to the experiments

// ExecutableName return the name of the running executable
func ExecutableName() (string, error) {
	name, err := os.Executable()
	if err != nil {
		return "", err
	}
	name = filepath.Base(name)
	return strings.TrimSuffix(name, filepath.Ext(name)), nil
}

//LogName returns the name of the running executable combind with suffix .log
func LogName() (string, error) {
	logname, err := ExecutableName()
	return logname + ".log", err
}

// DBName returns the name of the running executable combind with suffix .db
func DBName() (string, error) {
	dbname, err := ExecutableName()
	return dbname + ".db", err
}

// RemoveDB deletes the file with the provided file name
func RemoveDB(dbname string) error {
	return os.RemoveAll(dbname)
}

// OpenFreshDB removes an existing '<executable name>.db' file and create a new one.
// It returns an open nodedb.Connection. Use defer conn.Close() to perform close and clean up when main function terminates.
func OpenFreshDB() (*nodedb.Connection, error) {
	dbname, err := DBName()
	if err != nil {
		log.Fatalf("error db file name: %v", err)
		return nil, err
	}
	err = RemoveDB(dbname)
	if err != nil {
		log.Fatalf("error deleting db file: %v", err)
		return nil, err
	}
	return nodedb.Open(dbname)
}

// StartLogFile opens (or creates) a log file name after the running executable
func StartLogFile() (*os.File, error) {

	filename, err := LogName()

	if err != nil {
		log.Fatalf("error log file name: %v", err)
		return nil, err
	}

	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
		return nil, err
	}

	log.SetOutput(logFile)
	return logFile, err
}

// StartTimer provides a quick and easy timer, just pass a name and call the returned function when done to stop the timer
func StartTimer(name string) func() {
	start := time.Now()
	log.Println(name, "started")
	return func() {
		log.Println(name, "completed after", time.Now().Sub(start))
	}
}
