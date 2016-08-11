package stats

import (
	"log"

	"golang.org/x/net/context"

	"github.com/MediaMath/govent/graphite"
)

//RegisterGraphite registers a graphite event handler for any events that are sent
func RegisterGraphite(ctx context.Context, graphiteURL string, graphiteUser string, graphitePassword string, verbose bool) {
	log.Printf("Register graphite %v %v", graphiteUser, graphiteURL)
	g := &GraphiteEndpoint{graphiteURL, graphiteUser, graphitePassword, verbose}
	go g.Start(ctx, DefaultBroker)

	return
}

//GraphiteEndpoint is able to forward event style datum to a graphite http api
type GraphiteEndpoint struct {
	GraphiteURL      string
	GraphiteUser     string
	GraphitePassword string
	Verbose          bool
}

//Start registers the endpoint and begins looping over its data
func (e *GraphiteEndpoint) Start(ctx context.Context, broker Broker) {
	data := broker.RegisterEndpoint(ctx, 100)
	govent := graphite.NewVerbose(e.GraphiteUser, e.GraphitePassword, e.GraphiteURL, e.Verbose)

	for d := range data {
		err := d.SendToGraphiteEvents(govent)
		if err != nil {
			log.Printf("Error publishing %v to graphite: %v", d, err)
		}
	}
}
