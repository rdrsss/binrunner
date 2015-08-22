/*
 * @file : proc_test.go
 * @brief : bit runner unit tests
 * @author : Manuel A. Rodriguez (manuel.rdrs@gmail.com)
 */
package main

import (
	"testing"
)

// --------------------------------------------------------------
func TestAStart(t *testing.T) {
	_, err := StartProcess("top", "")
	if err != nil {
		t.Error(err)
	}
}

// --------------------------------------------------------------
func TestBListProcesses(t *testing.T) {
	pids := GetAllProcesses()
	for _, v := range pids {
		t.Log("pid : ", v)
	}
	if len(pids) < 1 {
		t.Error("no running processes")
	}
}

// --------------------------------------------------------------
func TestCKillProcesses(t *testing.T) {
	pids := GetAllProcesses()
	if len(pids) < 1 {
		t.Error("no running processes")
	} else {
		for _, v := range pids {
			err := KillProcess(v)
			if err != nil {
				t.Log(err)
			}
		}
	}
}
