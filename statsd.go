package stats

//Copyright 2016 MediaMath <http://www.mediamath.com>.  All rights reserved.
//Use of this source code is governed by a BSD-style
//license that can be found in the LICENSE file.

import "gopkg.in/alexcesaro/statsd.v2"

//StatsdEndpoint sends stats to a statsd client
func StatsdEndpoint(s *statsd.Client) Endpoint {
	return func(data <-chan interface{}) {
		for d := range data {
			switch t := d.(type) {
			case *count:
				s.Count(t.Name, t.Value)
			case *gauge:
				s.Gauge(t.Name, t.Value)
			case *biggauge:
				s.Gauge(t.Name, t.Value)
			case *timing:
				s.Timing(t.Name, t.Value)
			default:
			}
		}
	}
}
