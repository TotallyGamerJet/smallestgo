package runtime

type Error interface {
	error

	// RuntimeError is a no-op function but
	// serves to distinguish types that are run time
	// errors from ordinary errors: a type is a
	// run time error if it has a RuntimeError method.
	RuntimeError()
}

type errorString string

type plainError string

func (e plainError) RuntimeError() {}

func (e plainError) Error() string {
	return ""
}

//go:nosplit
func itoa(buf []byte, val uint64) []byte { return nil }

func (e errorString) RuntimeError() {}
func (e errorString) Error() string { return "" }
