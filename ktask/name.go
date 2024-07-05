package ktask

import (
	"errors"
	"regexp"
	"strings"
)

type Name []string

var nameLinePattern = regexp.MustCompile("^[\\p{Zs}\t]*$")

// NewEntry creates a Name from individual lines of text.
// Except for the first line, none of the lines can be empty or blank.
func NewName(line ...string) (Name, error) {
	for i, l := range line {
		if i == 0 {
			continue
		}
		if len(l) == 0 || nameLinePattern.MatchString(l) {
			return nil, errors.New("MALFORMED_SUMMARY")
		}
	}
	return line, nil
}

func (s Name) Lines() []string {
	return s
}

func (s Name) LinesWithoutFirstTag() []string {
	var ret []string
	found := false
	for _, l := range s {
		if !found {
			for _, m := range HashTagPattern.FindAllStringSubmatch(l, -1) {
				l = strings.ReplaceAll(l, m[0]+" ", "")
				l = strings.ReplaceAll(l, m[0], "")
				found = true
				break
			}
		}
		ret = append(ret, l)
	}
	return ret
}

func (s Name) Tags() *TagSet {
	tags := NewEmptyTagSet()
	for _, l := range s {
		for _, m := range HashTagPattern.FindAllStringSubmatch(l, -1) {
			tag, _ := NewTagFromString(m[0])
			tags.Put(tag)
		}
	}
	return &tags
}

func (s Name) Equals(name Name) bool {
	if len(s) != len(name) {
		return false
	}
	for i, l := range s {
		if l != name[i] {
			return false
		}
	}
	return true
}

// Append appends a text to an entry summary
func (s Name) Append(appendableText string) Name {
	if len(s) == 0 {
		return []string{appendableText}
	}
	delimiter := ""
	lastLine := s[len(s)-1]
	if len(lastLine) > 0 {
		delimiter = " "
	}
	s[len(s)-1] = lastLine + delimiter + appendableText
	return s
}
