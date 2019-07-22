package cmd

import (
	"github.com/integration-system/isp-journal/search"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"ispjournalctl/service"
	"ispjournalctl/util"
	"os"
	"time"
)

var (
	localCmd = &cobra.Command{
		Use:   "local",
		Short: "Read isp journal file on local machine",
		RunE:  localRun,
	}
)

type localCmdCfg struct {
	File  string
	Since time.Time
	Until time.Time
	Event []string
	Level []string
	N     int
	Gz    bool
	Out   string
}

func init() {
	localCmd.Flags().String("file", "", "target file to read")
	localCmd.Flags().Bool("gz", false, "source file is gzipped")
	localCmd.Flags().IntP("n", "n", 10, "log entries count from start, -1 - read all")
	localCmd.Flags().String("since", "", "since time in format 2018-06-15 [08:15:00]")
	localCmd.Flags().String("until", "", "until time in format 2018-06-15 [08:15:00]")
	localCmd.Flags().StringSlice("event", []string{}, "filtered events, format: [--event='file1' --event='file2']")
	localCmd.Flags().StringSlice("level", []string{}, "filtered log levels, format: [--level='OK' --level='WARN', --level='ERROR']")
	localCmd.Flags().StringP("out", "o", "csv", "output format in csv with ';' or json, example: --out='csv'")

	if err := viper.BindPFlags(localCmd.Flags()); err != nil {
		panic(err)
	}

	root.AddCommand(localCmd)
}

func localRun(cmd *cobra.Command, args []string) error {
	cfg := localCmdCfg{}
	if err := util.UnmarshalConfig(&cfg); err != nil {
		return err
	}

	if cfg.File != "" {
		if file, err := os.Open(cfg.File); err != nil {
			return err
		} else {
			cmd.SetIn(file)
		}
	}

	req := search.SearchRequest{
		From:  cfg.Since,
		To:    cfg.Until,
		Event: cfg.Event,
		Level: cfg.Level,
	}
	if req.From.IsZero() { //TODO hack
		req.From = req.From.Add(1 * time.Second)
	}
	filter, err := search.NewFilter(req)
	if err != nil {
		return err
	}
	reader, err := search.NewLogReader(cmd.InOrStdin(), cfg.Gz, filter)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := service.NewWriter(cfg.Out, os.Stdout)
	if err != nil {
		return err
	}
	defer writer.Close()

	written := 0
	for {
		if cfg.N > 0 && written >= cfg.N {
			return nil
		}

		entry, err := reader.FilterNext()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return errors.WithMessage(err, "invalid input format")
		}

		if entry != nil {
			if err := writer.Write(entry); err != nil {
				return errors.WithMessage(err, "could not write row")
			}
			written++
		}
	}
}
