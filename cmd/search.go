package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"ispjournalctl/command"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search isp journal file on jsp-journal-service",
	RunE:  command.Search.Run,
}

func init() {
	searchCmd.Flags().StringP("gate", "g", "", "gate to isp-journal-service in format '127.0.0.0:0000'")
	searchCmd.Flags().String("module", "", "module name")
	searchCmd.Flags().IntP("n", "n", -1, "log entries count from start, default: read all")
	searchCmd.Flags().String("since", "", "since time in format 2018-06-15 [08:15:00]")
	searchCmd.Flags().String("until", "", "until time in format 2018-06-15 [08:15:00]")
	searchCmd.Flags().StringSlice("host", []string{}, "filtered host, format: [--host='host1' --host='host2'], empty: show all")
	searchCmd.Flags().StringSlice("event", []string{}, "filtered events, format: [--event='event1' --event='event2'], empty: show all")
	searchCmd.Flags().StringSlice("level", []string{}, "filtered log levels, format: [--level='OK' --level='WARN', --level='ERROR'], empty: show all")
	searchCmd.Flags().StringP("out", "o", "text", "output format in csv with ';' or json, example: --out='csv'")

	searchCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if err := viper.BindPFlags(searchCmd.Flags()); err != nil {
			panic(err)
		}
	}
	root.AddCommand(searchCmd)
}
