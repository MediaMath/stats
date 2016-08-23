package stats

import (
	"log"

	"github.com/MediaMath/govent/graphite"
)

//GraphiteEndpoint sends event data to a graphite endpoint
func GraphiteEndpoint(govent *graphite.Graphite) Endpoint {
	return func(data <-chan interface{}) {
		for d := range data {
			event, is := d.(event)
			if is {
				err := govent.Publish(event.inner)
				log.Printf("Error publishing %v to graphite: %v", d, err)
			}
		}
	}
}
