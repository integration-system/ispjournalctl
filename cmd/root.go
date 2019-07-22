package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	root = &cobra.Command{Use: "ispjournalctl"}
)

func Execute(version string) {
	root.Version = version
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
