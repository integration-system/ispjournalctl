package command

import (
	domain "github.com/integration-system/isp-journal/search"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"ispjournalctl/service"
	"ispjournalctl/util"
	"os"
	"time"
)

var Read readCmdCfg

type readCmdCfg struct {
	File  string
	Since time.Time
	Until time.Time
	Event []string
	Level []string
	N     int
	Gz    bool
	Out   string
}

func (cfg readCmdCfg) Run(cmd *cobra.Command, args []string) error {
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

	req := domain.SearchRequest{
		From:  cfg.Since,
		To:    cfg.Until,
		Event: cfg.Event,
		Level: cfg.Level,
	}
	filter, err := domain.NewFilter(req)
	if err != nil {
		return err
	}

	reader, err := domain.NewLogReader(cmd.InOrStdin(), cfg.Gz, filter)
	if reader != nil {
		defer func() { _ = reader.Close() }()
	}
	if err != nil {
		return err
	}

	writer, err := service.NewWriter(cfg.Out, os.Stdout)
	defer func() { _ = writer.Close() }()
	if err != nil {
		return err
	}

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
			if err := writer.WriteRead(entry); err != nil {
				return errors.WithMessage(err, "could not write row")
			}
			written++
		}
	}
}
