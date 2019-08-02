package service

import (
	"fmt"
	"github.com/integration-system/isp-journal/entry"
	"github.com/integration-system/isp-journal/search"
	"io"
)

type Writer interface {
	WriteRead(entry *entry.Entry) error
	WriteSearch(entry *search.SearchResponse) error
	io.Closer
}

func NewWriter(t string, wr io.WriteCloser) (Writer, error) {
	switch t {
	case "csv":
		return NewCsvWriter(wr), nil
	case "json":
		return NewJsonWriter(wr), nil
	default:
		return nil, fmt.Errorf("unknown output format: %v", t)
	}
}
