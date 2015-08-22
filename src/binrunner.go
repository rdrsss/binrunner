/*
 * @file : processmanager.go
 * @brief : binrunner, start, kill, list, processes via http
 *			interface.
 * @author : Manuel A. Rodriguez (manuel.rdrs@gmail.com)
 */
package main

import (
	"encoding/json"
	"flag"
	"log"
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

// --------------------------------------------------------------
func loadConfigFile(cfg_path string) Config {
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
	flag.Parse()
	log.SetFlags(log.Lshortfile)

}
