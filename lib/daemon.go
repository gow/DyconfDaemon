package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/facebookgo/stackerr"
	"github.com/gow/dyconf"
)

type daemon struct {
	fileName    string
	log         *log.Logger
	host        string
	port        string
	confManager dyconf.ConfigManager
}

type DaemonOption func(*daemon)

func OptionFileName(fileName string) DaemonOption {
	return func(d *daemon) {
		d.fileName = fileName
	}
}

func OptionLogger(logger *log.Logger) DaemonOption {
	return func(d *daemon) {
		d.log = logger
	}
}

func OptionHost(host string) DaemonOption {
	return func(d *daemon) {
		d.host = host
	}
}

func OptionPort(port string) DaemonOption {
	return func(d *daemon) {
		d.port = port
	}
}

func NewDaemon(options ...DaemonOption) (*daemon, error) {
	d := &daemon{}
	for _, option := range options {
		option(d)
	}

	if d.log == nil {
		return nil, stackerr.New("A logger must be provided")
	}

	if d.host == "" {
		return nil, stackerr.New("Host name/IP must be provided")
	}

	if d.port == "" {
		return nil, stackerr.New("Host port must be provided")
	}

	if err := d.validateFile(); err != nil {
		return nil, err
	}

	cm, err := dyconf.NewManager(d.fileName)
	if err != nil {
		return nil, err
	}
	d.confManager = cm

	return d, nil
}

func (d *daemon) validateFile() error {
	if d.fileName == "" {
		return stackerr.Newf("Invalid config file. File path cannot be empty.")
	}
	if !filepath.IsAbs(d.fileName) {
		return stackerr.Newf("Invalid config file: [%s]. File path must be absolute path.", d.fileName)
	}
	dir := filepath.Dir(d.fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return stackerr.Newf("Invalid config file: [%s]. The enclosing directory [%s] does not exit.", d.fileName, dir)
	}

	if _, err := os.Stat(d.fileName); os.IsNotExist(err) {
		d.log.Printf("[INFO]: [%s] doesn't exist. Going to create a new one.", d.fileName)
	}
	return nil
}

func (d *daemon) Start() error {
	d.log.Println("Starting the daemon...")
	http.HandleFunc("/config/", d.configAPI())
	//http.HandleFunc("/config/getall/", d.getAllConfig)

	d.log.Fatal(
		"Failed to start the daemon.",
		http.ListenAndServe(fmt.Sprintf("%s:%s", d.host, d.port), nil),
	)
	return nil
}

// configAPI returns the API handler.
func (d *daemon) configAPI() http.HandlerFunc {
	return API(&configAPI{dmn: d})
}

type api interface {
	GET(http.ResponseWriter, *http.Request)
	PUT(http.ResponseWriter, *http.Request)
	POST(http.ResponseWriter, *http.Request)
	DELETE(http.ResponseWriter, *http.Request)
}

type configAPI struct {
	dmn *daemon
}

func (c *configAPI) GET(w http.ResponseWriter, req *http.Request) {
	keys := req.URL.Query()["key"]
	resp := json.NewEncoder(w)
	type kvPair struct {
		Key   string `json:"key"`
		Value []byte `json:"value"`
		Err   string `json:"error,omitempty"`
	}
	var kvp []kvPair
	for _, key := range keys {
		pair := kvPair{Key: key}
		val, err := c.dmn.confManager.Get(key)
		if err != nil {
			pair.Err = err.Error()
		} else {
			pair.Value = val
		}
		kvp = append(kvp, pair)
	}

	w.WriteHeader(http.StatusOK)
	resp.Encode(kvp)
}

func (c *configAPI) PUT(w http.ResponseWriter, req *http.Request) {
	resp := json.NewEncoder(w)
	key := req.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp.Encode(struct{ Error string }{"Invalid Key"})
		return
	}
	if _, err := c.dmn.confManager.Get(key); err == nil {
		resp.Encode(struct{ Error string }{"key already exists. Use a POST request to modify it"})
		return
	}

	val := req.FormValue("value")
	if err := c.dmn.confManager.Set(key, []byte(val)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp.Encode(true)
}

func (c *configAPI) DELETE(w http.ResponseWriter, req *http.Request) {
	resp := json.NewEncoder(w)
	key := req.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp.Encode(struct{ Error string }{"Invalid Key"})
		return
	}

	if _, err := c.dmn.confManager.Get(key); err != nil {
		resp.Encode(struct{ Error string }{"key doesn't exists."})
		return
	}

	if err := c.dmn.confManager.Delete(key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp.Encode(true)
}

func (c *configAPI) POST(w http.ResponseWriter, req *http.Request) {
	resp := json.NewEncoder(w)
	key := req.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp.Encode(struct{ Error string }{"Invalid Key"})
		return
	}

	if _, err := c.dmn.confManager.Get(key); err != nil {
		resp.Encode(struct{ Error string }{"key doesn't exists. Use PUT to add a new key."})
		return
	}

	val := req.FormValue("value")
	if err := c.dmn.confManager.Set(key, []byte(val)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp.Encode(true)
}

func API(a api) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "PUT":
			a.PUT(w, req)
		case "GET":
			a.GET(w, req)
		case "POST":
			a.POST(w, req)
		case "DELETE":
			a.DELETE(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
