package stats

import (
	"golang.org/x/net/context"
	"gopkg.in/alexcesaro/statsd.v2"
)

//StartStatsd registers the endpoint and begins looping over its data
func StartStatsd(ctx context.Context, broker Broker, s *statsd.Client) {
	data := broker.RegisterEndpoint(ctx, 100)
	for d := range data {
		switch t := d.(type) {
		case *count:
			s.Count(t.Name, t.Value)
		case *gauge:
			s.Gauge(t.Name, t.Value)
		case *timing:
			s.Timing(t.Name, t.Value)
		default:
		}
	}
}
