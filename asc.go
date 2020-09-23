package main

import (
	"github.com/kreuzwerker/asc/command"
)

var (
	build   string
	time    string
	version string
)

func main() {
	command.Execute(version, build, time)
}
