package parser

import (
	"ktask/ktask"
	"strings"
	"time"
)

// Serialiser is used when the output should be modified, e.g. coloured.
type Serialiser interface {
	Date(time.Time) string
	Stage(ktask.Stage) string
	Name(NameText) string
}

type Line struct {
	Text   string
	Record ktask.Record
	EntryI int
}

type Lines []Line

var canonicalLineEnding = "\n"
var canonicalIndentation = "    "

func (ls Lines) ToString() string {
	builder := strings.Builder{}
	for _, l := range ls {
		builder.WriteString(l.Text)
		builder.WriteString(canonicalLineEnding)
	}
	return builder.String()
}

// SerialiseRecords serialises records into the canonical string representation.
// (So it doesn’t and cannot restore the original formatting!)
func SerialiseRecords(s Serialiser, rs ...ktask.Record) Lines {
	var lines []Line
	for i, r := range rs {
		lines = append(lines, serialiseRecord(s, r)...)
		if i < len(rs)-1 {
			lines = append(lines, Line{"", nil, -1})
		}
	}
	return lines
}

func serialiseRecord(s Serialiser, r ktask.Record) []Line {
	var lines []Line
	headline := s.Stage(r.Stage())
	lines = append(lines, Line{headline, r, -1})
	for entryI, e := range r.Entries() {
		cValue := s.Date(e.CreatedAt())
		mValue := s.Date(e.ModifiedAt())
		lines = append(lines, Line{canonicalIndentation + cValue + " " + mValue, r, entryI})
		for i, l := range e.Name().Lines() {
			summaryText := s.Name([]string{l})
			if i == 0 && l != "" {
				lines[len(lines)-1].Text += " " + summaryText
			} else if i >= 1 {
				lines = append(lines, Line{canonicalIndentation + canonicalIndentation + summaryText, r, entryI})
			}
		}
	}
	return lines
}

type NameText []string

func (s NameText) ToString() string {
	return strings.Join(s, canonicalLineEnding)
}
