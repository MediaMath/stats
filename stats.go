package stats

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/MediaMath/govent/graphite"
)

//Broker is a coordination point for stats and endpoints
type Broker chan interface{}

//StartBroker starts the background goroutine that listens for stats and forwards them
func StartBroker(bufferSize int) Broker {
	s := Broker(make(chan interface{}, bufferSize))
	go brokerLoop(s, bufferSize)
	return s
}

func brokerLoop(s Broker, bufferSize int) {
	endpoints := []endpoint{}
	var allDone sync.WaitGroup
	var pill poison

	for act := range s {
		endpoints, pill = doEvent(endpoints, &allDone, bufferSize, act)
		if pill != nil {
			break
		}
	}

	for _, endpoint := range endpoints {
		close(endpoint)
	}

	close(pill)
	allDone.Wait()
}

func doEvent(endpoints []endpoint, allDone *sync.WaitGroup, bufferSize int, act interface{}) ([]endpoint, poison) {
	switch a := act.(type) {
	case poison:
		return endpoints, a
	case Endpoint:
		e := make(chan interface{}, bufferSize)
		endpoints = append(endpoints, e)
		allDone.Add(1)
		go func() {
			a(e)
			allDone.Done()
		}()
	default:
		for _, e := range endpoints {
			select {
			case e <- a:
			default:
			}
		}
	}

	return endpoints, nil
}

//Endpoint is a function that takes a channel of stats and reacts to them. It will be started in a go routine by the broker
type Endpoint func(<-chan interface{})
type endpoint chan<- interface{}

type poison chan<- error

//ErrActivityBufferFull is returned if the brokers buffer is full when attempting to register an endpoint or stop the broker
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

//Finish will attempt to shutdown the broker and all endpoints after sending buffered stats
func (s Broker) Finish(ctx context.Context) error {
	done := make(chan error)
	select {
	case s <- poison(done):
	default:
		return ErrActivityBufferFull
	}

	var err error
	select {
	case err = <-done:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return err
}

//Send will send the supplied datum
func (s Broker) Send(datum interface{}) {
	select {
	case s <- datum:
	default:
	}
}

//Count sends a count value for the given name
func (s Broker) Count(name string, value int64) {
	var (
		tags       []string
		sampleRate float64 = 1
	)
	s.Send(&count{Name: name, Value: value, Tags: tags, Rate: sampleRate})
}

//Incr increments a count by 1
func (s Broker) Incr(name string) {
	s.Count(name, 1)
}

//Gauge sends a gauge value for the given name
func (s Broker) Gauge(name string, value float64) {
	var (
		tags       []string
		sampleRate float64 = 1
	)
	s.Send(&gauge{Name: name, Value: value, Tags: tags, Rate: sampleRate})
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
func (s Broker) Timing(name string, value time.Duration) {
	var (
		tags       []string
		sampleRate float64 = 1
	)
	s.Send(&timing{Name: name, Value: value, Tags: tags, Rate: sampleRate})
}

//TimingDuration sends a timing value for the duration provided
func (s Broker) TimingDuration(name string, duration time.Duration) {
	timeMillis := time.Duration(duration.Nanoseconds() / 1000000)
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

//DEvent will send a datadog event
func (s Broker) DEvent(tag string, data string) {
	s.Send(&devent{Title: tag, Text: data})
}
