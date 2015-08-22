/*
 * @file : process.go
 * @brief : start, stop, etc... processes from proc.cfg.
 * @author : Manuel A. Rodriguez (manuel.rdrs@gmail.com)
 */
package main

import (
	"errors"
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
