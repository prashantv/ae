// Package ae (AppError) is used to wrap errors with additional information and stack traces.
// ae.Err can represent a root underlying error and annotate that error with additional
// errors. The stack trace where the original error was converted to an ae.Err is recorded.
package ae

import (
	"bytes"
	"fmt"
)

var (
	// IncludeStackInError controls whether to include the stack trace in the Error string.
	IncludeStackInError = true

	// PrintToLog controls errors should be printed to the log when Err.String is called.
	PrintToLog = true
)

// Err is the interface used for application errors with their cause and location.
type Err interface {
	error

	// Errors returns the list of errors in least recent to most recent order.
	Errors() []error

	// First returns the least recent error.
	First() error

	// Last returns the most recent error.
	Last() error

	// Stack returns the stack trace from when the error was wrapped.
	Stack() string
}

type appError struct {
	underlying []error
	stack      []uintptr
	frameCache []stackFrame
}

func (ae *appError) Error() string {
	if PrintToLog {
		ae.PrintTolog()
	}
	if IncludeStackInError {
		return ae.errorMsgs() + StackSeparator + ae.Stack()
	}
	return ae.errorMsgs()
}

func (ae *appError) First() error {
	return ae.underlying[0]
}

func (ae *appError) Last() error {
	return ae.underlying[len(ae.underlying)-1]
}

func (ae *appError) Errors() []error {
	return ae.underlying
}

// errorMsgs combines all the underlying errors from most recent to least recent.
func (ae *appError) errorMsgs() string {
	var buf bytes.Buffer
	// messages should be read in stack order: last first
	for i := len(ae.underlying) - 1; i >= 0; i-- {
		e := ae.underlying[i]
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		buf.WriteString(e.Error())
	}
	return buf.String()
}

// Errorf creates an error using fmt.Errorf, and then wraps that error.
func Errorf(format string, a ...interface{}) Err {
	return wrapF(fmt.Errorf(format, a...))
}

// Wrap wraps the underlying error and adds stack information.
func Wrap(underlying error) Err {
	return wrapF(underlying)
}

// Wrapf wraps the underlying error and a new error to the underlying errors.
func Wrapf(underlying error, format string, a ...interface{}) Err {
	// Leave early to avoid formatting the message if there's no actual error.
	if underlying == nil {
		return nil
	}
	return wrapF(underlying, fmt.Errorf(format, a...))
}

func wrapF(underlying error, extraErrors ...error) Err {
	if underlying == nil {
		return nil
	}

	// If the error we are wrapping is already an an appError, then we
	// just need to add a new error to the underlying error list.
	if aerr, ok := underlying.(*appError); ok {
		if len(extraErrors) > 0 {
			aerr.underlying = append(aerr.underlying, extraErrors...)
		}
		return aerr
	}

	underlyingErrs := []error{underlying}
	underlyingErrs = append(underlyingErrs, extraErrors...)
	return &appError{
		underlying: underlyingErrs,
		// skip is 2 since we want to ignore the wrapF and the caller of wrapF.
		stack: getStackPC(2 /* skip */),
	}
}
