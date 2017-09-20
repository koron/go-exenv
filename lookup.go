package exenv

// Lookup lookups a value of an environment variable and expand the value by
// its contents and external providers.  Currently support two expansions:
// "consul:some_key" and "raw:raw_text".
func Lookup(key string) (string, error) {
	return defaultProvider.Lookup(key)
}
