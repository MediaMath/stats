package stats

import (
	"fmt"

	"github.com/MediaMath/govent/graphite"
)

type event struct {
	inner *graphite.Event
}

func (e event) String() string {
	return fmt.Sprintf("%v|%v|e", e.inner.Tags, e.inner.Data)
}

type count struct {
	Name  string
	Value int
	Tags  []string
	Rate  float64
}

func (c *count) String() string {
	return fmt.Sprintf("%v|%v|c", c.Name, c.Value)
}

type gauge struct {
	Name  string
	Value int
	Tags  []string
	Rate  float64
}

func (g *gauge) String() string {
	return fmt.Sprintf("%v|%v|g", g.Name, g.Value)
}

type timing struct {
	Name  string
	Value int
	Tags  []string
	Rate  float64
}

func (t *timing) String() string {
	return fmt.Sprintf("%v|%v|t", t.Name, t.Value)
}
