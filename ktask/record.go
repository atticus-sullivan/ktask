package ktask

import "time"

type Record interface {
	Stage() Stage
	// Entries returns a list of all entries that are associated with this record.
	Entries() []Entry
	// SetEntries associates new entries with the record.
	SetEntries([]Entry)
	// AddEntry adds an entry to the record
	AddEntry(n Name, c, m time.Time, i int)
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
