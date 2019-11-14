package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"ispjournalctl/command"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read isp journal file on local machine",
	RunE:  command.Read.Run,
}

func init() {
	readCmd.Flags().String("file", "", "source file to read")
	readCmd.Flags().Bool("gz", false, "source file is gzipped")
	readCmd.Flags().IntP("n", "n", -1, "log entries count from start, default: read all")
	readCmd.Flags().String("since", "", "since time in format 2018-06-15 [08:15:00]")
	readCmd.Flags().String("until", "", "until time in format 2018-06-15 [08:15:00]")
	readCmd.Flags().StringSlice("event", []string{}, "filtered events, format: [--event='event1' --event='event2'], empty: show all")
	readCmd.Flags().StringSlice("level", []string{}, "filtered log levels, format: [--level='OK' --level='WARN', --level='ERROR'], empty: show all")
	readCmd.Flags().StringP("out", "o", "text", "output format in csv with ';' or json, example: --out='csv'")

	readCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if err := viper.BindPFlags(readCmd.Flags()); err != nil {
			panic(err)
		}
	}

	root.AddCommand(readCmd)
}
