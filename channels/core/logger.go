package core

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func SetUpLogger() {

	// Make sure to log what function it error happened
	log.SetReportCaller(true)
	// What kind of level of logging we want
	log.SetLevel(log.DebugLevel)

	// Create local output file
	f, err := os.Create("./logs.txt")

	if err != nil {
		log.Fatal(err)
	}

	// Say logger to write to it
	log.SetOutput(f)
}
