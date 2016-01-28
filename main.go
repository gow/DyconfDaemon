package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fileName := getConfigFilePath()
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Printf("[INFO]: [%s] doesn't exist. Going to create a new one.", fileName)
	}

	fmt.Println(fileName)
}

func getConfigFilePath() string {
	var fileName string
	flag.StringVar(&fileName, "file", "", "The full file path of the config file.")
	flag.Parse()
	if fileName == "" {
		log.Fatal("Invalid config file. File path cannot be empty.")
	}
	if !filepath.IsAbs(fileName) {
		log.Fatalf("Invalid config file: [%s]. File path must be absolute path.", fileName)
	}
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("Invalid config file: [%s]. The enclosing directory [%s] does not exit.", fileName, dir)
	}

	return fileName
}
