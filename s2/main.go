package main

import (
	"os"

	"github.com/dpatrie/sparkgrid/services"
)

func main() {
	toProcess := os.Getenv("FILE_DIR")
	if toProcess == "" {
		toProcess = "./requirements/"
	}
	doneProcessing := os.Getenv("PROCESSED_DIR")
	if doneProcessing == "" {
		doneProcessing = "./requirements/processed/"
	}

	//TOOD: Use inotify or something to see when new files are added
	(&services.S2{}).ProcessDir(toProcess, doneProcessing)
}
