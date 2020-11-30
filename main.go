package main

import (
	"github.com/inanme/vergo/cmd"
)

var (
	version, commit, date, builtBy string
)

func init() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
	cmd.BuiltBy = builtBy
}

func main() {
	cmd.Execute()
}
