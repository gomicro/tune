package tune

import (
	"github.com/peterbourgon/g2s"
)

// NewNoop returns a new client that will perform noops.
func NewNoop() *Client {
	c := &Client{
		Statter: g2s.Noop(),
	}

	return c
}
