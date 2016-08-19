package stats

import (
	"log"

	"golang.org/x/net/context"

	"github.com/MediaMath/govent/graphite"
)

//StartGraphite registers the endpoint and begins looping over its data
func StartGraphite(ctx context.Context, broker Broker, govent *graphite.Graphite) {
	data := broker.RegisterEndpoint(ctx, 100)
	for d := range data {
		event, is := d.(event)
		if is {
			err := govent.Publish(event.inner)
			log.Printf("Error publishing %v to graphite: %v", d, err)
		}
	}
}
