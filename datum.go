package stats

//Copyright 2016 MediaMath <http://www.mediamath.com>.  All rights reserved.
//Use of this source code is governed by a BSD-style
//license that can be found in the LICENSE file.

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
}

func (c *count) String() string {
	return fmt.Sprintf("%v|%v|c", c.Name, c.Value)
}

type gauge struct {
	Name  string
	Value int
}

func (g *gauge) String() string {
	return fmt.Sprintf("%v|%v|g", g.Name, g.Value)
}

type biggauge struct {
	Name  string
	Value uint64
}

func (g *biggauge) String() string {
	return fmt.Sprintf("%v|%v|g", g.Name, g.Value)
}

type timing struct {
	Name  string
	Value int
}

func (t *timing) String() string {
	return fmt.Sprintf("%v|%v|t", t.Name, t.Value)
}
