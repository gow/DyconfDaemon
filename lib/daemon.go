package daemon

import (
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
	//http.HandleFunc("/config/", d.configServer)
	//http.HandleFunc("/config/getall/", d.getAllConfig)

	d.log.Fatal(
		"Failed to start the daemon.",
		http.ListenAndServe(fmt.Sprintf("%s:%s", d.host, d.port), nil),
	)
	return nil
}
