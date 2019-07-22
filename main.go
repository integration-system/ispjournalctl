package main

import "ispjournalctl/cmd"

var (
	version = "1.0.0"
)

func main() {
	cmd.Execute(version)
}
