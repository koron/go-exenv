package exenv

import (
	"os"
	"sync"
)

type item struct {
	v   string
	err error
}

type cache struct {
	l sync.Mutex
	m map[string]item
	p provider
}

func (c *cache) Lookup(key string) (string, error) {
	c.l.Lock()
	defer c.l.Unlock()
	if n, ok := c.m[key]; ok {
		return n.v, n.err
	}
	v, err := c.p.Lookup(key)
	c.m[key] = item{v: v, err: err}
	return v, err
}

var defaultProvider = &cache{
	m: map[string]item{},
	p: &expander{
		p: env(os.LookupEnv),
	},
}
