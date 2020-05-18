package tune

import (
	"fmt"
	"time"

	"github.com/peterbourgon/g2s"
)

// NewStatsD takes a protocol, host, and prefix to initialize a statsD client with. It
// returns an error if initializing the client encounters any errors.
func NewStatsD(proto, host, prefix string) (*Client, error) {
	s, err := g2s.DialWithPrefix(proto, host, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize statter: %v", err.Error())
	}

	c := &Client{
		Statter: s,
	}

	return c, nil
}

// NewAsyncStatsD takes a protocol, hnost, prefix, and retry interval to initialize a
// statsD client with. It will retry on the interval given and returns the
// client on a channel.
func NewAsyncStatsD(proto, host, prefix string, retry time.Duration) <-chan *Client {
	out := make(chan *Client)

	if retry == 0 {
		retry = defaultRetry
	}

	go func() {
		for {
			s, err := NewStatsD(proto, host, prefix)
			if err != nil {
				<-time.After(retry)
				continue
			}

			out <- s
			break
		}
	}()

	return out
}
