package exenv

type env func(string) (string, bool)

func (e env) Lookup(key string) (string, error) {
	v, ok := e(key)
	if !ok {
		return "", &NotFoundError{Key: key}
	}
	return v, nil
}
