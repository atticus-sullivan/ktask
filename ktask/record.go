package ktask

import (
	"errors"
	"time"
)

type Record interface {
	Stage() Stage
	// Entries returns a list of all entries that are associated with this record.
	Entries() []Entry
	// SetEntries associates new entries with the record.
	SetEntries([]Entry)
	// AddEntry adds an entry to the record
	AddEntry(n Name, c, m time.Time, i int)
	addEntry(e Entry)
	SplitOnFunc(pred func(e *Entry) bool) (Record, Record)
	Merge(rs ...Record) error
}

func NewRecord(stage Stage) Record {
	return &record{
		stage: stage,
	}
}

type record struct {
	stage   Stage
	entries []Entry
}

func (r *record) Stage() Stage {
	return r.stage
}

func (r *record) Entries() []Entry {
	return r.entries
}

func (r *record) SetEntries(es []Entry) {
	r.entries = es
}

func (r *record) AddEntry(n Name, c, m time.Time, i int) {
	r.entries = append(r.entries, NewEntry(n, c, m, i))
}

func (r *record) addEntry(e Entry) {
	r.entries = append(r.entries, e)
}

func (r *record) SplitOnFunc(pred func(e *Entry) bool) (Record, Record) {
	r1 := NewRecord(r.stage)
	r2 := NewRecord(r.stage)

	for _, i := range r.entries {
		if pred(&i) {
			r1.addEntry(i)
		} else {
			r2.addEntry(i)
		}
	}

	return r1, r2
}

func (r *record) Merge(rs ...Record) error {
	for _, ro := range rs {
		if ro.Stage() != r.Stage() {
			return errors.New("mismatching record stages ocurred during merge")
		}
		r.entries = append(r.entries, ro.Entries()...)
	}
	return nil
}
