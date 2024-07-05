package parser

import (
	"ktask/ktask"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var parsers = []Parser{
	NewSerialParser(),
	NewParallelParser(1),
	NewParallelParser(2),
	NewParallelParser(4),
	NewParallelParser(8),
	NewParallelParser(15),
	NewParallelParser(50),
}

func TestParseMinimalDocument(t *testing.T) {
	text := `todo`
	for _, p := range parsers {
		rs, _, errs := p.Parse(text)
		require.Nil(t, errs)
		require.Len(t, rs, 1)
		assert.Equal(t, ktask.Todo, rs[0].Stage())
	}
}

func TestParseMultipleRecords(t *testing.T) {
	text := `
todo

done
  1999-06-03 1999-06-02 already done
`
	for _, p := range parsers {
		rs, _, errs := p.Parse(text)
		require.Nil(t, errs)
		require.Len(t, rs, 2)

		assert.Equal(t, ktask.Todo, rs[0].Stage())
		assert.Len(t, rs[0].Entries(), 0)

		assert.Equal(t, ktask.Done, rs[1].Stage())
		assert.Len(t, rs[1].Entries(), 1)
	}
}

func TestParseCompleteRecord(t *testing.T) {
	text := `
todo
    1970-08-29 1970-08-28 needs to be done
    1970-08-27 1970-08-26 needs to be done with
        multiline summary
`
	for _, p := range parsers {
		rs, _, errs := p.Parse(text)
		require.Nil(t, errs)
		require.Len(t, rs, 1)

		r := rs[0]
		assert.Equal(t, ktask.Todo, r.Stage())
		assert.Len(t, r.Entries(), 2)

		assert.Equal(t, time.Date(1970, 8, 29, 0, 0, 0, 0, time.UTC), r.Entries()[0].CreatedAt())
		assert.Equal(t, time.Date(1970, 8, 28, 0, 0, 0, 0, time.UTC), r.Entries()[0].ModifiedAt())
		assert.Equal(t, ktask.Name([]string{"needs to be done"}), r.Entries()[0].Name())
		assert.Equal(t, int(0), r.Entries()[0].Index())

		assert.Equal(t, time.Date(1970, 8, 27, 0, 0, 0, 0, time.UTC), r.Entries()[1].CreatedAt())
		assert.Equal(t, time.Date(1970, 8, 26, 0, 0, 0, 0, time.UTC), r.Entries()[1].ModifiedAt())
		assert.Equal(t, ktask.Name([]string{"needs to be done with", "multiline summary"}), r.Entries()[1].Name())
		assert.Equal(t, int(1), r.Entries()[1].Index())
	}
}

func TestParseEmptyOrBlankDocument(t *testing.T) {
	for _, text := range []string{
		"",
		"    ",
		"\n\n\n\n\n",
		"\n\t     \n \n         ",
	} {
		for _, p := range parsers {
			rs, _, errs := p.Parse(text)
			require.Nil(t, errs)
			require.Len(t, rs, 0)
		}
	}
}

func TestParseWindowsAndUnixLineEndings(t *testing.T) {
	text := "todo\r\n\r\ndone\n\nin progress"
	for _, p := range parsers {
		rs, _, errs := p.Parse(text)
		require.Nil(t, errs)
		require.Len(t, rs, 3)

		assert.Equal(t, ktask.Todo, rs[0].Stage())
		assert.Len(t, rs[0].Entries(), 0)

		assert.Equal(t, ktask.Done, rs[1].Stage())
		assert.Len(t, rs[1].Entries(), 0)

		assert.Equal(t, ktask.InProgress, rs[2].Stage())
		assert.Len(t, rs[1].Entries(), 0)
	}
}

//
// func TestParseUtf8Document(t *testing.T) {
// 	text := `
// 2018-01-04
// 	1h Домашня робота 🏡...
// 	2h Сьогодні я дзвонив
// 		Дімі і складав плани
//
// 2018-01-05
// मुख्य रूपमा काम
// 	10:00-12:30 बगैचा खन्नुहोस्
// 	13:00-15:00 कर घोषणा
//
// 2018-01-06
// 	3h sázet květiny
// 	14:00-? jít na procházku, vynést
// 		odpadky, přines noviny
// `
// 	for _, p := range parsers {
// 		rs, _, errs := p.Parse(text)
// 		require.Nil(t, errs)
// 		require.Len(t, rs, 3)
// 	}
// }
//
// func TestParseMultipleRecordsWhenBlankLineContainsWhitespace(t *testing.T) {
// 	text := "2018-01-01\n    1h\n" + "    \n" + "2019-01-01\n     \n2019-01-02"
// 	for _, p := range parsers {
// 		rs, _, errs := p.Parse(text)
// 		require.Nil(t, errs)
// 		require.Len(t, rs, 3)
// 	}
// }
//
// func TestParseAlternativeFormatting(t *testing.T) {
// 	text := `
// 1999/05/31
// 	8:00-13:00
//
// 1999-05-31
// 	8:00am-1:00pm
// `
// 	for _, p := range parsers {
// 		rs, _, errs := p.Parse(text)
// 		require.Nil(t, errs)
// 		require.Len(t, rs, 2)
//
// 		assert.True(t, rs[0].Date().IsEqualTo(rs[1].Date()))
// 		assert.Equal(t, rs[0].Entries()[0].Duration(), rs[1].Entries()[0].Duration())
// 	}
// }
//
// func TestAcceptTabOrSpacesAsIndentation(t *testing.T) {
// 	for _, x := range []string{
// 		"2000-01-01\n\t8h",
// 		"2000-01-01\n\t8h\n\t15m",
// 		"2000-05-31\n  6h",
// 		"2000-05-31\n  6h\n  20m",
// 		"2000-05-31\n   6h",
// 		"2000-05-31\n    6h",
// 	} {
// 		for _, p := range parsers {
// 			rs, _, errs := p.Parse(x)
// 			require.Nil(t, errs)
// 			require.Len(t, rs, 1)
// 		}
// 	}
// }
//
// func TestParseDocumentSucceedsWithCorrectEntryValues(t *testing.T) {
// 	for _, test := range []struct {
// 		text        string
// 		expectEntry any
// 	}{
// 		// Durations
// 		{"1234-12-12\n\t5h", ktask.NewDuration(5, 0)},
// 		{"1234-12-12\n\t2m", ktask.NewDuration(0, 2)},
// 		{"1234-12-12\n\t2h30m", ktask.NewDuration(2, 30)},
//
// 		// Durations with sign
// 		{"1234-12-12\n\t+5h", ktask.Ɀ_ForceSign_(ktask.NewDuration(5, 0))},
// 		{"1234-12-12\n\t+2h30m", ktask.Ɀ_ForceSign_(ktask.NewDuration(2, 30))},
// 		{"1234-12-12\n\t+2m", ktask.Ɀ_ForceSign_(ktask.NewDuration(0, 2))},
// 		{"1234-12-12\n\t-5h", ktask.NewDuration(-5, -0)},
// 		{"1234-12-12\n\t-2h30m", ktask.NewDuration(-2, -30)},
// 		{"1234-12-12\n\t-2m", ktask.NewDuration(-0, -2)},
//
// 		// Ranges
// 		{"1234-12-12\n\t3:05 - 11:59", ktask.Ɀ_Range_(ktask.Ɀ_Time_(3, 5), ktask.Ɀ_Time_(11, 59))},
// 		{"1234-12-12\n\t22:00 - 24:00", ktask.Ɀ_Range_(ktask.Ɀ_Time_(22, 0), ktask.Ɀ_TimeTomorrow_(0, 0))},
// 		{"1234-12-12\n\t9:00am - 1:43pm", ktask.Ɀ_Range_(ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(9, 00)), ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(13, 43)))},
// 		{"1234-12-12\n\t9:00am-1:43pm", ktask.Ɀ_NoSpaces_(ktask.Ɀ_Range_(ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(9, 00)), ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(13, 43))))},
// 		{"1234-12-12\n\t9:00am-9:05", ktask.Ɀ_NoSpaces_(ktask.Ɀ_Range_(ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(9, 00)), ktask.Ɀ_Time_(9, 05)))},
//
// 		// Ranges with shifted times
// 		{"1234-12-12\n\t9:00am-8:12am>", ktask.Ɀ_NoSpaces_(ktask.Ɀ_Range_(ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(9, 00)), ktask.Ɀ_IsAmPm_(ktask.Ɀ_TimeTomorrow_(8, 12))))},
// 		{"1234-12-12\n\t<22:00 - <24:00", ktask.Ɀ_Range_(ktask.Ɀ_TimeYesterday_(22, 0), ktask.Ɀ_Time_(0, 0))},
// 		{"1234-12-12\n\t<23:30 - 0:10", ktask.Ɀ_Range_(ktask.Ɀ_TimeYesterday_(23, 30), ktask.Ɀ_Time_(0, 10))},
// 		{"1234-12-12\n\t22:17 - 1:00>", ktask.Ɀ_Range_(ktask.Ɀ_Time_(22, 17), ktask.Ɀ_TimeTomorrow_(1, 00))},
// 		{"1234-12-12\n\t22:17   -        1:00>", ktask.Ɀ_Range_(ktask.Ɀ_Time_(22, 17), ktask.Ɀ_TimeTomorrow_(1, 00))},
// 		{"1234-12-12\n\t<23:00-1:00>", ktask.Ɀ_NoSpaces_(ktask.Ɀ_Range_(ktask.Ɀ_TimeYesterday_(23, 00), ktask.Ɀ_TimeTomorrow_(1, 00)))},
// 		{"1234-12-12\n\t<23:00-<23:10", ktask.Ɀ_NoSpaces_(ktask.Ɀ_Range_(ktask.Ɀ_TimeYesterday_(23, 00), ktask.Ɀ_TimeYesterday_(23, 10)))},
// 		{"1234-12-12\n\t12:01>-13:59>", ktask.Ɀ_NoSpaces_(ktask.Ɀ_Range_(ktask.Ɀ_TimeTomorrow_(12, 01), ktask.Ɀ_TimeTomorrow_(13, 59)))},
//
// 		// Open ranges
// 		{"1234-12-12\n\t12:01 - ?", ktask.NewOpenRange(ktask.Ɀ_Time_(12, 1))},
// 		{"1234-12-12\n\t6:45pm - ?", ktask.NewOpenRange(ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(18, 45)))},
// 		{"1234-12-12\n\t6:45pm   -         ?", ktask.NewOpenRange(ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(18, 45)))},
// 		{"1234-12-12\n\t18:45 - ???", ktask.Ɀ_QuestionMarks_(ktask.NewOpenRange(ktask.Ɀ_Time_(18, 45)), 2)},
// 		{"1234-12-12\n\t<3:12-??????", ktask.Ɀ_QuestionMarks_(ktask.Ɀ_NoSpacesO_(ktask.NewOpenRange(ktask.Ɀ_TimeYesterday_(3, 12))), 5)},
// 	} {
// 		for _, p := range parsers {
// 			rs, _, errs := p.Parse(test.text)
// 			require.Nil(t, errs, test.text)
// 			require.Len(t, rs, 1, test.text)
// 			require.Len(t, rs[0].Entries(), 1, test.text)
// 			value := ktask.Unbox(&rs[0].Entries()[0],
// 				func(r ktask.Range) any { return r },
// 				func(d ktask.Duration) any { return d },
// 				func(o ktask.OpenRange) any { return o },
// 			)
// 			assert.Equal(t, test.expectEntry, value, test.text)
// 		}
// 	}
// }
//
// func TestParsesDocumentsWithEntrySummaries(t *testing.T) {
// 	for _, test := range []struct {
// 		text          string
// 		expectEntry   any
// 		expectSummary ktask.EntrySummary
// 	}{
// 		// Single line entries
// 		{"1234-12-12\n\t5h Some remark", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("Some remark")},
// 		{"1234-12-12\n\t3:05 - 11:59 Did this and that", ktask.Ɀ_Range_(ktask.Ɀ_Time_(3, 5), ktask.Ɀ_Time_(11, 59)), ktask.Ɀ_EntrySummary_("Did this and that")},
// 		{"1234-12-12\n\t9:00am-8:12am> Things", ktask.Ɀ_NoSpaces_(ktask.Ɀ_Range_(ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(9, 00)), ktask.Ɀ_IsAmPm_(ktask.Ɀ_TimeTomorrow_(8, 12)))), ktask.Ɀ_EntrySummary_("Things")},
// 		{"1234-12-12\n\t18:45 - ? Just started something", ktask.NewOpenRange(ktask.Ɀ_Time_(18, 45)), ktask.Ɀ_EntrySummary_("Just started something")},
// 		{"1234-12-12\n\t5h    Some remark", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("   Some remark")},
// 		{"1234-12-12\n\t5h\tSome remark", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("Some remark")},
// 		{"1234-12-12\n\t9:00am-9:05 Mixed styles", ktask.Ɀ_NoSpaces_(ktask.Ɀ_Range_(ktask.Ɀ_IsAmPm_(ktask.Ɀ_Time_(9, 00)), ktask.Ɀ_Time_(9, 05))), ktask.Ɀ_EntrySummary_("Mixed styles")},
// 		{"1234-12-12\n\t3:05 - 11:59\tFoo", ktask.Ɀ_Range_(ktask.Ɀ_Time_(3, 5), ktask.Ɀ_Time_(11, 59)), ktask.Ɀ_EntrySummary_("Foo")},
// 		{"1234-12-12\n\t<22:00 - <24:00\tFoo", ktask.Ɀ_Range_(ktask.Ɀ_TimeYesterday_(22, 0), ktask.Ɀ_Time_(0, 0)), ktask.Ɀ_EntrySummary_("Foo")},
// 		{"1234-12-12\n\t22:00 - 24:00\tFoo", ktask.Ɀ_Range_(ktask.Ɀ_Time_(22, 0), ktask.Ɀ_TimeTomorrow_(0, 0)), ktask.Ɀ_EntrySummary_("Foo")},
// 		{"1234-12-12\n\t18:45 - ???       ASDF", ktask.Ɀ_QuestionMarks_(ktask.NewOpenRange(ktask.Ɀ_Time_(18, 45)), 2), ktask.Ɀ_EntrySummary_("      ASDF")},
// 		{"1234-12-12\n\t18:45 - ?\tFoo", ktask.NewOpenRange(ktask.Ɀ_Time_(18, 45)), ktask.Ɀ_EntrySummary_("Foo")},
//
// 		// Multiline-summary entries
// 		{"1234-12-12\n\t5h Some remark\n\t\twith more text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("Some remark", "with more text")},
// 		{"1234-12-12\n\t8:00-9:00 Some remark\n\t\twith more text", ktask.Ɀ_NoSpaces_(ktask.Ɀ_Range_(ktask.Ɀ_Time_(8, 00), ktask.Ɀ_Time_(9, 00))), ktask.Ɀ_EntrySummary_("Some remark", "with more text")},
// 		{"1234-12-12\n\t5h Some remark\n\t\twith\n\t\tmore\n\t\ttext", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("Some remark", "with", "more", "text")},
// 		{"1234-12-12\n  5h Some remark\n    with more text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("Some remark", "with more text")},
// 		{"1234-12-12\n   5h Some remark\n      with more text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("Some remark", "with more text")},
// 		{"1234-12-12\n    5h Some remark\n        with more text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("Some remark", "with more text")},
// 		{"1234-12-12\n    5h Some remark\n        with\n        more\n        text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("Some remark", "with", "more", "text")},
//
// 		// Multiline-summary entries where first summary line is empty
// 		{"1234-12-12\n\t5h\n\t\twith more text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("", "with more text")},
// 		{"1234-12-12\n\t5h \n\t\twith more text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("", "with more text")},
// 		{"1234-12-12\n\t5h  \n\t\twith more text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_(" ", "with more text")},
// 		{"1234-12-12\n\t5h\n\t\t with more text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("", " with more text")},
// 		{"1234-12-12\n\t5h\n\t\t\twith more text", ktask.NewDuration(5, 0), ktask.Ɀ_EntrySummary_("", "\twith more text")},
// 	} {
// 		for _, p := range parsers {
// 			rs, _, errs := p.Parse(test.text)
// 			require.Nil(t, errs, test.text)
// 			require.Len(t, rs, 1, test.text)
// 			require.Len(t, rs[0].Entries(), 1, test.text)
// 			value := ktask.Unbox(&rs[0].Entries()[0],
// 				func(r ktask.Range) any { return r },
// 				func(d ktask.Duration) any { return d },
// 				func(o ktask.OpenRange) any { return o },
// 			)
// 			assert.Equal(t, test.expectEntry, value, test.text)
// 			assert.Equal(t, test.expectSummary, rs[0].Entries()[0].Summary(), test.text)
// 		}
// 	}
// }
//
// func TestMalformedRecord(t *testing.T) {
// 	text := `
// 1999-05-31
// 	5h30m This and that
// Why is there a summary at the end?
// `
// 	for _, p := range parsers {
// 		rs, _, errs := p.Parse(text)
// 		require.Nil(t, rs)
// 		require.NotNil(t, errs)
// 		require.Len(t, errs, 1)
// 		assert.Equal(t, ErrorIllegalIndentation().toErrData(4, 0, 34), toErrData(errs[0]))
// 	}
// }
//
// func TestReportErrorsInHeadline(t *testing.T) {
// 	for _, test := range []struct {
// 		text   string
// 		expect errData
// 	}{
// 		{"Hello 123", ErrorInvalidDate().toErrData(1, 0, 5)},
// 		{" 2020-01-01", ErrorIllegalIndentation().toErrData(1, 0, 11)},
// 		{"   2020-01-01", ErrorIllegalIndentation().toErrData(1, 0, 13)},
// 		{"2020-01-01 ()", ErrorMalformedPropertiesSyntax().toErrData(1, 12, 1)},
// 		{"2020-01-01 (asdf)", ErrorUnrecognisedProperty().toErrData(1, 12, 4)},
// 		{"2020-01-01 (asdf!)", ErrorMalformedShouldTotal().toErrData(1, 12, 4)},
// 		{"2020-01-01 5h30m!", ErrorUnrecognisedTextInHeadline().toErrData(1, 11, 6)},
// 		{"2020-01-01 (5h30m!", ErrorMalformedPropertiesSyntax().toErrData(1, 18, 1)},
// 		{"2020-01-01 (", ErrorMalformedPropertiesSyntax().toErrData(1, 12, 1)},
// 		{"2020-01-01 (5h!) foo", ErrorUnrecognisedTextInHeadline().toErrData(1, 17, 3)},
// 		{"2020-01-01 (5h! asdf)", ErrorUnrecognisedProperty().toErrData(1, 16, 4)},
// 		{"2020-01-01 (5h!!!)", ErrorUnrecognisedProperty().toErrData(1, 15, 2)},
// 	} {
// 		for _, p := range parsers {
// 			rs, _, errs := p.Parse(test.text)
// 			require.Nil(t, rs)
// 			require.NotNil(t, errs)
// 			require.Len(t, errs, 1)
// 			assert.Equal(t, test.expect, toErrData(errs[0]), test.text)
// 		}
// 	}
// }
//
// func TestReportErrorsInSummary(t *testing.T) {
// 	text := `
// 2020-01-01
// This is a summary that contains
//  whitespace at the beginning of the line.
// That is not allowed.
//  Other kinds of blank characters are not allowed there neither.
//  And neither are fake blank lines:
//
// End.
// `
// 	for _, p := range parsers {
// 		rs, _, errs := p.Parse(text)
// 		require.Nil(t, rs)
// 		require.NotNil(t, errs)
// 		require.Len(t, errs, 4)
// 		assert.Equal(t, ErrorMalformedSummary().toErrData(4, 0, 41), toErrData(errs[0]))
// 		assert.Equal(t, ErrorMalformedSummary().toErrData(6, 0, 63), toErrData(errs[1]))
// 		assert.Equal(t, ErrorMalformedSummary().toErrData(7, 0, 34), toErrData(errs[2]))
// 		assert.Equal(t, ErrorMalformedSummary().toErrData(8, 0, 4), toErrData(errs[3]))
// 	}
// }
//
// func TestReportErrorsIfIndentationIsIncorrect(t *testing.T) {
// 	for _, test := range []struct {
// 		text   string
// 		expect errData
// 	}{
// 		// To few characters (that’s actually a malformed summary, though):
// 		{"2020-01-01\n 8h", ErrorMalformedSummary().toErrData(2, 0, 3)},
//
// 		// Not exactly one indentation level:
// 		{"2020-01-01\n\t 8h", ErrorIllegalIndentation().toErrData(2, 0, 4)},
// 		{"2020-01-01\n\t\t8h", ErrorIllegalIndentation().toErrData(2, 0, 4)},
// 		{"2020-01-01\n     8h", ErrorIllegalIndentation().toErrData(2, 0, 7)},
//
// 		// Mixed styles for entries within one record:
// 		{"2020-01-01\n    8h\n\t2h", ErrorIllegalIndentation().toErrData(3, 0, 3)},
// 		{"2020-01-01\n  8h\n\t2h", ErrorIllegalIndentation().toErrData(3, 0, 3)},
// 		{"2020-01-01\n\t8h\n    2h", ErrorIllegalIndentation().toErrData(3, 0, 6)},
// 		{"2020-01-01\n\t8h\n  2h", ErrorIllegalIndentation().toErrData(3, 0, 4)},
// 		{"2020-01-01\n    8h\n  2h", ErrorIllegalIndentation().toErrData(3, 0, 4)},
// 		{"2020-01-01\n  8h\n   2h", ErrorIllegalIndentation().toErrData(3, 0, 5)},
//
// 		// Mixed styles for entry summaries within one record:
// 		{"2020-01-01\n  8h Foo\n\tbar baz", ErrorIllegalIndentation().toErrData(3, 0, 8)},
// 		{"2020-01-01\n    8h Foo\n       bar baz", ErrorIllegalIndentation().toErrData(3, 0, 14)},
// 		{"2020-01-01\n    8h Foo\n      bar baz", ErrorIllegalIndentation().toErrData(3, 0, 13)},
// 		{"2020-01-01\n    8h Foo\n    \tbar baz", ErrorIllegalIndentation().toErrData(3, 0, 12)},
// 		{"2020-01-01\n   8h Foo\n     bar baz", ErrorIllegalIndentation().toErrData(3, 0, 12)},
// 		{"2020-01-01\n  8h Foo\n   bar baz", ErrorIllegalIndentation().toErrData(3, 0, 10)},
// 		{"2020-01-01\n  8h\n  8h Foo\n   bar baz", ErrorIllegalIndentation().toErrData(4, 0, 10)},
// 	} {
// 		for _, p := range parsers {
// 			rs, _, errs := p.Parse(test.text)
// 			require.Nil(t, rs, test.text)
// 			require.NotNil(t, errs, test.text)
// 			require.Len(t, errs, 1, test.text)
// 			assert.Equal(t, test.expect, toErrData(errs[0]), test.text)
// 		}
// 	}
// }
//
// func TestAcceptMixingIndentationStylesAcrossDifferentRecords(t *testing.T) {
// 	text := `
// 2020-01-01
//   4h This is two spaces
//   2h So is this
//
// 2020-01-02
//     6h This is 4 spaces
//
// 2020-01-03
// 	12m This is a tab
// `
// 	for _, p := range parsers {
// 		rs, _, errs := p.Parse(text)
// 		require.Nil(t, errs)
// 		require.Len(t, rs, 3)
// 	}
// }
//
// func TestReportErrorsInEntries(t *testing.T) {
// 	for _, test := range []struct {
// 		text   string
// 		expect errData
// 	}{
// 		// Malformed syntax
// 		{"2020-01-01\n\t5h1", ErrorMalformedEntry().toErrData(2, 1, 3)},
// 		{"2020-01-01\n\tasdf Test 123", ErrorMalformedEntry().toErrData(2, 1, 4)},
// 		{"2020-01-01\n\t15:30", ErrorMalformedEntry().toErrData(2, 6, 1)},
// 		{"2020-01-01\n\t08:00-", ErrorMalformedEntry().toErrData(2, 7, 1)},
// 		{"2020-01-01\n\t08:00-asdf", ErrorMalformedEntry().toErrData(2, 7, 4)},
// 		{"2020-01-01\n\t08:00 - ?asdf", ErrorMalformedEntry().toErrData(2, 10, 4)},
// 		{"2020-01-01\n\t-18:00", ErrorMalformedEntry().toErrData(2, 1, 6)},
// 		{"2020-01-01\n\t5h Test\n\t15:30 Foo Bar Baz", ErrorMalformedEntry().toErrData(3, 7, 1)},
// 		{"2020-01-01\n\t5h Hello\n\t\tFoo\n\t15:30 Foo Bar Baz", ErrorMalformedEntry().toErrData(4, 7, 1)},
// 		{"2020-01-01\n\t12:76 - 13:00", ErrorMalformedEntry().toErrData(2, 1, 5)},
// 		{"2020-01-01\n\t12:00 - 44:00", ErrorMalformedEntry().toErrData(2, 9, 5)},
// 		{"2020-01-01\n\t23:00> - 25:61>", ErrorMalformedEntry().toErrData(2, 10, 6)},
// 		{"2020-01-01\n\t12:00> - 24:00>", ErrorMalformedEntry().toErrData(2, 10, 6)},
//
// 		// Logical errors
// 		{"2020-01-01\n\t08:00- ?\n\t09:00 - ?", ErrorDuplicateOpenRange().toErrData(3, 1, 9)},
// 		{"2020-01-01\n\t15:00 - 14:00", ErrorIllegalRange().toErrData(2, 1, 13)},
// 		{"2020-01-01\n\t15:00 - 14:00", ErrorIllegalRange().toErrData(2, 1, 13)},
// 	} {
// 		for _, p := range parsers {
// 			rs, _, errs := p.Parse(test.text)
// 			require.Nil(t, rs, test.text)
// 			require.NotNil(t, errs, test.text)
// 			require.Len(t, errs, 1, test.text)
// 			assert.Equal(t, test.expect, toErrData(errs[0]), test.text)
// 		}
// 	}
// }
//
// func TestParseLongDocumentWithMultipleErrors(t *testing.T) {
// 	text := `
// 2019-08-15
//     16:00-19:41 Something
//     20:02-?
//
// 2019-08-16
//     Entry without value
//     8h
//     -12m Break
//
// 2019-08-17 (8h)
// Record summary
//     11:00-?
//       Open range
//
//
// 2019-08-38
// What date is this?!?
// `
// 	for _, p := range parsers {
// 		_, _, errs := p.Parse(text)
// 		require.Len(t, errs, 4)
// 		assert.Equal(t, ErrorMalformedEntry().toErrData(7, 4, 5), toErrData(errs[0]))
// 		assert.Equal(t, ErrorUnrecognisedProperty().toErrData(11, 12, 2), toErrData(errs[1]))
// 		assert.Equal(t, ErrorIllegalIndentation().toErrData(14, 0, 16), toErrData(errs[2]))
// 		assert.Equal(t, ErrorInvalidDate().toErrData(17, 0, 10), toErrData(errs[3]))
// 	}
// }
