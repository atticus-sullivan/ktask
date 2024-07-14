package ktask

import (
	"errors"
	"fmt"
	"strings"

	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/parser/txt"
)

var Reflower = tf.NewReflower(80, "\n")

func NewParserErrors(errs []txt.Error) ParserErrors {
	return parserErrors{errs}
}

type ParserErrors interface {
	Error
	All() []txt.Error
}

// Error is a representation of an application error.
type Error interface {
	// Error returns the error message.
	Error() string

	// Details returns additional details, such as a hint how to solve the problem.
	Details() string

	// Original returns the original underlying error, if it exists.
	Original() error

	// Code returns the error code.
	Code() Code
}

type Code int

type parserErrors struct {
	errors []txt.Error
}

func (pe parserErrors) Error() string {
	return fmt.Sprintf("%d parsing error(s)", len(pe.errors))
}

func (pe parserErrors) Details() string {
	return fmt.Sprintf("%d parsing error(s)", len(pe.errors))
}

func (pe parserErrors) Original() error {
	return nil
}

func (pe parserErrors) Code() Code {
	return LOGICAL_ERROR
}

func (pe parserErrors) All() []txt.Error {
	return pe.errors
}

const (
	// GENERAL_ERROR should be used for generic, otherwise unspecified errors.
	GENERAL_ERROR Code = iota + 1

	// NO_INPUT_ERROR should be used if no input was specified.
	NO_INPUT_ERROR

	// NO_TARGET_FILE should be used if no target file was specified.
	NO_TARGET_FILE

	// IO_ERROR should be used for errors during I/O processes.
	IO_ERROR

	// CONFIG_ERROR should be used for config-folder-related problems.
	CONFIG_ERROR

	// NO_SUCH_BOOKMARK_ERROR should be used if the specified an unknown bookmark name.
	NO_SUCH_BOOKMARK_ERROR

	// NO_SUCH_FILE should be used if the specified file does not exit.
	NO_SUCH_FILE

	// LOGICAL_ERROR should be used syntax or logical violations.
	LOGICAL_ERROR
)

// PrettifyParsingError turns a parsing error into a coloured and well-structured form.
func PrettifyParsingError(err ParserErrors, styler tf.Styler) error {
	message := ""
	INDENT := "    "
	for _, e := range err.All() {
		message += "\n"
		message += fmt.Sprintf(
			styler.Props(tf.StyleProps{Background: tf.RED, Color: tf.RED}).Format("[")+
				styler.Props(tf.StyleProps{Background: tf.RED, Color: tf.TEXT_INVERSE}).Format("SYNTAX ERROR")+
				styler.Props(tf.StyleProps{Background: tf.RED, Color: tf.RED}).Format("]")+
				styler.Props(tf.StyleProps{Color: tf.RED}).Format(" in line %d"),
			e.LineNumber(),
		)
		if e.Origin() != "" {
			message += fmt.Sprintf(
				styler.Props(tf.StyleProps{Color: tf.RED}).Format(" of file %s"),
				e.Origin(),
			)
		}
		message += "\n"
		message += fmt.Sprintf(
			styler.Props(tf.StyleProps{Color: tf.SUBDUED}).Format(INDENT+"%s"),
			// Replace all tabs with one space each, otherwise the carets might
			// not be in line with the text anymore (since we can’t know how wide
			// a tab is).
			strings.Replace(e.LineText(), "\t", " ", -1),
		) + "\n"
		message += fmt.Sprintf(
			styler.Props(tf.StyleProps{Color: tf.RED}).Format(INDENT+"%s%s"),
			strings.Repeat(" ", e.Position()), strings.Repeat("^", e.Length()),
		) + "\n"
		message += fmt.Sprintf(
			styler.Props(tf.StyleProps{Color: tf.YELLOW}).Format("%s"),
			Reflower.Reflow(e.Message(), []string{INDENT}),
		) + "\n"
	}
	return errors.New(message)
}

func NewErrorWithCode(code Code, message string, details string, original error) Error {
	return AppError{code, message, details, original}
}

func (e AppError) Error() string {
	return e.message
}

func (e AppError) Details() string {
	return e.details
}

func (e AppError) Original() error {
	return e.original
}

func (e AppError) Code() Code {
	return e.code
}

type AppError struct {
	code     Code
	message  string
	details  string
	original error
}

func NewError(message string, details string, original error) Error {
	return NewErrorWithCode(GENERAL_ERROR, message, details, original)
}

// PrettifyAppError prints app errors including details.
func PrettifyAppError(err Error, isDebug bool) error {
	message := "Error: " + err.Error() + "\n"
	message += Reflower.Reflow(err.Details(), nil)
	if isDebug && err.Original() != nil {
		message += "\n\nOriginal Error:\n" + err.Original().Error()
	}
	return errors.New(message)
}
