package service

import (
	"fmt"
	"github.com/integration-system/isp-journal/entry"
	"io"
)

type Writer interface {
	Write(entry *entry.Entry) error
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
