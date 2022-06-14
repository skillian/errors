// Package errors helps with propogating error messages up the call stack.
// When unhandled, the errors.Error struct prints the wrapped error message,
// the error cause and its context (the error being handled when the
// errors.Error.Err error was produced).
package errors

import (
	goerrors "errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

const pathSep = string(byte(os.PathSeparator))

var (
	goPath = func(v string) string {
		if len(v) >= len(pathSep) {
			if v[len(v)-len(pathSep):] != pathSep {
				return v + pathSep
			}
		}
		return v
	}(path.Clean(os.Getenv("GOPATH")))
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

// As forwards the call to the go errors package.
func As(err error, target interface{}) bool {
	return goerrors.As(err, target)
}

// Is forwards the call to the go errors package.
func Is(err, target error) bool {
	return goerrors.Is(err, target)
}

// New calls the Go errors package's New function so that you don't have to
// import both packages.
func New(text string) error {
	return goerrors.New(text)
}

// WrapDeferred wraps a deferred function to ensure its returned error value
// isn't discarded.
func WrapDeferred(pe *error, deferred func() error) {
	if err := deferred(); err != nil {
		if pe == nil {
			*pe = err
		} else {
			*pe = CreateError(err, nil, *pe, 0)
		}
	}
	return
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
func Errorf(format string, args ...interface{}) Error {
	return CreateError(Message{Fmt: format, Args: args}, nil, nil, 1)
}

// ErrorfWithCause creates an Error with a formatted error string and  then
// states the error's cause within an Error struct.
func ErrorfWithCause(cause error, format string, args ...interface{}) Error {
	return CreateError(Message{Fmt: format, Args: args}, cause, nil, 1)
}

// ErrorfWithContext creates an Error with a formatted error string and  then
// states the error's context within an Error struct.
func ErrorfWithContext(context error, format string, args ...interface{}) Error {
	return CreateError(Message{Fmt: format, Args: args}, nil, context, 1)
}

// ErrorfWithCauseAndContext creates an Error with a formatted error string and
// then states the error's context within an Error struct.
func ErrorfWithCauseAndContext(cause, context error, format string, args ...interface{}) Error {
	return CreateError(Message{Fmt: format, Args: args}, cause, context, 1)
}

// Error implements the builtin error interface that includes information about
// the Err Cause and the Context.
func (e Error) Error() string {
	var ca string
	if e.Cause != nil {
		ca = fmt.Sprintf("Cause:  %v", e.Cause)
	}
	var co string
	if e.Context != nil {
		co = fmt.Sprintf("Context:  %v", e.Context)
	}
	stackTrace := formatStackTrace(e)
	return strings.Join([]string{
		e.Err.Error(),
		stackTrace,
		ca,
		co,
	}, "\n")
}

func (e Error) As(target interface{}) bool {
	if p, ok := target.(*Error); ok {
		*p = e
		return true
	}
	return As(e.Err, target) || (e.Cause != nil && As(e.Cause, target)) || (e.Context != nil && As(e.Context, target))
}

func (e Error) Is(err error) bool {
	if e2, ok := err.(Error); ok {
		return Is(e.Err, e2.Err) && ((e.Cause == nil && e2.Cause == nil) || (e.Cause != nil && e2.Cause != nil && Is(e.Cause, e2.Cause))) && ((e.Context == nil && e2.Context == nil) || (e.Context != nil && e2.Context != nil && Is(e.Context, e2.Context)))
	}
	return Is(e.Err, err) || (e.Cause != nil && Is(e.Cause, err)) || (e.Context != nil && Is(e.Context, err))
}

func setErrorPCs(skip int, e *Error) {
	var cache [32]uintptr
	pcs := cache[:]
	for count := 0; ; {
		count += runtime.Callers(skip+2+count, pcs[count:])
		if count < len(pcs) {
			pcs = pcs[:count]
			break
		}
		pcs = append(pcs, 0)
		pcs = pcs[:cap(pcs)]
	}
	e.pcs = make([]uintptr, len(pcs))
	copy(e.pcs, pcs)
}

func formatStackTrace(e Error) string {
	if e.pcs == nil {
		return ""
	}
	formattedFrames := make([]string, 0, len(e.pcs))
	frames := runtime.CallersFrames(e.pcs)
	for {
		frame, more := frames.Next()
		file := frame.File
		if strings.HasPrefix(file, goPath) {
			file = file[len(goPath):]
		}
		formattedFrames = append(formattedFrames, fmt.Sprintf(
			"%v\n\t%v:%d",
			frame.Function,
			file,
			frame.Line))
		if !more {
			break
		}
	}
	return strings.Join(formattedFrames, "\n")
}

// Message defines an error message with parameters
type Message struct {
	// Fmt holds a string with its formatting parameters
	Fmt string

	// Args are the parameters to the Fmt string in the message
	Args []interface{}
}

// Error implements the error interface so that printing strings formats
// the arguments.
func (m Message) Error() string {
	return fmt.Sprintf(m.Fmt, m.Args...)
}
