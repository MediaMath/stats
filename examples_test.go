package stats

import (
	"context"
	"log"
)

func Example() {
	//stats will be prefixed with com.example.foo
	ctx := SetPrefix(context.Background(), "com.example.foo")

	//statsd location
	ctx = SetStatsdURL(ctx, "stats.example.com:57475")

	//data dog location
	ctx = SetDatadogURL(ctx, "https://app.datadoghq.com/api/")

	// will start forwarding stats to data dog and statsd (as they are in the context) but not graphite
	err := RegisterStatsContext(ctx)
	if err != nil {
		log.Printf("Error registering stats context you won't get any stats: %v", err)
	}

	// this will start a goroutine that forwards standard golang runtime stats (like gc/thread counts etc)
	err = RegisterRuntimeStatsContext(ctx)
	if err != nil {
		log.Printf("Error registering runtime stats context you won't get runtime stats: %v", err)
	}

	// will send an event named com.example.foo.foo with the bar.baz data.  Will not go to graphite as it isn't set up in the context
	// will not go to statsd as it doesn't handle event stats, will go to data dog as it does
	Event("foo", "bar.baz")

	//Sends com.example.foo.boom|7|c
	Count("boom", 7)

	//Sends com.example.bam|1|c
	Incr("bam")
}
