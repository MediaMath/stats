package stats

//Copyright 2016 MediaMath <http://www.mediamath.com>.  All rights reserved.
//Use of this source code is governed by a BSD-style
//license that can be found in the LICENSE file.

import (
	"log"

	"context"
)

//RegisterStatsLogger starts a stats receiver that logs the stats. This is useful for tests but would probably destroy a production system.
func RegisterStatsLogger(ctx context.Context) {
	DefaultBroker.RegisterEndpoint(LogEndpoint())
}

//LogEndpoint is an endpoint that logs data
func LogEndpoint() Endpoint {
	return func(data <-chan interface{}) {
		for d := range data {
			log.Printf("%v", d)
		}
	}
}
