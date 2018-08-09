package data

import (
	"github.com/BurntSushi/toml"
	"io"
	"time"
)

type Record struct {
	Start time.Time `toml:"start,omitempty"`
	Stop  time.Time `toml:"stop,omitempty"`
}

func NewRecord(start time.Time, stop time.Time) *Record {
	return &Record{start, stop}
}

func (r *Record) Duration() time.Duration {
	if r.Stop.IsZero() {
		return 0
	}

	return r.Stop.Sub(r.Start)
}

type RecordTable struct {
	Tag     string `toml:",omitempty"`
	Records []*Record
}

func Read(reader io.Reader) (*RecordTable, error) {
	var table RecordTable

	_, err := toml.DecodeReader(reader, &table)

	return &table, err
}

func (rt *RecordTable) Write(writer io.Writer) error {
	return toml.NewEncoder(writer).Encode(rt)
}

func (rt *RecordTable) Add(r *Record) {
	rt.Records = append(rt.Records, r)
}

func (rt *RecordTable) Update(stop time.Time) {
	last := len(rt.Records) - 1
	lastRecord := rt.Records[last]

	if lastRecord.Stop.IsZero() {
		lastRecord.Stop = stop
	}
}
