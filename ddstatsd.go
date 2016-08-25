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
				s.Count(t.Name, int64(t.Value), t.Tags, t.Rate)
			case *gauge:
				s.Gauge(t.Name, float64(t.Value), t.Tags, t.Rate)
			case *timing:
				s.Timing(t.Name, time.Duration(t.Value), t.Tags, t.Rate)
			default:
			}
		}
	}
}
