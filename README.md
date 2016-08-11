## stats - a stats library for golang

This library allows you to instrument your code so that it can send stats (or not)

```go

func main() {
	//will setup the default stats calls to send to statsd with the prefix 'com.example.stats' ignore events
	stats.RegisterStatsd(context.Background(), "stats.example.com:57475", "com.example.stats")

	//ignored
	stats.Event("foo", "bar.baz")

	//sends com.example.stats.boom count 7
	stats.Count("boom", 7)
	//increments com.example.stats.bam
	stats.Incr("bam")
}

```
