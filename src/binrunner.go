/*
 * @file : binrunner.go
 * @brief : binrunner, start, kill, list, processes via http
 *			interface.
 * @author : Manuel A. Rodriguez (manuel.rdrs@gmail.com)
 */
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

/*
 * Config file for service
 */
type Config struct {
	HttpListenPort string
	MaxProcesses   int
}

type ProcConfig struct {
	Entries []EntryConfig
}

type EntryConfig struct {
	Alias string
	Cmd   string
	Args  string
}

/*
 * Global vars
 */
var (
	//
	_cfg_dir       = flag.String("cfg_dir", "../cfg", "path to config directory")
	_cfg_path      = flag.String("cfg_path", *_cfg_dir+"/binrunner.cfg", "path to config file")
	_cfg_proc_path = flag.String("cfg_proc_path", *_cfg_dir+"/proc.cfg", "path to proc config file")
	//
	_cfg      Config
	_cfg_proc ProcConfig
)

// --------------------------------------------------------------
func handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

	} else {
		// 400
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Method"))
	}
}

// --------------------------------------------------------------
func handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

	} else {
		// 400
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Method"))
	}
}

// --------------------------------------------------------------
func handleRestart(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

	} else {
		// 400
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Method"))
	}
}

// --------------------------------------------------------------
func handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else {
		// 400
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Method"))
	}
}

// --------------------------------------------------------------
func handleProfiles(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else {
		// 400
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Method"))
	}
}

// --------------------------------------------------------------
func loadConfig(cfg_path string) Config {
	file, errOpen := os.Open(cfg_path)
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	decoder := json.NewDecoder(file)

	cfg := Config{}
	errDecode := decoder.Decode(&cfg)
	if errDecode != nil {
		log.Fatal(errDecode)
	}

	log.Println("Loaded config file : ", cfg)
	return cfg
}

// --------------------------------------------------------------
func loadProcConfig(cfg_path string) ProcConfig {
	file, errOpen := os.Open(cfg_path)
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	decoder := json.NewDecoder(file)

	cfg := ProcConfig{}
	errDecode := decoder.Decode(&cfg)
	if errDecode != nil {
		log.Fatal(errDecode)
	}

	log.Println("Loaded process config file : ", cfg)
	return cfg
}

func main() {
	// parse flags
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	// load config files
	_cfg = loadConfig(*_cfg_path)
	_cfg_proc = loadProcConfig(*_cfg_proc_path)

	// handle requests
	http.HandleFunc("/run", handleRun)
	http.HandleFunc("/stop", handleStop)
	http.HandleFunc("/restart", handleRestart)
	http.HandleFunc("/info", handleInfo)
	http.HandleFunc("/profiles", handleProfiles)

	http.ListenAndServe(":"+_cfg.HttpListenPort, nil)
}
