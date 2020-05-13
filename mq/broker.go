package mq

import (
	"encoding/json"
	"sync"
)

type Broker struct {
	topics sync.Map // topics is a map[string]Topic
}

// NewBroker creates new instance of Message Broker.
func NewBroker() *Broker {
	return &Broker{topics: sync.Map{}}
}

// HandleNewMessage puts the message to the topic if it exist.
func (broker *Broker) HandleNewMessage(topic string, data json.RawMessage) {
	raw, present := broker.topics.Load(topic)
	tReg, ok := raw.(*Topic)
	if !present || !ok {
		return
	}

	tReg.PutMessage(data)
	broker.topics.Store(topic, tReg)
}

// Subscribe adds the subscriber to the provided topic.
// The topic will be created if it does not already exist.
func (broker *Broker) Subscribe(topic, subscriber string) {
	raw, present := broker.topics.Load(topic)
	tReg, ok := raw.(*Topic)
	if !present || !ok {
		tReg = NewTopic()
	}

	tReg.Subscribe(subscriber)
	broker.topics.Store(topic, tReg)
}

// Unsubscribe removes the subscriber from the provided topic and
// deletes the topic if there are no subscribers.
func (broker *Broker) Unsubscribe(topic, subscriber string) {
	raw, present := broker.topics.Load(topic)
	tReg, ok := raw.(*Topic)
	if !present || !ok {
		return
	}

	tReg.Unsubscribe(subscriber)
	if tReg.subCount <= 0 {
		broker.topics.Delete(topic)
		return
	}

	broker.topics.Store(topic, tReg)
}

// Poll fetch the next unseen message or no message if everything is seen,
// or `false` if the subscription is not found.
func (broker *Broker) Poll(topic, subscriber string) (json.RawMessage, bool) {
	raw, present := broker.topics.Load(topic)
	tReg, ok := raw.(*Topic)
	if !present || !ok {
		return nil, false
	}

	msg, subscribed := tReg.Poll(subscriber)
	broker.topics.Store(topic, tReg)
	return msg, subscribed
}
