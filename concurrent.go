package errors

import (
	"strings"
	"sync"
)

// Concurrent represents a collection of errors from concurrent
// Goroutines.
type Concurrent struct {
	mutex  sync.Mutex
	errors []error
}

// Add one or more errors to the Concurrent collection
func (e *Concurrent) Add(errs ...error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.errors = append(e.errors, errs...)
}

// Err bundles the errors in the Concurrent slice together and returns
// a single error.  If there are no errors in the Concurrent slice,
// a nil error is returned.
func (e *Concurrent) Err() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if len(e.errors) == 0 {
		return nil
	}
	return Errorf(
		"%d errors:\n    %v",
		len(e.errors),
		unjoined(e.errors))
}

type unjoined []error

// String implements fmt.Stringer and allows a collection of errors to be
// joined lazily.  Only when a ConcurrentErrors.Err result is formatted are
// the errors joined
func (u unjoined) String() string {
	strs := make([]string, len(u))
	for i, err := range u {
		strs[i] = err.Error()
	}
	return strings.Join(strs, "\n    ")
}
