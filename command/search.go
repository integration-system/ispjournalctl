package command

import (
	"fmt"
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
	Gate   string
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

	if cfg.Gate == "" {
		return errors.New("Gate name is required")
	} else {
		service.JournalServiceClient.ReceiveConfiguration(cfg.Gate)
	}

	writer, err := service.NewWriter(cfg.Out, os.Stdout)
	defer func() { _ = writer.Close() }()
	if err != nil {
		return err
	}

	if err := cfg.searchLogs(writer); err != nil {
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

func (cfg searchCmdCfg) searchLogs(writer service.Writer) error {
	req := cfg.getSearchRequest()
	if req.From.IsZero() { //TODO hack
		req.From = req.From.Add(1 * time.Second)
	}

	currentRows := 0
	for {
		if response, err := service.JournalServiceClient.Search(req); err != nil {
			return err
		} else {
			if len(response) == 0 {
				break
			}
			for _, value := range response {
				if currentRows == cfg.N {
					break
				} else if err := writer.WriteSearch(&value); err != nil {
					return err
				}
				currentRows++
			}
			req.Offset = req.Offset + req.Limit
		}
	}
	if currentRows == 0 {
		fmt.Println("entries not found")
	}
	return nil
}
