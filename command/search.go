package command

import (
	domain "github.com/integration-system/isp-journal/search"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"ispjournalctl/service"
	"ispjournalctl/util"
	"os"
	"time"
)

const (
	limit  = 1000
	offset = 0
)

var Search searchCmdCfg

type searchCmdCfg struct {
	Module string
	Since  time.Time
	Until  time.Time
	Host   []string
	Event  []string
	Level  []string
	Out    string
	N      int
}

func (cfg searchCmdCfg) Run(cmd *cobra.Command, args []string) error {
	if err := util.UnmarshalConfig(&cfg); err != nil {
		return err
	}

	if cfg.Module == "" {
		return errors.New("Module name is required")
	}

	writer, err := service.NewWriter(cfg.Out, os.Stdout)
	defer func() { _ = writer.Close() }()
	if err != nil {
		return err
	}

	request := cfg.getSearchRequest()
	if err := cfg.searchLogs(request, writer, 0, cfg.N); err != nil {
		return err
	}
	return nil
}

func (cfg searchCmdCfg) getSearchRequest() domain.SearchRequest {
	searchRequest := domain.SearchRequest{
		ModuleName: cfg.Module,
		Host:       cfg.Host,
		Event:      cfg.Event,
		Level:      cfg.Level,
		Limit:      limit,
		Offset:     offset,
	}
	if !cfg.Since.IsZero() {
		searchRequest.From = cfg.Since
	}
	if !cfg.Until.IsZero() {
		searchRequest.To = cfg.Until
	}
	return searchRequest
}

func (cfg searchCmdCfg) searchLogs(req domain.SearchRequest, writer service.Writer, currentRows, limitRows int) error {
	if response, err := service.JournalServiceClient.Search(req); err != nil {
		return err
	} else {
		if len(response) == 0 {
			return nil
		}
		for _, value := range response {
			if currentRows == limitRows {
				return nil
			} else {
				currentRows++
			}
			if err := writer.WriteSearch(&value); err != nil {
				return err
			}
		}
		req.Offset = req.Offset + req.Limit
		if err := cfg.searchLogs(req, writer, currentRows, limitRows); err != nil {
			return err
		} else {
			return nil
		}
	}
}
