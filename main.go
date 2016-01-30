package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "[DyconfDaemon]: ", log.Lshortfile)
	fileName := getConfigFilePath()
	d := &daemon{}
	if err := d.start(logger, fileName, "localhost", "9009"); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Daemon successfully started")
}

func getConfigFilePath() string {
	var fileName string
	flag.StringVar(&fileName, "file", "", "The full file path of the config file.")
	flag.Parse()
	return fileName
}
