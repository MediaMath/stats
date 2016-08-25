package stats

import (
	"time"

	"golang.org/x/net/context"

	"github.com/MediaMath/govent/graphite"
)

//DefaultBroker is started by default and runs in the background
var DefaultBroker Broker

func init() {
	DefaultBroker = StartBroker(100)
}

//Finish will finish the DefaultBroker
func Finish(ctx context.Context) error {
	return DefaultBroker.Finish(ctx)
}

//Send will send the supplied datum
func Send(datum interface{}) {
	DefaultBroker.Send(datum)
}

//Event will send an event
func Event(tag string, data string) {
	DefaultBroker.Event(tag, data)
}

//DEvent will send a datadog event
func DEvent(tag string, data string) {
	DefaultBroker.DEvent(tag, data)
}

//GraphiteEvent will send a graphite event
func GraphiteEvent(e *graphite.Event) {
	DefaultBroker.GraphiteEvent(e)
}

//Incr increments a count by 1
func Incr(name string) {
	DefaultBroker.Incr(name)
}

//Count sends a count value for the given name
func Count(name string, i int64) {
	DefaultBroker.Count(name, i)
}

//On sends a 1 gauge
func On(name string) {
	DefaultBroker.On(name)
}

//Off sends a 0 gauge
func Off(name string) {
	DefaultBroker.Off(name)
}

//Gauge sends a gauge value for the given name
func Gauge(name string, i float64) {
	DefaultBroker.Gauge(name, i)
}

//Timing sends a timing value for the given name
func Timing(name string, i time.Duration) {
	DefaultBroker.Timing(name, i)
}

//TimingDuration sends a timing value for the duration provided
func TimingDuration(name string, duration time.Duration) {
	DefaultBroker.TimingDuration(name, duration)
}

//TimingPeriod sends a timing value for the given start and end
func TimingPeriod(name string, start time.Time, end time.Time) {
	DefaultBroker.TimingPeriod(name, start, end)
}
