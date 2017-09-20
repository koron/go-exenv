package exenv

// NotFoundError represents that key wan't found.
type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return "not found key: " + e.Key
}

// IsNotFound checks an error is NotFoundError. If it was, returns true and
// key.
func IsNotFound(err error) (bool, string) {
	e, ok := err.(*NotFoundError)
	if !ok {
		return false, ""
	}
	return ok, e.Key
}
