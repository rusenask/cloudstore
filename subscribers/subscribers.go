package subscribers

import "github.com/lileio/pubsub"

// See https://godoc.org/github.com/lileio/pubsub#Handler for an example
// of what an subscriber handler should look like

type CloudstoreServiceSubscriber struct{}

func (s *CloudstoreServiceSubscriber) Setup(c *pubsub.Client) {
	// https://godoc.org/github.com/lileio/pubsub#Client.On
	// c.On(pubsub.HandlerOptions{
	// 	Topic:    "some_topic",
	// 	Name:     "service_name",
	// 	Handler:  s.SomeMethod,
	// 	Deadline: 10 * time.Second,
	// 	AutoAck:  true,
	// })
}
