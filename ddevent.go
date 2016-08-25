package stats

import (
	"log"

	ddStatsd "github.com/DataDog/datadog-go/statsd"
)

//DDEventEndpoint sends event data to a datadog endpoint
func DDEventEndpoint(client *ddStatsd.Client) Endpoint {
	return func(data <-chan interface{}) {
		for d := range data {
			event, is := d.(event)
			if is {
				err := client.SimpleEvent(event.inner.Tags, event.inner.Data)
				log.Printf("Error publishing event %v to datadog: %v", d, err)
			}
		}
	}
}
