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

	nsInMs := uint64(time.Millisecond)

	for {
		select {
		case <-time.After(sleep):
			runtime.ReadMemStats(memStats)

			now := time.Now()

			Gauge("runtime.NumGoroutine", float64(runtime.NumGoroutine()))

			Gauge("memory.Alloc", float64(memStats.Alloc))
			Gauge("memory.Mallocs", float64(memStats.Mallocs))
			Gauge("memory.Frees", float64(memStats.Frees))

			Gauge("memory.HeapSys", float64(memStats.HeapSys))
			Gauge("memory.HeapAlloc", float64(memStats.HeapAlloc))
			Gauge("memory.HeapIdle", float64(memStats.HeapIdle))
			Gauge("memory.HeapInUse", float64(memStats.HeapInuse))
			Gauge("memory.HeapReleased", float64(memStats.HeapReleased))
			Gauge("memory.HeapObject", float64(memStats.HeapObjects))

			Gauge("memory.StackInUse", float64(memStats.StackInuse))
			Gauge("memory.StackSys", float64(memStats.StackSys))
			Gauge("memory.MSpanSys", float64(memStats.MSpanSys))
			Gauge("memory.MSpanInUse", float64(memStats.MSpanInuse))
			Gauge("memory.MCacheInUse", float64(memStats.MCacheInuse))
			Gauge("memory.MCacheSys", float64(memStats.MCacheSys))

			Gauge("memory.NextGC", float64(memStats.NextGC))
			Gauge("memory.LastGC", float64(memStats.LastGC))
			Gauge("memory.gc.total_pauses", float64((memStats.PauseTotalNs)/nsInMs))

			if lastPauseNs > 0 {
				pauseSinceLastSample := memStats.PauseTotalNs - lastPauseNs
				Gauge("memory.gc.pauses_per_second", float64((pauseSinceLastSample)/nsInMs/uint64(sleep.Seconds())))
			}
			lastPauseNs = memStats.PauseTotalNs

			countGc := float64(memStats.NumGC - lastNumGc)
			if lastNumGc > 0 {
				diff := countGc
				diffTime := now.Sub(lastSampleTime).Seconds()
				Gauge("memory.gc.gcs_per_second", diff/float64(diffTime))
			}

			if countGc > 0 {
				if countGc > 256 {
					//We're missing some gc pause times reset to 256
					countGc = 256
				}

				for i := float64(0); i < countGc; i++ {
					idx := ((memStats.NumGC - uint32(i)) + 255) % 256
					pause := memStats.PauseNs[idx]
					Timing("memory.gc.pause", time.Duration(pause/nsInMs))
				}
			}

			lastNumGc = memStats.NumGC
			lastSampleTime = now
		case <-ctx.Done():
			return
		}
	}
}
