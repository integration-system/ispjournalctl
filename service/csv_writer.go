package service

import (
	"encoding/csv"
	"github.com/integration-system/isp-journal/entry"
	"io"
)

var (
	headers = []string{"module_name", "host", "event", "level", "time", "request", "response", "error_text"}
)

type csvWriter struct {
	wr            io.WriteCloser
	csvWr         *csv.Writer
	headerWritten bool
}

func NewCsvWriter(wr io.WriteCloser) Writer {
	writer := csv.NewWriter(wr)
	writer.Comma = ';'
	return &csvWriter{
		wr:            wr,
		csvWr:         writer,
		headerWritten: false,
	}
}

func (w *csvWriter) Write(entry *entry.Entry) error {
	if !w.headerWritten {
		if err := w.csvWr.Write(headers); err != nil {
			return err
		}
		w.headerWritten = true
	}

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
	return w.csvWr.Write(row)
}

func (w *csvWriter) Close() error {
	w.csvWr.Flush()

	return w.wr.Close()
}
