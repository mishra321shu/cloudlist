package main

import (
	"github.com/mishra321shu/cloudlist/internal/runner"
	"github.com/projectdiscovery/gologger"
)

func main() {
	options := runner.ParseOptions()
	runner, err := runner.New(options)
	if err != nil {
		gologger.Fatalf("Could not create runner: %s\n", err)
	}
	runner.Enumerate()
}
