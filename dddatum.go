package stats

import (
	"fmt"
	"time"
)

type ddevent struct {
	Title string
	Text  string
}

func (e ddevent) String() string {
	return fmt.Sprintf("%s|%s|e", e.Title, e.Text)
}

type ddcount struct {
	Name  string
	Value int64
	Tags  []string
	Rate  float64
}

func (c ddcount) String() string {
	return fmt.Sprintf("%v|%v|c", c.Name, c.Value)
}

type ddgauge struct {
	Name  string
	Value float64
	Tags  []string
	Rate  float64
}

func (g ddgauge) String() string {
	return fmt.Sprintf("%v|%v|g", g.Name, g.Value)
}

type ddtiming struct {
	Name  string
	Value time.Duration
	Tags  []string
	Rate  float64
}

func (t ddtiming) String() string {
	return fmt.Sprintf("%v|%v|t", t.Name, t.Value)
}
