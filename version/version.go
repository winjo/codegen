package version

import (
	"fmt"

	"github.com/winjo/codegen/base"
)

const version = "v0.0.1"

var CmdVersion = &base.Command{
	Name: "version",
}

func init() {
	CmdVersion.Run = run
}

func run(cmd *base.Command, args []string) {
	fmt.Printf("version: %s\n", version)
}
