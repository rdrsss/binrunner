/*
 * @file : process.go
 * @brief : start, stop, etc... processes from proc.cfg.
 * @author : Manuel A. Rodriguez (manuel.rdrs@gmail.com)
 */
package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"sync"
)

type ProcessMap struct {
	proc_map   map[int]*exec.Cmd // pid - cmd
	pm_rw_lock sync.RWMutex
}

/*
 * Global vars
 */
var (
	_pmap = ProcessMap{proc_map: make(map[int]*exec.Cmd)}
)

// --------------------------------------------------------------
func (pm ProcessMap) Append(cmd *exec.Cmd) {
	pm.pm_rw_lock.Lock()
	pm.proc_map[cmd.Process.Pid] = cmd
	pm.pm_rw_lock.Unlock()
}

// --------------------------------------------------------------
func (pm ProcessMap) Remove(pid int) bool {
	var ret bool = false
	pm.pm_rw_lock.Lock()
	delete(pm.proc_map, pid)
	pm.pm_rw_lock.Unlock()
	return ret
}

// --------------------------------------------------------------
func (pm ProcessMap) HasProcess(pid int) bool {
	var ret bool = false
	pm.pm_rw_lock.RLock()
	if _, ok := pm.proc_map[pid]; ok == true {
		ret = true
	}
	pm.pm_rw_lock.RUnlock()
	return ret
}

// --------------------------------------------------------------
func (pm ProcessMap) GetProcess(pid int) (*os.Process, error) {
	var err error = nil
	var proc *os.Process = nil

	pm.pm_rw_lock.RLock()
	if _, ok := pm.proc_map[pid]; ok == true {
		proc = pm.proc_map[pid].Process
	} else {
		err = errors.New("Not valid process id within process map")
	}
	pm.pm_rw_lock.RUnlock()

	return proc, err
}

// --------------------------------------------------------------
func (pm ProcessMap) GetProcesses(pids []int) ([]*os.Process, error) {
	var err error = nil
	var procs [](*os.Process)

	pm.pm_rw_lock.RLock()
	for _, pid := range pids {
		if _, ok := pm.proc_map[pid]; ok == true {
			procs = append(procs, pm.proc_map[pid].Process)
		} else {
			err = errors.New("Not valid process id within process map")
		}
	}
	pm.pm_rw_lock.RUnlock()

	return procs, err
}

// --------------------------------------------------------------
func (pm ProcessMap) ListPids() []int {
	pm.pm_rw_lock.RLock()
	if len(pm.proc_map)-1 < 0 {
		pm.pm_rw_lock.RUnlock()
		return nil
	}

	var pids = make([]int, len(pm.proc_map)-1)
	for key, _ := range pm.proc_map {
		pids = append(pids, key)
	}
	pm.pm_rw_lock.RUnlock()
	return pids
}

// --------------------------------------------------------------
func (pm ProcessMap) Count() int {
	var count int = -1
	return count
}

/* public functions */

// --------------------------------------------------------------
func StartProcess(cmd_str, args string) (*exec.Cmd, error) {
	var err error = nil
	var cmd *exec.Cmd = nil
	cmd = exec.Command(cmd_str, args)
	if cmd != nil {
		log.Println("Starting process [", cmd_str, "] args[", args, "]")
		err = cmd.Start()
		if err == nil {
			// no error
			_pmap.Append(cmd)
		} else {
			// Log stderr, and stdout
			log.Println(err)
			log.Println("stdout : [%s]", cmd.Stdout)
			log.Println("stderr : [%s]", cmd.Stderr)
		}
	} else {
		err = errors.New("Failed to create an executable command")
	}

	return cmd, err
}

// --------------------------------------------------------------
func KillProcess(pid int) error {
	var err error = nil
	if _pmap.HasProcess(pid) == true {
		// kill the process
		proc, _ := _pmap.GetProcess(pid)
		err = proc.Kill()
		// log pid, and track app info
		go func() {
			pstate, perr := proc.Wait()
			if perr != nil {
				log.Println(perr)
			}
			if pstate != nil {
				log.Println(pstate.String())
			}

			_pmap.Remove(proc.Pid)
		}()

	} else {
		err = errors.New("Could not find pid in process map")
	}
	return err
}

// --------------------------------------------------------------
func KillAllProcesses() {
	pids := _pmap.ListPids()
	procs, _ := _pmap.GetProcesses(pids)
	for _, proc := range procs {
		// kill the process
		err := proc.Kill()
		if err == nil {
			// log pid, and track app info
			go func() {
				pstate, perr := proc.Wait()
				if perr != nil {
					log.Println(perr)
				}
				if pstate != nil {
					log.Println(pstate.String())
				}

				_pmap.Remove(proc.Pid)
			}()
		}
	}
}

// --------------------------------------------------------------
func ProcessInfo(pid int) error {
	var err error = nil
	// TODO:: display info for pid
	return err
}

// --------------------------------------------------------------
func GetAllProcesses() []int {
	return _pmap.ListPids()
}
