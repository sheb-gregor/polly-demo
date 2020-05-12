package client

import "encoding/json"

type PollyClient interface {
	PollyPublisher
	PollySubscriber
}

type PollyPublisher interface {
	Publish(topic string, data json.RawMessage) error
}

type PollySubscriber interface {
	Subscribe(topic, subscriber string) error
	Unsubscribe(topic, subscriber string) error
	Poll(topic, subscriber string) (json.RawMessage, error)
}
