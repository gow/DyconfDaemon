package main

import (
	"encoding/json"
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
	http.HandleFunc("/config/", d.configServer)
	d.log.Fatal(
		"Failed to start the daemon.",
		http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil),
	)
	return nil
}

func (d *daemon) stop() error {
	return d.cm.Close()
}

func (d *daemon) configServer(w http.ResponseWriter, req *http.Request) {
	if req.URL.EscapedPath() != "/config/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch req.Method {
	case "PUT":
		d.putConfig(w, req)
	case "GET":
		d.getConfig(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (d *daemon) putConfig(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (d *daemon) getConfig(w http.ResponseWriter, req *http.Request) {
	keys := req.URL.Query()["key"]
	d.log.Print(keys)
	resp := json.NewEncoder(w)
	m, err := d.cm.Map()
	if err != nil {
		resp.Encode(err)
	}

	type kvpair struct {
		Key string `json:"key"`
		Val []byte `json:"value"`
		Err string `json:"error"`
	}
	var kvPairs []kvpair
	for _, key := range keys {
		pair := kvpair{Key: key}
		val, ok := m[key]
		if !ok {
			pair.Err = "Not Found"
		} else {
			pair.Val = val
		}
		kvPairs = append(kvPairs, pair)
	}
	w.WriteHeader(http.StatusOK)
	resp.Encode(kvPairs)
}

/*
type apiHandler interface {
	handleGet(w http.ResponseWriter, req *http.Request)
	handlePut(w http.ResponseWriter, req *http.Request)
	handlePost(w http.ResponseWriter, req *http.Request)
	handleDelete(w http.ResponseWriter, req *http.Request)
}

type api struct {
	h apiHandler
}

func (a *api) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		a.h.handleGet(w, req)
	case "PUT":
		a.h.handlePut(w, req)
	case "POST":
		a.h.handlePost(w, req)
	case "DELETE":
		a.h.handleDelete(w, req)
	}
}
*/
