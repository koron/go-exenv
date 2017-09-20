package exenv

type provider interface {
	Lookup(string) (string, error)
}
