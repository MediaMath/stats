package stats

//Copyright 2016 MediaMath <http://www.mediamath.com>.  All rights reserved.
//Use of this source code is governed by a BSD-style
//license that can be found in the LICENSE file.

import (
	"log"

	"github.com/MediaMath/govent/graphite"
)

//GraphiteEndpoint sends event data to a graphite endpoint
func GraphiteEndpoint(govent *graphite.Graphite) Endpoint {
	return func(data <-chan interface{}) {
		for d := range data {
			event, is := d.(event)
			if is {
				err := govent.Publish(event.inner)
				log.Printf("Error publishing %v to graphite: %v", d, err)
			}
		}
	}
}
