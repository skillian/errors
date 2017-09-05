// Errors is a package that helps with propogating error messages up the
// call stack.  When unhandled, the errors.Error struct prints the wrapped
// error message, the error cause and its context (the error being handled
// when the errors.Error.Err error was produced).
package errors

import (
	"fmt"
)

// Error bundles a builtin error with its causing error (Cause) or the error
// that was being handled at the time that the Err error occurred (Context).
//
// This API attempts to match that of Python's __cause__ and __context__
// Exception attributes.
type Error struct {
	// Err contains the error message that's being wrapped.
	Err error
	// Cause contains an error value that directly resulted in Err being
	// returned
	Cause error
	// Context contains the error value that was being handled when Err
	// was generated
	Context error
}

// ErrorfWithCause creates an Error with a formatted error string and  then
// states the error's cause within an Error struct.
func ErrorfWithCause(cause error, format string, args... interface{}) Error {
	return Error{Err: ErrorMessage{Fmt: format, Args: args}, Cause: cause}
}

// ErrorfWithContext creates an Error with a formatted error string and  then
// states the error's context within an Error struct.
func ErrorfWithContext(context error, format string, args... interface{}) Error {
	return Error{Err: ErrorMessage{Fmt: format, Args: args}, Context: context}
}

// ErrorfWithCauseAndContext creates an Error with a formatted error string and
// then states the error's context within an Error struct.
func ErrorfWithCauseAndContext(cause, context error, format string, args... interface{}) Error {
	return Error{
		Err: ErrorMessage{
			Fmt: format,
			Args: args,
		},
		Cause: cause,
		Context: context,
	}
}

// Error implements the builtin error interface that includes information about
// the Err Cause and the Context.
func (e Error) Error() string {
	var ca string
	if e.Cause == nil {
		ca = ""
	} else {
		ca = fmt.Sprintf("\nCause:  %v", e.Cause)
	}
	var co string
	if e.Context == nil {
		co = ""
	} else {
		ca = fmt.Sprintf("\nContext:  %v", e.Context)
	}
	return fmt.Sprintf("%v%v%v", e.Err, ca, co)
}

// Message defines a message with parameters
type Message struct {
	// Fmt holds a string with its formatting parameters
	Fmt string
	// Args are the parameters to the Fmt string in the message
	Args []interface{}
}

// String implements the io.Stringer protocol so that printing strings formats
// the arguments.
func (m Message) String() string {
	return fmt.Sprintf(m.Fmt, m.Args...)
}

// ErrorMessage is the same thing as a message with a separate type so that
// it can be treated as an error.
type ErrorMessage Message

// Error implements the builtin error interface to treat messages as error
// messages
func (m ErrorMessage) Error() string {
	return Message(m).String()
}


