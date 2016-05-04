package main

/*

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	d.log.Println("Starting the daemon...")
	http.HandleFunc("/config/", d.configServer)
	http.HandleFunc("/config/getall/", d.getAllConfig)
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
	case "POST":
		d.postConfig(w, req)
	case "DELETE":
		d.deleteConfig(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (d *daemon) putConfig(w http.ResponseWriter, req *http.Request) {
	resp := json.NewEncoder(w)
	key := req.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp.Encode(struct{ Error string }{"Invalid Key"})
		return
	}
	val := req.FormValue("value")

	if _, err := d.cm.Get(key); err == nil {
		resp.Encode(struct{ Error string }{"key already exists. Use a POST request to modify it"})
		return
	}
	if err := d.cm.Set(key, []byte(val)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	resp.Encode(true)
}

func (d *daemon) getConfig(w http.ResponseWriter, req *http.Request) {
	keys := req.URL.Query()["key"]
	resp := json.NewEncoder(w)

	type kvpair struct {
		Key string `json:"key"`
		Val []byte `json:"value"`
		Err string `json:"error",omitempty`
	}
	var kvPairs []kvpair
	for _, key := range keys {
		pair := kvpair{Key: key}
		val, err := d.cm.Get(key)
		if err != nil {
			pair.Err = err.Error()
		} else {
			pair.Val = val
		}
		kvPairs = append(kvPairs, pair)
	}
	w.WriteHeader(http.StatusOK)
	resp.Encode(kvPairs)
}

func (d *daemon) getAllConfig(w http.ResponseWriter, req *http.Request) {
	if req.URL.EscapedPath() != "/config/getall/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	all, err := d.cm.Map()
	resp := json.NewEncoder(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	resp.Encode(all)

}
func (d *daemon) postConfig(w http.ResponseWriter, req *http.Request) {
	resp := json.NewEncoder(w)
	key := req.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp.Encode(struct{ Error string }{"Invalid Key"})
		return
	}
	val := req.FormValue("value")

	if _, err := d.cm.Get(key); err != nil {
		resp.Encode(struct{ Error string }{"key doesn't exists. Use a PUT request to add it."})
		return
	}
	if err := d.cm.Set(key, []byte(val)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	resp.Encode(true)
}

func (d *daemon) deleteConfig(w http.ResponseWriter, req *http.Request) {
	resp := json.NewEncoder(w)
	key := req.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp.Encode(struct{ Error string }{"Invalid Key"})
		return
	}

	if _, err := d.cm.Get(key); err != nil {
		resp.Encode(struct{ Error string }{"key doesn't exists."})
		return
	}

	if err := d.cm.Delete(key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	resp.Encode(true)
}
*/
