package stats

import (
	"fmt"

	"github.com/MediaMath/govent/graphite"

	"gopkg.in/alexcesaro/statsd.v2"
)

type sendableDatum struct {
	Datum
}

func (d *sendableDatum) Action(endpoints map[endpoint]bool) (map[endpoint]bool, error) {
	for e := range endpoints {
		select {
		case e.data <- d.Datum:
		default:
		}
	}

	return endpoints, nil
}

//Datum is a individual statistic
type Datum interface {
	SendToStatsd(s *statsd.Client)
	SendToGraphiteEvents(g *graphite.Graphite) error
}

type event struct {
	inner *graphite.Event
}

func (e event) String() string {
	return fmt.Sprintf("%v|%v|e", e.inner.Tags, e.inner.Data)
}

func (e event) SendToStatsd(s *statsd.Client) {
}

func (e event) SendToGraphiteEvents(g *graphite.Graphite) error {
	return g.Publish(e.inner)
}

type count struct {
	Name  string
	Value int
}

func (c *count) String() string {
	return fmt.Sprintf("%v|%v|c", c.Name, c.Value)
}

func (c *count) SendToStatsd(s *statsd.Client) {
	s.Count(c.Name, c.Value)
}

func (c *count) SendToGraphiteEvents(g *graphite.Graphite) error {
	return nil
}

type gauge struct {
	Name  string
	Value int
}

func (g *gauge) String() string {
	return fmt.Sprintf("%v|%v|g", g.Name, g.Value)
}

func (g *gauge) SendToStatsd(s *statsd.Client) {
	s.Gauge(g.Name, g.Value)
}

func (g *gauge) SendToGraphiteEvents(gr *graphite.Graphite) error {
	return nil
}

type timing struct {
	Name  string
	Value int
}

func (t *timing) String() string {
	return fmt.Sprintf("%v|%v|t", t.Name, t.Value)
}

func (t *timing) SendToStatsd(s *statsd.Client) {
	s.Timing(t.Name, t.Value)
}

func (t *timing) SendToGraphiteEvents(g *graphite.Graphite) error {
	return nil
}
