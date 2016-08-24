package stats

import ddStatsd "github.com/DataDog/datadog-go/statsd"

//DatadogStatsdEndpoint sends stats to a statsd client
func DatadogStatsdEndpoint(s *ddStatsd.Client) Endpoint {
	return func(data <-chan interface{}) {
		for d := range data {
			switch t := d.(type) {
			case *ddcount:
				s.Count(t.Name, t.Value, t.Tags, t.Rate)
			case *ddgauge:
				s.Gauge(t.Name, t.Value, t.Tags, t.Rate)
			case *ddtiming:
				s.Timing(t.Name, t.Value, t.Tags, t.Rate)
			case *ddevent:
				s.SimpleEvent(t.Title, t.Text)
			default:
			}
		}
	}
}
