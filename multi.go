package errors

import (
	"errors"
	"math/bits"
	"strconv"
	"strings"
)

type Errors struct {
	errors []error
}

var _ interface {
	error
	As(target interface{}) bool
	Is(target error) bool
	Unwrap() error
} = Errors{}

// Aggregate multiple errors together into a single error.  Any nil
// errors are ommitted.  Any aggregate errors are "flattened" into this
// error.
func Aggregate(errs ...error) error {
	es := Errors{errors: make([]error, 0, 1<<bits.Len(uint(len(errs))))}
	es.appendErrors(errs)
	if len(es.errors) == 0 {
		return nil
	}
	return es
}

func (es *Errors) appendErrors(errs []error) {
	for _, err := range errs {
		if err == nil {
			continue
		}
		if es2, ok := err.(Errors); ok {
			es.appendErrors(es2.errors)
			continue
		}
		if es2, ok := err.(*Errors); ok {
			es.appendErrors(es2.errors)
			continue
		}
		es.errors = append(es.errors, err)
	}
}

func (es Errors) As(target interface{}) bool {
	if es2, ok := target.(*Errors); ok {
		*es2 = es
		return true
	}
	if es2, ok := target.(**Errors); ok {
		if es2 == nil {
			return false
		}
		if *es2 == nil {
			*es2 = &Errors{}
		}
		*(*es2) = es
		return true
	}
	for _, err := range es.errors {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (es Errors) Error() string {
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(len(es.errors)))
	sb.WriteString(" errors: ")
	for i, err := range es.errors {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

func (es Errors) Is(target error) bool {
	if es2, ok := target.(*Errors); ok {
		return es.errorsEqual(es2.errors)
	}
	if es2, ok := target.(Errors); ok {
		return es.errorsEqual(es2.errors)
	}
	for _, err := range es.errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (es Errors) errorsEqual(errs []error) bool {
	if len(es.errors) != len(errs) {
		return false
	}
	for i, err := range es.errors {
		if !errors.Is(err, errs[i]) {
			return false
		}
	}
	return true
}

func (es Errors) Unwrap() error {
	es2 := Errors{errors: es.errors[1:]}
	if len(es2.errors) == 0 {
		return nil
	}
	return es2
}
