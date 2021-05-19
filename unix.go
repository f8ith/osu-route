// +build linux darwin
package main

import (
	"os/exec"
	"log"
)

func run(name string, args ...string) {
	err := exec.Command(name, args...).Run()
	if err != nil {
		log.Fatal(err)
	}
}
