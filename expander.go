package exenv

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	consul "github.com/hashicorp/consul/api"
)

type expander struct {
	l sync.Mutex
	p provider
	c *consul.Client
}

func (ex *expander) Lookup(key string) (string, error) {
	v, err := ex.p.Lookup(key)
	if err != nil {
		return "", err
	}
	switch ex.prefix(v) {
	case "consul:":
		s, err := ex.queryConsul(v[7:])
		if ok, k2 := IsNotFound(err); ok {
			return "", fmt.Errorf("consul: not found key %q for %s", k2, key)
		}
		if err != nil {
			return "", err
		}
		return s, nil
	case "raw:":
		return v[4:], nil
	default:
		return v, nil
	}
}

func (ex *expander) prefix(s string) string {
	n := strings.IndexRune(s, ':')
	if n < 0 {
		return ""
	}
	return strings.ToLower(s[:n+1])
}

func (ex *expander) queryConsul(key string) (string, error) {
	c, err := ex.client()
	if err != nil {
		return "", err
	}
	p, _, err := c.KV().Get(key, nil)
	if err != nil {
		return "", err
	} else if p == nil {
		return "", &NotFoundError{Key: key}
	}
	return string(p.Value), nil
}

func (ex *expander) client() (*consul.Client, error) {
	ex.l.Lock()
	defer ex.l.Unlock()
	if ex.c != nil {
		return ex.c, nil
	}
	// open consul client.
	s, err := ex.p.Lookup("_CONSUL_URL")
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	conf := consul.DefaultConfig()
	conf.Address = u.Host
	conf.Scheme = u.Scheme
	c, err := consul.NewClient(conf)
	if err != nil {
		return nil, err
	}
	ex.c = c
	return ex.c, nil
}
