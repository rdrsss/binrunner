/*
 * @file : client.go
 * @brief : Client to test the web frontend of the binrunner service.
 * @author : Manuel A. Rodriguez (manuel.rdrs@gmail.com)
 */
package main

import (
	"encoding/json"
	"flag"
	"github.com/flynn/go-docopt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	HttpAddr string
}

var (
	_cfg_dir     = flag.String("config_dir", "../cfg", "path to config directory")
	_cfg_path    = flag.String("config_path", *_cfg_dir+"/client.cfg", "path to config file")
	_server_addr = flag.String("server_addr", "", "server to request from")
	_cfg         Config
)

// --------------------------------------------------------------
func loadConfigFile() Config {
	file, errOpen := os.Open(*_cfg_path)
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

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	// Read in config file
	_cfg = loadConfigFile()
	if len(*_server_addr) <= 0 {
		*_server_addr = _cfg.HttpAddr
	}

	usage := `client.
	Usage:
		client <command> [<args>...]

	Options:
		-h, --help

	Commands:
		run 	attempts to run a command passed into the client.
		stop	stops a command via passed in pid.
		list	requests a list of running processes.`

	args, _ := docopt.Parse(usage, nil, true, "client", true)

	cmd := args.String["<command>"]
	cmdArgs := args.All["<args>"].([]string)

	log.Println("command passed in : ", cmd)
	log.Println("args passed in : ", cmdArgs)

	switch cmd {
	case "run":
		if len(cmdArgs) == 0 { // help
			log.Println(usage)
		} else {
			var run, count string
			for _, val := range cmdArgs {
				println(val)
				if strings.HasPrefix(val, "--alias=") {
					sp := strings.Split(val, "--alias=")
					if len(sp) > 0 {
						run = sp[1]
					}
				}

				if strings.HasPrefix(val, "--count=") {
					sp := strings.Split(val, "--count=")
					if len(sp) > 0 {
						count = sp[1]
					}
				}
			}
			run = strings.Trim(run, " ")
			log.Println("runnable command : ", run)
			v := url.Values{}
			v.Set("alias", run)
			v.Set("proc_count", count)
			resp, err := http.PostForm(*_server_addr+"/run", v)
			if err == nil {
				log.Println(resp)
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				log.Println(string(body))
			} else {
				log.Println(err)
			}
		}
	case "stop":
		if len(cmdArgs) == 0 { // help
			log.Println(usage)
		} else if cmdArgs[0] == "--pid" {
		}
	case "list":
	default:
		log.Println(cmd, "Is an unknown command, see usage")
	}

}
