package main

import (
	"fmt"
	"os"

	"github.com/winjo/codegen/base"
	"github.com/winjo/codegen/dal"
	"github.com/winjo/codegen/version"
)

var commands = []*base.Command{
	version.CmdVersion,
	dal.CmdDal,
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		mainUsage()
	}

	cmdName := args[0]
	if cmdName == "-h" || cmdName == "help" {
		mainUsage()
	}

	for _, cmd := range commands {
		if cmd.Name != cmdName {
			continue
		}
		err := cmd.Flag.Parse(args[1:])
		if err == nil {
			args = cmd.Flag.Args()
			cmd.Run(cmd, args)
		}
		os.Exit(0)
		return
	}

	fmt.Fprintf(os.Stderr, "codegen %s: unknown command\nRun 'codegen -h' for usage.\n", cmdName)
	os.Exit(2)
}

func mainUsage() {
	fmt.Fprintln(os.Stderr, `codegen is a tool for generating go code

Usage:
	codegen dal    generate dal dao code
			`)
	os.Exit(2)
}
