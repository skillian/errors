package errors_test

import (
	"errors"
	"testing"

	skerrors "github.com/skillian/errors"
)

var (
	HelloWorld            = "Hello, world"
	HelloWorldWithContext = `Hello, world
Context:  Hello, world`
	HelloWorldWithCause = `Hello, world
Cause:  Hello, world`
	HelloWorldWithCauseAndContext = `Hello, world
Cause:  Hello, world
Context:  Hello, world`
	HelloWorldWithStackTrace = `Hello, world
github.com/skillian/errors_test.TestWithStackTrace
	github.com/skillian/errors/errors_test.go:60
testing.tRunner
`

	ErrHelloWorld = errors.New(HelloWorld)
)

func TestExact(t *testing.T) {
	t.Parallel()
	e := skerrors.Error{Err: ErrHelloWorld}
	es := e.Error()
	if es != HelloWorld {
		t.Errorf(
			"Error without Cause or Context should be the same string! (got: %q, expected %q)",
			es,
			HelloWorld)
	}
}

func TestWithContext(t *testing.T) {
	t.Parallel()
	e := skerrors.Error{Err: ErrHelloWorld, Context: ErrHelloWorld}
	es := e.Error()
	if es != HelloWorldWithContext {
		t.Errorf("Error with Context should be: %q, got: %q", HelloWorldWithContext, es)
	}
}

func TestWithCause(t *testing.T) {
	t.Parallel()
	e := skerrors.Error{Err: ErrHelloWorld, Cause: ErrHelloWorld}
	es := e.Error()
	if es != HelloWorldWithCause {
		t.Errorf("Error with Cause should be: %q, got: %q", HelloWorldWithCause, es)
	}
}

func TestWithStackTrace(t *testing.T) {
	t.Parallel()
	e := skerrors.CreateError(ErrHelloWorld, nil, nil, 0)
	es := e.Error()
	i := findDiffIndex(HelloWorldWithStackTrace, es)
	if i != -1 {
		start := i - 10
		if start < 0 {
			start = 0
		}
		end := i + 10
		if i >= len(es) {
			end = len(es)
		}
		t.Errorf("Error with stack trace should be: %q, got: %q", HelloWorldWithStackTrace, es[start:end])
	}
}

func findDiffIndex(a, b string) int {
	for i, r := range []byte(a) {
		if b[i] != r {
			return i
		}
	}
	return -1
}
