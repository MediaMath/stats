package stats

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

type key int

const (
	prefixKey key = iota
	statsdURLKey
	runtimeIntervalKey
	graphiteURLKey
	graphiteUserKey
	graphitePasswordKey
	graphiteVerboseKey
)

//SetPrefix sets the stats prefix
func SetPrefix(ctx context.Context, prefix string) context.Context {
	return context.WithValue(ctx, prefixKey, prefix)
}

//GetPrefix gets the prefix
func GetPrefix(ctx context.Context) string {
	return getString(ctx, prefixKey, "")
}

//SetStatsdURL sets the stats prefix
func SetStatsdURL(ctx context.Context, url string) context.Context {
	return context.WithValue(ctx, statsdURLKey, url)
}

//SetGraphite sets the graphite client
func SetGraphite(ctx context.Context, url, user, password string, verbose bool) context.Context {
	ctx = context.WithValue(ctx, graphiteURLKey, url)
	ctx = context.WithValue(ctx, graphiteUserKey, user)
	ctx = context.WithValue(ctx, graphitePasswordKey, password)
	ctx = context.WithValue(ctx, graphiteVerboseKey, verbose)

	return ctx
}

//SetRuntimeInterval sets the runtime stats collector interval
func SetRuntimeInterval(ctx context.Context, interval time.Duration) context.Context {
	return context.WithValue(ctx, runtimeIntervalKey, interval)
}

//HasStats checks if the statsd url and graphite url are set
func HasStats(ctx context.Context) (hasStatsdURL bool, hasGraphiteURL bool) {
	statsdURL := getString(ctx, statsdURLKey, "")
	graphiteURL := getString(ctx, graphiteURLKey, "")

	return statsdURL != "", graphiteURL != ""
}

//RegisterStatsContext starts statsd and graphite based on the context
func RegisterStatsContext(ctx context.Context) error {
	prefix := GetPrefix(ctx)
	if prefix == "" {
		return fmt.Errorf("No prefix not starting stats consumers")
	}

	statsdURL := getString(ctx, statsdURLKey, "")
	if statsdURL == "" {
		return fmt.Errorf("No statsd URL not starting stats consumers")
	}

	graphiteURL := getString(ctx, graphiteURLKey, "")
	if graphiteURL == "" {
		return fmt.Errorf("No graphite URL not starting stats consumers")
	}

	RegisterStatsd(ctx, statsdURL, prefix)

	graphiteUser := getString(ctx, graphiteUserKey, "")
	graphitePassword := getString(ctx, graphitePasswordKey, "")
	graphiteVerbose, _ := ctx.Value(graphiteVerboseKey).(bool)

	RegisterGraphite(ctx, graphiteURL, graphiteUser, graphitePassword, graphiteVerbose)
	return nil
}

//RegisterRuntimeStatsContext starts runtime stats reporting based on the context
func RegisterRuntimeStatsContext(ctx context.Context) error {
	interval, has := ctx.Value(runtimeIntervalKey).(time.Duration)
	if !has {
		return fmt.Errorf("No runtime interval not reporting runtime stats")
	}

	return ReportRuntimeStats(interval, ctx.Done())
}

func getString(ctx context.Context, key key, def string) string {
	val, has := ctx.Value(key).(string)
	if !has {
		return def
	}

	return val
}
