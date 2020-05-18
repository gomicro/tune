package tune

import (
	"fmt"
	"net/http"
	"time"

	"github.com/peterbourgon/g2s"
)

// Client represents a client wrapper for a statsD client
type Client struct {
	g2s.Statter
}

var (
	defaultRetry = 1 * time.Second
)

// StatEndpoint wraps a http handler to collect and send stats to the
// aggregator. It will send counts and timing metrics to be aggregated.
func (c *Client) StatEndpoint(fn http.HandlerFunc, label string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		n := time.Now()

		fn(w, r)

		c.Timing(1.0, fmt.Sprintf("%v.time", label), time.Since(n)*time.Millisecond)
		c.Counter(1.0, fmt.Sprintf("%v.count", label), 1)
	}
}
