package stats

import (
	"time"

	"golang.org/x/net/context"
)

//DDBroker is a datadog statsd broker that runs in backgorund
var DDBroker Broker

func init() {
	DDBroker = StartBroker(100)
}

//DDFinish will finish the DDBroker
func DDFinish(ctx context.Context) error {
	return DDBroker.Finish(ctx)
}

//DDSend will send the supplied datum
func DDSend(datum interface{}) {
	DDBroker.Send(datum)
}

//DDEvent will send an event
func DDEvent(tag string, data string) {
	DDBroker.DDEvent(tag, data)
}

//DDIncr increments a count by 1
func DDIncr(name string) {
	DDBroker.DDIncr(name)
}

//DDCount sends a count value for the given name
func DDCount(name string, i int64) {
	DDBroker.DDCount(name, i)
}

//DDOn sends a 1 gauge
func DDOn(name string) {
	DDBroker.DDOn(name)
}

//DDOff sends a 0 gauge
func DDOff(name string) {
	DDBroker.DDOff(name)
}

//DDGauge sends a gauge value for the given name
func DDGauge(name string, i float64) {
	DDBroker.DDGauge(name, i)
}

//DDTiming sends a timing value for the given name
func DDTiming(name string, i time.Duration) {
	DDBroker.DDTiming(name, i)
}

//DDTimingDuration sends a timing value for the duration provided
func DDTimingDuration(name string, duration time.Duration) {
	DDBroker.DDTimingDuration(name, duration)
}

//DDTimingPeriod sends a timing value for the given start and end
func DDTimingPeriod(name string, start time.Time, end time.Time) {
	DDBroker.DDTimingPeriod(name, start, end)
}
