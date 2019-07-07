package errors

// UnexpectedType is returned when a value is received in a context where a
// value of another type was expected.
type UnexpectedType struct {
	// Message contains the error message value.
	Message

	// Expected contains a value of the expected type.  It's value may
	// or may not mean anything depending on the message included with
	// the error, but its type actually matters and is included in the
	// error message.
	//
	// If Expected's value's type is AnyType, then error messages state
	// that the expected type is any of the given slice of
	Expected interface{}

	// Actual contains a value of the expected type.  It's value may
	// or may not mean anything depending on the message included with
	// the error, but its type actually matters and is included in the
	// error message.
	Actual interface{}

	// slice is an array of the Expected and Actual values to avoid an
	// allocation outside of the UnexpectedType value itself.
	slice [5]interface{}
}

const errUnexpectedTypeDefaultString = "%T: expected: %v (type: %T), but found: %v (type: %T)"

// NewUnexpectedType returns an UnexpectedType error from the given pair of
// values.
func NewUnexpectedType(expected, actual interface{}) *UnexpectedType {
	err := &UnexpectedType{
		Message:  Message{},
		Expected: expected,
		Actual:   actual,
		slice: [...]interface{}{
			nil, // the UnexpectdType error itself,
			expected,
			expected,
			actual,
			actual,
		},
	}
	// set the first element of the slice so the error formatting works
	// correctly.
	err.slice[0] = err
	err.Message = Message{
		Fmt: errUnexpectedTypeDefaultString,
		// slicing an array keeps the array itself instead of allocating
		// a new slice:
		Args: err.slice[:],
	}
	return err
}

// AnyType is used in TypeError to indicate that the expected type is any
// of the given types.
type AnyType []interface{}
