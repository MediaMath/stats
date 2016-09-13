package stats

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/MediaMath/govent/graphite"
)

func TestFunctionalGraphite(t *testing.T) {
	t.Skip("Unskip this to see if events are getting published")

	username := os.Getenv("GRAPHITE_USER")
	password := os.Getenv("GRAPHITE_PASSWORD")
	url := os.Getenv("GRAPHITE_URL")

	if testing.Short() {
		t.Skipf("skipped because is an integration test")
	}

	if username == "" || password == "" || url == "" {
		t.Skipf("skipped because missing creds")
	}

	govent := &graphite.Graphite{
		Username: username,
		Password: password,
		Addr:     url,
		Client:   &http.Client{Timeout: time.Second * 10},
		Verbose:  true,
		Prefix:   "com.mediamath.stats.tests",
	}

	broker := StartBroker(100)
	defer broker.Finish(context.Background())

	broker.RegisterEndpoint(GraphiteEndpoint(govent))
	broker.Event("tag.foo", "data.bar")
	<-time.After(5 * time.Second)
}
