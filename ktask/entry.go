package ktask

import (
	"strings"
	"time"
)

type Entry struct {
	name       Name
	createdAt  time.Time
	modifiedAt time.Time
	index      int
}

func NewEntry(name Name, createdAt, modifiedAt time.Time, index int) Entry {
	return Entry{
		name:       name,
		createdAt:  createdAt,
		modifiedAt: modifiedAt,
		index:      index,
	}
}

func (e *Entry) Index() int {
	return e.index
}

func (e *Entry) Name() Name {
	return e.name
}

func (e *Entry) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Entry) ModifiedAt() time.Time {
	return e.modifiedAt
}

func (e Entry) FilterValue() string {
	builder := strings.Builder{}

	first := true
	for _, n := range e.name {
		if !first {
			builder.WriteRune(' ')
		} else {
			first = false
		}
		builder.WriteString(n)
	}

	for k := range e.name.Tags().ForLookup() {
		builder.WriteString(k.ToString())
	}

	return builder.String()
}

// define how this should be rendered with the default delegate
func (e Entry) Title() string {
	builder := strings.Builder{}

	first := true
	for _, n := range e.name.LinesWithoutFirstTag() {
		if !first {
			builder.WriteRune(' ')
		} else {
			first = false
		}
		builder.WriteString(n)
	}
	return builder.String()
}

// define how this should be rendered with the default delegate
func (e Entry) Description() string {
	if e.name.Tags().IsEmpty() {
		return ""
	}
	return e.name.Tags().original[0].ToString()
}
