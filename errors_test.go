package errors_test

import (
	"errors"
	"testing"

	skerrors "github.com/skillian/errors"
)

var (
	HelloWorld = "Hello, world"
	HelloWorldWithContext = `Hello, world
	Context:  Hello, world`
	HelloWorldWithCause = `Hello, world
	Cause:  Hello, world`
	HelloWorldWithCauseAndContext = `Hello, world
	Cause:  Hello, world
	Context:  Hello, world`

	ErrHelloWorld = errors.New(HelloWorld)
)

func TestExact(t *testing.T) {
	t.Parallel()
	e := skerrors.Error{Err: ErrHelloWorld}
	es := e.Error()
	if es != HelloWorld {
		t.Errorf(
			"Error without Cause or Context should be the same string!(got: %q, expected %q)",
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
