package main

import (
	"github.com/zjy-dev/gcovr-json-util/cmd"
)

func main() {
	// Pass version info to cmd package
	cmd.SetVersionInfo(Version, GitCommit, BuildDate)
	cmd.Execute()
}
