package stats

import (
	"time"

	ddStatsd "github.com/DataDog/datadog-go/statsd"
)

//DatadogStatsdEndpoint sends stats to a statsd client
func DatadogStatsdEndpoint(s *ddStatsd.Client) Endpoint {
	return func(data <-chan interface{}) {
		for d := range data {
			switch t := d.(type) {
			case *count:
				s.Count(t.Name, int64(t.Value), nil, 1)
			case *gauge:
				s.Gauge(t.Name, float64(t.Value), nil, 1)
			case *timing:
				s.Timing(t.Name, time.Duration(t.Value), nil, 1)
			case *event:
				s.SimpleEvent(t.inner.Tags, t.inner.Data)
			default:
			}
		}
	}
}
