/*
Package parser contains the logic how to convert Record objects from and to plain text.
*/
package parser

import (
	"ktask/ktask"
	"time"

	"github.com/jotaen/klog/klog/parser/txt"
)

func parse(block txt.Block) (ktask.Record, []txt.Error) {
	lines, initialLineOffset, _ := block.SignificantLines()
	initialLineCount := len(lines) // Capture current value
	nr := func(lines []txt.Line) int {
		return initialLineOffset + initialLineCount - len(lines)
	}
	var errs []txt.Error

	// ========== HEADLINE ==========
	record := func() ktask.Record {
		headline := txt.NewParseable(lines[0], 0)

		// There is no leading whitespace allowed in the headline.
		if txt.IsSpaceOrTab(headline.Peek()) {
			errs = append(errs, ErrorIllegalIndentation().New(block, nr(lines), 0, headline.Length()))
			return nil
		}

		// Parse the stage
		stageText, _ := headline.PeekUntil(func(_ rune) bool { return false }) // Move forward until end of line
		rStage := ktask.Stage(stageText.ToString())
		sErr := rStage.Valid()
		if sErr != nil {
			errs = append(errs, ErrorInvalidStage().New(block, nr(lines), headline.PointerPosition, stageText.Length()))
			return nil
		}
		headline.Advance(stageText.Length())
		headline.SkipWhile(txt.IsSpaceOrTab)
		r := ktask.NewRecord(rStage)

		// Make sure there is no other text left in the headline.
		headline.SkipWhile(txt.IsSpaceOrTab)
		if headline.RemainingLength() > 0 {
			errs = append(errs, ErrorUnrecognisedTextInHeadline().New(block, nr(lines), headline.PointerPosition, headline.RemainingLength()))
		}
		return r
	}()
	lines = lines[1:]

	if record == nil {
		// In case there was an error, generate dummy record to ensure that we have something to
		// work with during parsing. That allows us to continue even if there are errors early on.
		dummyStage := ktask.Stage("ToDo")
		record = ktask.NewRecord(dummyStage)
	}

	var indentator *txt.Indentator

	// ========== ENTRIES ==========
	index := int(-1)
	for len(lines) > 0 {
		l := lines[0]
		indentator = txt.NewIndentator(txt.Indentations, l)
		if indentator == nil {
			// We should never make it here if the indentation could not be determined.
			panic("Could not detect indentation")
		}

		// Check for correct indentation.
		entry := indentator.NewIndentedParseable(l, 1)
		if entry == nil || txt.IsSpaceOrTab(entry.Peek()) {
			errs = append(errs, ErrorIllegalIndentation().New(block, nr(lines), 0, len(l.Text)))
			break
		}

		// Parse entry value.
		createEntry, evErr := func() (func(ktask.Name, int) txt.Error, txt.Error) {
			// Try to interpret the entry value as stage.
			createdAtCandidate, _ := entry.PeekUntil(txt.IsSpaceOrTab)
			createdAt, dErr := time.Parse("2006-01-02", createdAtCandidate.ToString())
			if dErr != nil {
				return nil, ErrorMalformedEntry().New(block, nr(lines), entry.PointerPosition, createdAtCandidate.Length())
			}
			entry.Advance(createdAtCandidate.Length())
			entry.SkipWhile(txt.IsSpaceOrTab)

			modifiedAtCandidate, _ := entry.PeekUntil(txt.IsSpaceOrTab)
			if modifiedAtCandidate.Length() == 0 {
				return nil, ErrorMalformedEntry().New(block, nr(lines), entry.PointerPosition, 1)
			}
			modifiedAt, dErr := time.Parse("2006-01-02", modifiedAtCandidate.ToString())
			if dErr != nil {
				return nil, ErrorMalformedEntry().New(block, nr(lines), entry.PointerPosition, createdAtCandidate.Length())
			}
			entry.Advance(modifiedAtCandidate.Length())

			return func(n ktask.Name, i int) txt.Error {
				record.AddEntry(n, createdAt, modifiedAt, i)
				return nil
			}, nil
		}()
		lines = lines[1:]
		index++

		// Check for error while parsing the entry value.
		if evErr != nil {
			errs = append(errs, evErr)
			continue
		}

		// Parse entry summary.
		entryName, esErr := func() (ktask.Name, txt.Error) {
			var result ktask.Name

			// Parse first line of entry summary.
			if txt.IsSpaceOrTab(entry.Peek()) {
				entry.Advance(1)
				nameText := entry.Remainder()
				firstLine, sErr := ktask.NewName(nameText.ToString())
				if sErr != nil {
					return nil, ErrorMalformedSummary().New(block, nr(lines), 0, nameText.Length())
				}
				result = firstLine
			} else {
				result, _ = ktask.NewName("")
			}

			// Parse subsequent lines of multiline name.
			for len(lines) > 0 {
				nextNameLine := indentator.NewIndentedParseable(lines[0], 2)
				if nextNameLine == nil {
					break
				}
				lines = lines[1:]
				additionalText, _ := nextNameLine.PeekUntil(func(_ rune) bool {
					return false // Move forward until end of line
				})
				newEntrySummary, sErr := ktask.NewName(append(result, additionalText.ToString())...)
				if sErr != nil {
					return nil, ErrorMalformedSummary().New(block, nr(lines), 0, nextNameLine.Length())
				}
				result = newEntrySummary
			}

			return result, nil
		}()

		// Check for error while parsing the entry summary.
		if esErr != nil {
			errs = append(errs, esErr)
			continue
		}

		// Check for error when eventually applying the entry.
		eErr := createEntry(entryName, index)
		if eErr != nil {
			errs = append(errs, eErr)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return record, nil
}
