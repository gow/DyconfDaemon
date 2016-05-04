package main

import (
	"flag"
	"log"
	"os"

	"github.com/gow/DyconfDaemon/lib"
)

func main() {
	logger := log.New(os.Stdout, "[DyconfDaemon]: ", log.Lshortfile)
	fileName := getConfigFilePath()
	d, err := daemon.NewDaemon(
		daemon.OptionFileName(fileName),
		daemon.OptionLogger(logger),
		daemon.OptionHost("localhost"),
		daemon.OptionPort("9009"),
	)
	if err != nil {
		logger.Fatal(err)
	}

	if err := d.Start(); err != nil {
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
