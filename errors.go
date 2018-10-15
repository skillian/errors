// Package errors helps with propogating error messages up the call stack.
// When unhandled, the errors.Error struct prints the wrapped error message,
// the error cause and its context (the error being handled when the
// errors.Error.Err error was produced).
package errors

import (
    "fmt"
    "runtime"
    "strings"
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

    // pcs holds a slice of program counters that can be turned into a stack
    // trace.
    pcs []uintptr
}

// Cause gets the root cause of the given error.  If the error is an
// errors.Error or *errors.Error, Cause attempts to recursively unpack the
// errors.Error.Cause until it gets to a Cause that is not an errors.Error or
// *errors.Error.
func Cause(err error) error {
    for {
        e, ok := err.(Error)
        if !ok {
            pe, ok := err.(*Error)
            if !ok {
                return err
            }
            e = *pe
        }
        if e.Cause == nil {
            return e.Err
        }
        err = e.Cause
    }
}

// WrapDefer wraps a deferred function to ensure its returned error value is
// not discarded.
func WrapDefer(pe *error, deferred func() error) {
    if pe == nil {
        panic("WrapDefer requires a pointer to an error")
    }
    err := deferred()
    if err != nil {
        if (*pe) != nil {
            err = CreateError(err, nil, *pe, 0)
        }
        *pe = err
    }
}

// CreateError creates and returns an Error object after initializing its
// program counters slice starting at the frame specified by skip where a value
// of 0 will start from the caller of CreateError
func CreateError(err, cause, context error, skip int) Error {
    e := Error{
        Err:     err,
        Cause:   cause,
        Context: context,
    }
    setErrorPCs(skip+1, &e)
    return e
}

// Errorf returns an error message without a Cause or Context
func Errorf(format string, args... interface{}) Error {
    return Error{Err: ErrorMessage{Fmt: format, Args: args}}
}

// ErrorfWithCause creates an Error with a formatted error string and  then
// states the error's cause within an Error struct.
func ErrorfWithCause(cause error, format string, args ...interface{}) Error {
    return CreateError(ErrorMessage{Fmt: format, Args: args}, cause, nil, 1)
}

// ErrorfWithContext creates an Error with a formatted error string and  then
// states the error's context within an Error struct.
func ErrorfWithContext(context error, format string, args ...interface{}) Error {
    return CreateError(ErrorMessage{Fmt: format, Args: args}, nil, context, 1)
}

// ErrorfWithCauseAndContext creates an Error with a formatted error string and
// then states the error's context within an Error struct.
func ErrorfWithCauseAndContext(cause, context error, format string, args ...interface{}) Error {
    return CreateError(ErrorMessage{Fmt: format, Args: args}, cause, context, 1)
}

// Error implements the builtin error interface that includes information about
// the Err Cause and the Context.
func (e Error) Error() string {
    var ca string
    if e.Cause == nil {
        ca = ""
    } else {
        ca = fmt.Sprintf("\n    Cause:  %v", e.Cause)
    }
    var co string
    if e.Context == nil {
        co = ""
    } else {
        co = fmt.Sprintf("\n    Context:  %v", e.Context)
    }
    stackTrace := formatStackTrace(e)
    return fmt.Sprintf(
        "%v%v%v%v",
        e.Err,
        stackTrace,
        ca,
        co)
}

func setErrorPCs(skip int, e *Error) {
    var pcs [32]uintptr
    count := runtime.Callers(skip+2, pcs[:])
    e.pcs = make([]uintptr, count)
    for i, pc := range pcs[:count] {
        e.pcs[i] = pc
    }
}

func formatStackTrace(e Error) string {
    if e.pcs == nil {
        return ""
    }
    formattedFrames := make([]string, 0, len(e.pcs))
    frames := runtime.CallersFrames(e.pcs)
    for {
        frame, more := frames.Next()
        formattedFrames = append(formattedFrames, fmt.Sprintf(
            "   at %v in %v, line %d",
            frame.Function,
            frame.File,
            frame.Line))
        if !more {
            break
        }
    }
    return strings.Join(formattedFrames, "\n")
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

// Errorf works like fmt.Errorf but returns an ErrorMessage
func Errorf(format string, args ...interface{}) ErrorMessage {
    return ErrorMessage(Message{Fmt: format, Args: args})
}

// Error implements the builtin error interface to treat messages as error
// messages
func (m ErrorMessage) Error() string {
    return Message(m).String()
}
