package stats

import (
	"log"

	"golang.org/x/net/context"
	"gopkg.in/alexcesaro/statsd.v2"
)

//RegisterStatsd adds a listener that forwards data to a statsd address
func RegisterStatsd(ctx context.Context, url string, prefix string) {
	s, err := statsd.New(statsd.Address(url), statsd.Prefix(prefix))
	if err != nil {
		log.Printf("Error with statsd: %v", err)
	} else {
		endpoint := &StatsdEndpoint{s}
		go endpoint.Start(ctx, DefaultBroker)
	}

}

//StatsdEndpoint receives stats and forwards them to a statsd receiver
type StatsdEndpoint struct {
	*statsd.Client
}

//Start will register the endpoint and loop through the received data
func (s *StatsdEndpoint) Start(ctx context.Context, broker Broker) {
	data := broker.RegisterEndpoint(ctx, 100)
	for d := range data {
		d.SendToStatsd(s.Client)
	}
}
