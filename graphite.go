package stats

import (
	"log"

	"golang.org/x/net/context"

	"github.com/MediaMath/govent/graphite"
)

//StartGraphite registers the endpoint and begins looping over its data
func StartGraphite(ctx context.Context, broker Broker, govent *graphite.Graphite) {
	log.Printf("Starting graphite %v %v", govent.Username, govent.Addr)
	data := broker.RegisterEndpoint(ctx, 100)
	for d := range data {
		err := d.SendToGraphiteEvents(govent)
		if err != nil {
			log.Printf("Error publishing %v to graphite: %v", d, err)
		}
	}
}
