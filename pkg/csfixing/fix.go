package csfixing

import "fmt"

func Fix(conf ApplicationConfig) {
	fmt.Println(conf.getCommandToRun())
	// @todo do some stuff here
}