package stats

import (
	"fmt"
	"time"

	"github.com/MediaMath/govent/graphite"

	"golang.org/x/net/context"
)

//Broker is a coordination point for stats and endpoints
type Broker chan interface{}

//NewBroker sets up a broker with the provided buffer size
func NewBroker(bufferSize int) Broker {
	return Broker(make(chan interface{}, 100))
}

//RegisterEndpoint will add an endpoint to the list, the provided context will be listened to for cancellation
func (s Broker) RegisterEndpoint(ctx context.Context, bufferSize int) <-chan interface{} {
	data := make(chan interface{}, bufferSize)

	select {
	case s <- &endpoint{data, ctx}:
	case <-ctx.Done():
	}

	return data
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

//Start starts the background goroutine that listens for stats and forwards them
func (s Broker) Start(ctx context.Context) {
	go func() {
		endpoints := []*endpoint{}
		for act := range s {
			if isDone(ctx) {
				break
			}

			endpoints = cleanupEndpoints(endpoints)

			switch a := act.(type) {
			case *endpoint:
				endpoints = append(endpoints, a)
			default:
				for _, e := range endpoints {
					e.send(a)
				}
			}

		}

		for _, endpoint := range endpoints {
			close(endpoint.data)
		}
	}()
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func cleanupEndpoints(endpoints []*endpoint) []*endpoint {
	cleaned := []*endpoint{}
	for _, endpoint := range endpoints {
		select {
		case <-endpoint.ctx.Done():
			close(endpoint.data)
		default:
			cleaned = append(cleaned, endpoint)
		}
	}

	return cleaned
}

type endpoint struct {
	data chan<- interface{}
	ctx  context.Context
}

func (e *endpoint) send(d interface{}) {
	select {
	case e.data <- d:
	default:
	}
}

//ErrActivityBufferFull is returned anytime the activity buffer is full in stats
var ErrActivityBufferFull = fmt.Errorf("stats activity buffer full")
