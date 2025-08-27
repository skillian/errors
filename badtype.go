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
	slice [3]interface{}
}

const errUnexpectedTypeDefaultString = "%[1]T: expected: %[2]v (type: %[2]T), but found: %[3]v (type: %[3]T)"

// NewUnexpectedType returns an UnexpectedType error from the given pair of
// values.
func NewUnexpectedType(expected, actual interface{}) *UnexpectedType {
	err := &UnexpectedType{
		Message: Message{
			Fmt: errUnexpectedTypeDefaultString,
		},
		Expected: expected,
		Actual:   actual,
		slice: [...]interface{}{
			nil, // the UnexpectdType error itself,
			expected,
			actual,
		},
	}
	err.slice[0] = err
	err.Message.Args = err.slice[:]
	return err
}

// AnyType is used in TypeError to indicate that the expected type is any
// of the given types.
type AnyType []interface{}
