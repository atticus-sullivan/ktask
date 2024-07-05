package parser

import "github.com/jotaen/klog/klog/parser/txt"

type HumanError struct {
	code    string
	title   string
	details string
}

func (e HumanError) New(b txt.Block, line int, start int, length int) txt.Error {
	return txt.NewError(b, line, start, length, e.code, e.title, e.details)
}

func ErrorInvalidStage() HumanError {
	return HumanError{
		"ErrorInvalidStage",
		"Invalid stage",
		"The highlighted value is not recognised as stage. " +
			"Please check if the stage is formatted correctly",
	}
}

func ErrorIllegalIndentation() HumanError {
	return HumanError{
		"ErrorIllegalIndentation",
		"Unexpected indentation",
		"Please correct the indentation of this line. Indentation must be 2-4 spaces or one tab. " +
			"You cannot mix different indentation styles within the same record.",
	}
}

func ErrorUnrecognisedTextInHeadline() HumanError {
	return HumanError{
		"ErrorUnrecognisedTextInHeadline",
		"Malformed headline",
		"The highlighted text in the headline is not recognised. " +
			"Please make sure to surround the should-total with parentheses, e.g.: (5h!) " +
			"You generally cannot put arbitrary text into the headline.",
	}
}

func ErrorMalformedSummary() HumanError {
	return HumanError{
		"ErrorMalformedSummary",
		"Malformed summary",
		"Summary lines cannot start with blank characters, such as non-breaking spaces.",
	}
}

func ErrorMalformedEntry() HumanError {
	return HumanError{
		"ErrorMalformedEntry",
		"Malformed entry",
		"Please review the syntax of the entry. " +
			"It must start with a duration or a time range. " +
			"Valid examples would be: 3h20m or 8:00-10:00 or 8:00-? " +
			"or <23:00-6:00 or 18:00-0:30>",
	}
}
