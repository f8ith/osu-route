package main

import (
	"os/exec"
	"syscall"
)

func run(name string, args ...string) {
	args = append([]string{"/c", name}, args...)
	cmdInstance := exec.Command("C:\\Windows\\system32\\cmd.exe", args...)
	cmdInstance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmdInstance.Run()
	if err != nil {
		log.Fatal(err)
	}
}
