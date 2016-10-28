## stats - a stats facade for golang 
[![GoDoc](https://godoc.org/github.com/MediaMath/stats?status.png)](https://godoc.org/github.com/MediaMath/stats)

Package stats provides a facade around push based stats instrumentation.  It allows you to put simple instrumentation calls in your code:

```go
	stats.Incr("get_request")
```

Then you can turn this off, change where the stats go, etc at startup time.

See the [GoDoc](https://godoc.org/github.com/MediaMath/stats) for more details.

