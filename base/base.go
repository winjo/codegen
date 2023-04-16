package base

import "flag"

type Command struct {
	Name string
	Run  func(cmd *Command, args []string)
	Flag flag.FlagSet
}
