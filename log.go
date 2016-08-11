package stats

import (
	"log"

	"golang.org/x/net/context"
)

//RegisterStatsLogger starts a stats receiver that logs the stats. This is useful for tests but would probably destroy a production system.
func RegisterStatsLogger(ctx context.Context) {
	go LogStats(ctx, DefaultBroker)
}

//LogStats registers the enddpoint and begins logging the received data
func LogStats(ctx context.Context, broker Broker) {
	data := broker.RegisterEndpoint(ctx, 100)
	for d := range data {
		log.Printf("%v", d)
	}
}
