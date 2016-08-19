package stats

import (
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestDatum(t *testing.T) {
	ctx, clean := context.WithTimeout(context.Background(), time.Second)
	defer clean()

	broker := NewBroker(100)
	go broker.Start(ctx)
	go LogStats(ctx, broker)

	data := broker.RegisterEndpoint(ctx, 100)

	type counts struct {
		count  int
		gauge  int
		timing int
	}

	sum := make(chan *counts)

	go func() {
		received := 0
		c := &counts{}
		for datum := range data {
			switch t := datum.(type) {
			case *count:
				c.count += t.Value
			case *gauge:
				c.gauge += t.Value
			case *timing:
				c.timing += t.Value
			}

			received++
			if received >= 6 {
				break
			}
		}

		sum <- c
		close(sum)
	}()

	broker.Count("C", 1)
	broker.Gauge("G", 3)
	broker.Timing("T", 3)
	broker.Count("C", 8)
	broker.Timing("T", 1)
	broker.Gauge("G", 3)

	select {
	case c := <-sum:
		if c == nil {
			t.Fatal("No c")
		}

		if c.count != 9 || c.gauge != 6 || c.timing != 4 {
			t.Errorf("Wrong: %v", c)
		}
	case <-ctx.Done():
		t.Fatal("Timed out")
	}

}

func TestMultipleRegistered(t *testing.T) {
	ctx, clean := context.WithTimeout(context.Background(), 1*time.Second)
	defer clean()

	broker := NewBroker(100)
	go broker.Start(ctx)

	data1 := broker.RegisterEndpoint(ctx, 100)
	data2 := broker.RegisterEndpoint(ctx, 100)

	broker.Count("C", 7)
	broker.Count("C", 8)
	broker.Count("C", 6)

	sum1 := make(chan int)
	go func() {
		s := 0
		for c := range data1 {
			t := c.(*count)
			s += t.Value

			if s >= 21 {
				break
			}
		}
		sum1 <- s
		close(sum1)
	}()

	sum2 := make(chan int)
	go func() {
		s := 0
		for c := range data2 {
			t := c.(*count)
			s += t.Value

			if s >= 21 {
				break
			}
		}
		sum2 <- s
		close(sum2)
	}()

	select {
	case s := <-sum1:
		if s != 21 {
			t.Errorf("Wrong: %v", s)
		}
	case <-ctx.Done():
		t.Fatal("Timed out")
	}

	select {
	case s := <-sum2:
		if s != 21 {
			t.Errorf("Wrong: %v", s)
		}
	case <-ctx.Done():
		t.Fatal("Timed out")
	}
}
