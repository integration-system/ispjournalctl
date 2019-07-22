package service

import (
	"github.com/integration-system/isp-journal/entry"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/modern-go/reflect2"
	"io"
	"time"
	"unicode"
	"unsafe"
)

const (
	FullDateFormat = "2006-01-02T15:04:05.999-07:00"
)

var (
	ji = jsoniter.ConfigFastest
)

func init() {
	extra.SetNamingStrategy(toCamelCase)

	tc := &timeCoder{}
	timeType := reflect2.TypeByName("time.Time")

	encExt := jsoniter.EncoderExtension{timeType: tc}
	decExt := jsoniter.DecoderExtension{timeType: tc}
	ji.RegisterExtension(encExt)
	ji.RegisterExtension(decExt)
}

type row struct {
	*entry.Entry
	Request  string
	Response string
}

type jsonWriter struct {
	wr io.WriteCloser
}

func (w *jsonWriter) Write(entry *entry.Entry) error {
	r := row{Entry: entry, Request: string(entry.Request), Response: string(entry.Response)}
	bytes, err := ji.Marshal(r)
	if err != nil {
		return err
	}

	if _, err := w.wr.Write(append(bytes, '\n')); err != nil {
		return err
	} else {
		return nil
	}
}

func (w *jsonWriter) Close() error {
	return w.wr.Close()
}

func NewJsonWriter(wr io.WriteCloser) Writer {
	return &jsonWriter{
		wr: wr,
	}
}

type timeCoder struct {
}

func (codec *timeCoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	t, err := time.Parse(FullDateFormat, iter.ReadString())
	if err != nil {
		iter.ReportError("string -> time.Time", err.Error())
	} else {
		*((*time.Time)(ptr)) = t
	}
}

func (codec *timeCoder) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Time)(ptr))
	return ts.IsZero()
}

func (codec *timeCoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Time)(ptr))
	stream.WriteString(ts.Format(FullDateFormat))
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	arr := []rune(s)
	arr[0] = unicode.ToLower(arr[0])
	return string(arr)
}
