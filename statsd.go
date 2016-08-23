package stats

import "gopkg.in/alexcesaro/statsd.v2"

//StatsdEndpoint sends stats to a statsd client
func StatsdEndpoint(s *statsd.Client) Endpoint {
	return func(data <-chan interface{}) {
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
}
