# [stats](https://github.com/MediaMath/stats) &middot; [![CircleCI Status](https://circleci.com/gh/MediaMath/stats.svg?style=shield)](https://circleci.com/gh/MediaMath/stats) [![GitHub license](https://img.shields.io/badge/license-BSD3-blue.svg)](https://github.com/MediaMath/stats/blob/master/LICENSE) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/MediaMath/stats/blob/master/CONTRIBUTING.md)

## stats - a stats facade for golang 
[![GoDoc](https://godoc.org/github.com/MediaMath/stats?status.png)](https://godoc.org/github.com/MediaMath/stats)

Package stats provides a facade around push based stats instrumentation.  It allows you to put simple instrumentation calls in your code:

```go
	stats.Incr("get_request")
```

Then you can turn this off, change where the stats go, etc at startup time.

See the [GoDoc](https://godoc.org/github.com/MediaMath/stats) for more details.

