package stats

import (
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"
)

var running int32

//ReportRuntimeStats publishes runtime stats to statsd
func ReportRuntimeStats(ctx context.Context, sleep time.Duration) error {
	if !atomic.CompareAndSwapInt32(&running, 0, 1) {
		return fmt.Errorf("runtime stats are already running")
	}

	log.Printf("Updating runtime stats every: %v", sleep)
	go sendMemStats(ctx, sleep)
	return nil
}

func sendMemStats(ctx context.Context, sleep time.Duration) {
	var lastPauseNs uint64
	var lastNumGc uint32

	memStats := &runtime.MemStats{}
	lastSampleTime := time.Now()

	nsInMs := int(time.Millisecond)

	for {
		select {
		case <-time.After(sleep):
			runtime.ReadMemStats(memStats)

			now := time.Now()

			Gauge("runtime.NumGoroutine", runtime.NumGoroutine())

			Gauge("memory.Alloc", int(memStats.Alloc))
			Gauge("memory.Mallocs", int(memStats.Mallocs))
			Gauge("memory.Frees", int(memStats.Frees))

			Gauge("memory.HeapSys", int(memStats.HeapSys))
			Gauge("memory.HeapAlloc", int(memStats.HeapAlloc))
			Gauge("memory.HeapIdle", int(memStats.HeapIdle))
			Gauge("memory.HeapInUse", int(memStats.HeapInuse))
			Gauge("memory.HeapReleased", int(memStats.HeapReleased))
			Gauge("memory.HeapObject", int(memStats.HeapObjects))

			Gauge("memory.StackInUse", int(memStats.StackInuse))
			Gauge("memory.StackSys", int(memStats.StackSys))
			Gauge("memory.MSpanSys", int(memStats.MSpanSys))
			Gauge("memory.MSpanInUse", int(memStats.MSpanInuse))
			Gauge("memory.MCacheInUse", int(memStats.MCacheInuse))
			Gauge("memory.MCacheSys", int(memStats.MCacheSys))

			Gauge("memory.NextGC", int(memStats.NextGC))
			Gauge("memory.LastGC", int(memStats.LastGC))
			Gauge("memory.gc.total_pauses", int(memStats.PauseTotalNs)/nsInMs)

			if lastPauseNs > 0 {
				pauseSinceLastSample := memStats.PauseTotalNs - lastPauseNs
				Gauge("memory.gc.pauses_per_second", int(pauseSinceLastSample)/nsInMs/int(sleep.Seconds()))
			}
			lastPauseNs = memStats.PauseTotalNs

			countGc := int(memStats.NumGC - lastNumGc)
			if lastNumGc > 0 {
				diff := countGc
				diffTime := now.Sub(lastSampleTime).Seconds()
				Gauge("memory.gc.gcs_per_second", diff/int(diffTime))
			}

			if countGc > 0 {
				if countGc > 256 {
					//We're missing some gc pause times reset to 256
					countGc = 256
				}

				for i := 0; i < countGc; i++ {
					idx := int((memStats.NumGC-uint32(i))+255) % 256
					pause := int(memStats.PauseNs[idx])
					Timing("memory.gc.pause", pause/nsInMs)
				}
			}

			lastNumGc = memStats.NumGC
			lastSampleTime = now
		case <-ctx.Done():
			return
		}
	}
}
