package errors

// WrapDeferred wraps a deferred function and captures its error into the
// provided error pointer.
func WrapDeferred(pErr *error, f func() error) {
    if pErr == nil {
        panic("pErr error pointer is nil")
    }
    err := f()
    if err != nil {
        if *pErr == nil {
            *pErr = err
        } else {
            *pErr = Error{
                Err: err,
                Context: *pErr,
            }
        }
    }
}
