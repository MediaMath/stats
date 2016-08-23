package stats

import (
	"fmt"
	"time"

	"github.com/MediaMath/govent/graphite"
)

//Broker is a coordination point for stats and endpoints
type Broker chan interface{}

//StartBroker starts the background goroutine that listens for stats and forwards them
func StartBroker(bufferSize int) Broker {
	s := Broker(make(chan interface{}, bufferSize))

	go func() {
		endpoints := []endpoint{}
		for act := range s {
			switch a := act.(type) {
			case Endpoint:
				e := make(chan interface{}, bufferSize)
				endpoints = append(endpoints, e)
				go a(e)
			default:
				for _, e := range endpoints {
					select {
					case e <- a:
					default:
					}
				}
			}

		}

		for _, endpoint := range endpoints {
			close(endpoint)
		}

	}()

	return s
}

//Endpoint is a function that takes a channel of stats and reacts to them. It will be started in a go routine by the broker
type Endpoint func(<-chan interface{})
type endpoint chan<- interface{}

//ErrActivityBufferFull is returned if the brokers buffer is full when attempting to register an endpoint
var ErrActivityBufferFull = fmt.Errorf("stats activity buffer full")

//RegisterEndpoint will add an endpoint to the list, the provided context will be listened to for cancellation
func (s Broker) RegisterEndpoint(e Endpoint) error {
	select {
	case s <- e:
	default:
		return ErrActivityBufferFull
	}

	return nil
}

//Send will send the supplied datum
func (s Broker) Send(datum interface{}) {
	select {
	case s <- datum:
	default:
	}
}

//Count sends a count value for the given name
func (s Broker) Count(name string, i int) {
	s.Send(&count{Name: name, Value: i})
}

//Incr increments a count by 1
func (s Broker) Incr(name string) {
	s.Count(name, 1)
}

//Gauge sends a gauge value for the given name
func (s Broker) Gauge(name string, i int) {
	s.Send(&gauge{Name: name, Value: i})
}

//On sends a 1 gauge
func (s Broker) On(name string) {
	s.Gauge(name, 1)
}

//Off sends a 0 gauge
func (s Broker) Off(name string) {
	s.Gauge(name, 0)
}

//Timing sends a timing value for the given name
func (s Broker) Timing(name string, i int) {
	s.Send(&timing{Name: name, Value: i})
}

//TimingDuration sends a timing value for the duration provided
func (s Broker) TimingDuration(name string, duration time.Duration) {
	timeMillis := int(duration.Nanoseconds() / 1000000)
	s.Timing(name, timeMillis)
}

//TimingPeriod sends a timing value for the given start and end
func (s Broker) TimingPeriod(name string, start time.Time, end time.Time) {
	s.TimingDuration(name, end.Sub(start))
}

//GraphiteEvent will send a graphite event
func (s Broker) GraphiteEvent(e *graphite.Event) {
	s.Send(&event{e})
}

//Event will send an event
func (s Broker) Event(tag string, data string) {
	s.GraphiteEvent(graphite.NewTaggedEvent(tag, data))
}
