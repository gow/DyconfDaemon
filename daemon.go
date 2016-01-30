package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/gow/dyconf"
)

type daemon struct {
	fileName string
	cm       dyconf.ConfigManager
	log      *log.Logger
}

func (d *daemon) start(logger *log.Logger, fileName string, host string, port string) error {
	d.log = logger
	if fileName == "" {
		return fmt.Errorf("Invalid config file. File path cannot be empty.")
	}
	if !filepath.IsAbs(fileName) {
		return fmt.Errorf("Invalid config file: [%s]. File path must be absolute path.", fileName)
	}
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("Invalid config file: [%s]. The enclosing directory [%s] does not exit.", fileName, dir)
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		d.log.Printf("[INFO]: [%s] doesn't exist. Going to create a new one.", fileName)
	}

	// Open and close the config file.
	var err error
	d.cm, err = dyconf.NewManager(fileName)
	if err != nil {
		return err
	}
	d.log.Println("Existing map: ", spew.Sdump(d.cm.Map()))
	d.log.Println("Starting the daemon...")
	http.Handle("/", d)
	d.log.Fatal(
		"Failed to start the daemon.",
		http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil),
	)
	return nil
}

func (d *daemon) stop() error {
	return d.cm.Close()
}

func (d *daemon) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "PUT":
		d.handlePut(w, req)
	case "GET":
		d.handleGet(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (d *daemon) handlePut(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (d *daemon) handleGet(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
