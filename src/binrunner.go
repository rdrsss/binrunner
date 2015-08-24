/*
 * @file : binrunner.go
 * @brief : binrunner, start, kill, list, processes via http
 *			interface.
 * @author : Manuel A. Rodriguez (manuel.rdrs@gmail.com)
 */
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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
	Alias           string
	Cmd             string
	Args            string
	RestartAttempts int
}

type ProfileCache struct {
	Entries    map[string]EntryConfig
	pc_rw_lock sync.RWMutex
}

// response structs
type CmdResponse struct {
	Pid []int
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
	// cmd
	_pcache = ProfileCache{Entries: make(map[string]EntryConfig)}
)

// **************************************************************
func (pc ProfileCache) HasProfile(alias string) bool {
	var ret = false
	pc.pc_rw_lock.RLock()
	if _, ok := pc.Entries[alias]; ok == true {
		ret = ok
	}
	pc.pc_rw_lock.RUnlock()
	return ret
}

// **************************************************************
func (pc ProfileCache) GetProfile(alias string) (EntryConfig, error) {
	var entry EntryConfig
	var err error = nil
	pc.pc_rw_lock.RLock()
	if _, ok := pc.Entries[alias]; ok == true {
		entry = pc.Entries[alias]
	} else {
		err = errors.New("No such profile in cache")
	}
	pc.pc_rw_lock.RUnlock()
	return entry, err
}

// **************************************************************
func (pc ProfileCache) AddProfile(entry EntryConfig) {
	pc.pc_rw_lock.Lock()
	pc.Entries[entry.Alias] = entry
	pc.pc_rw_lock.Unlock()
}

// **************************************************************
func (pc ProfileCache) RemoveProfile(alias string) {
	pc.pc_rw_lock.Lock()
	if _, ok := pc.Entries[alias]; ok == true {
		delete(pc.Entries, alias)
	}
	pc.pc_rw_lock.Unlock()
}

// **************************************************************
func (pc ProfileCache) ClearCache() {
	pc.pc_rw_lock.Lock()
	// just overwrite the map and let the gc take care of the rest
	pc.Entries = make(map[string]EntryConfig)
	pc.pc_rw_lock.Unlock()
}

// --------------------------------------------------------------
func handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		alias := r.FormValue("alias")
		count, _ := strconv.Atoi(r.FormValue("proc_count"))
		if count <= 0 {
			count = 1
		}
		//_ := r.FormValue("args")
		if len(alias) > 0 {
			// find lookup alias
			p, err := _pcache.GetProfile(alias)
			if err != nil {
				// 400
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("No such alias in cache"))
			} else {
				var pids []int
				var success bool = true
				for i := 0; i < count; i++ {
					pcmd, err := StartProcess(p.Cmd, p.Args)
					if err != nil {
						// 400
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(err.Error()))
						success = false
						break

					}
					pids = append(pids, pcmd.Process.Pid)
				}
				if success {
					var r = CmdResponse{pids}
					md, _ := json.Marshal(r)
					// 200
					w.WriteHeader(http.StatusOK)
					w.Write(md)
				}
			}
		} else {
			// 400
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No cmd passed"))
		}

	} else {
		// 400
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Method"))
	}
}

// --------------------------------------------------------------
func handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		pid := r.FormValue("pid")
		if len(pid) > 0 {
			if strings.ToLower(pid) == "all" {

			} else {
				ipid, err := strconv.Atoi(pid)
				if err != nil {
					// 400
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
				} else {
					err = KillProcess(ipid)
					if err != nil {
						// 400
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(err.Error()))
					} else {
						// 200
						w.WriteHeader(http.StatusOK)
					}
				}
			}
		} else {
			// 400
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No pid passed"))
		}

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
		pid := r.FormValue("pid")
		if len(pid) > 0 {
			// TODO : list info specific to the pid
		} else {
			// return list
			pids := _pmap.ListPids()
			pcount := strconv.Itoa(len(pids))

			var jspids string = "{\n\"count\":\"" + pcount + "\""
			if len(pids) > 0 {
				jspids += ",\n"
				for k, v := range pids {
					pos := strconv.Itoa(k)
					pid := strconv.Itoa(v)
					jspids += "\"" + pos + "\":\"" + pid + "\""
					if k < len(pids) {
						jspids += ",\n"
					}
				}
				jspids += "\n}"
			}
			// 200
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(jspids))
		}

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
func reloadProfileCache() {
	_pcache.ClearCache()
	_pcache.pc_rw_lock.Lock()
	for _, v := range _cfg_proc.Entries {
		log.Println("adding entry : ", v.Alias, " [", v, "]")
		_pcache.Entries[v.Alias] = v
	}
	_pcache.pc_rw_lock.Unlock()

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

	//
	reloadProfileCache()

	// handle requests
	http.HandleFunc("/run", handleRun)
	http.HandleFunc("/stop", handleStop)
	http.HandleFunc("/restart", handleRestart)
	http.HandleFunc("/info", handleInfo)
	http.HandleFunc("/profiles", handleProfiles)

	http.ListenAndServe(":"+_cfg.HttpListenPort, nil)
}
