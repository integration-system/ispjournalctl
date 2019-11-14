package service

import (
	"github.com/integration-system/isp-journal/entry"
	"github.com/integration-system/isp-journal/search"
	"io"
	"strings"
)

type textWriter struct {
	wr io.WriteCloser
}

func NewTextWriter(wr io.WriteCloser) Writer {
	return &textWriter{
		wr: wr,
	}
}

func (w *textWriter) WriteRead(entry *entry.Entry) error {
	row := []string{
		entry.ModuleName,
		entry.Host,
		entry.Event,
		entry.Level,
		entry.Time,
		string(entry.Request),
		string(entry.Response),
		entry.ErrorText,
	}
	s := strings.Join(row, ";") + "\n"
	_, err := w.wr.Write([]byte(s))
	return err
}

func (w *textWriter) WriteSearch(entry *search.SearchResponse) error {
	row := []string{
		entry.ModuleName,
		entry.Host,
		entry.Event,
		entry.Level,
		entry.Time,
		entry.Request,
		entry.Response,
		entry.ErrorText,
	}
	s := strings.Join(row, ";") + "\n"
	_, err := w.wr.Write([]byte(s))
	return err
}

func (w *textWriter) Close() error {
	return w.wr.Close()
}
