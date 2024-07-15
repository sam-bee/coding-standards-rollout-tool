package csfixing

import "fmt"

func Fix(conf ApplicationConfig, git gitInterface) {
	fmt.Println(conf.getCommandToRun())
	// @todo do some stuff here
}
