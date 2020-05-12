package server

import (
	"container/list"
	"encoding/json"
	"sync"
)

type Broker struct {
	topics sync.Map // topics is a map[string]Topic
}

func NewBroker() *Broker {
	return &Broker{
		topics: sync.Map{},
	}
}

func (broker *Broker) HandleNewMessage(topic string, data json.RawMessage) {
	raw, present := broker.topics.Load(topic)
	tReg, ok := raw.(*Topic)
	if !present || !ok {
		return
	}

	tReg.PutMessage(data)
	broker.topics.Store(topic, tReg)
	return
}

func (broker *Broker) Subscribe(topic, subscriber string) {
	raw, present := broker.topics.Load(topic)
	tReg, ok := raw.(*Topic)
	if !present || !ok {
		tReg = NewTopic()
	}

	tReg.Subscribe(subscriber)
	broker.topics.Store(topic, tReg)
}

func (broker *Broker) Unsubscribe(topic, subscriber string) {
	raw, present := broker.topics.Load(topic)
	tReg, ok := raw.(*Topic)
	if !present || !ok {
		return
	}

	tReg.Unsubscribe(subscriber)
	broker.topics.Store(topic, tReg)
}

func (broker *Broker) Poll(topic, subscriber string) json.RawMessage {
	raw, present := broker.topics.Load(topic)
	tReg, ok := raw.(*Topic)
	if !present || !ok {
		return nil
	}

	msg := tReg.Poll(subscriber)
	broker.topics.Store(topic, tReg)
	return msg
}

type Topic struct {
	sync.Mutex

	subCount int64
	lastId   int64

	// subscribers - this map contains Unread Message Queues (FIFOs) for each user,
	// the value of the list item is the message identifier.
	subscribers map[string]*list.List
	// unreadCount map contains counters that show how many subscribers have not yet received each message.
	unreadCount map[int64]int64
	// messages map with messages, key is unique ID of message
	messages map[int64]json.RawMessage
}

func NewTopic() *Topic {
	return &Topic{
		subscribers: map[string]*list.List{},
		unreadCount: map[int64]int64{},
		messages:    map[int64]json.RawMessage{},
	}
}

func (topic *Topic) Poll(subscriber string) json.RawMessage {
	topic.Lock()
	defer topic.Unlock()

	queue, ok := topic.subscribers[subscriber]
	if !ok {
		return nil
	}

	if queue.Len() == 0 {
		return nil
	}

	el := queue.Front()
	queue.Remove(el)

	id, ok := el.Value.(int64)
	if !ok {
		return nil
	}

	data, found := topic.messages[id]
	if !found {
		return nil
	}

	topic.unreadCount[id] -= 1
	if topic.unreadCount[id] <= 0 {
		delete(topic.messages, id)
		delete(topic.unreadCount, id)
	}

	return data
}

func (topic *Topic) PutMessage(data json.RawMessage) {
	topic.Lock()

	topic.lastId += 1
	topic.messages[topic.lastId] = data
	topic.unreadCount[topic.lastId] = topic.subCount

	for name := range topic.subscribers {
		topic.subscribers[name].PushBack(topic.lastId)
	}

	topic.Unlock()
}

func (topic *Topic) Subscribe(subscriber string) {
	topic.Lock()
	if _, ok := topic.subscribers[subscriber]; !ok {
		topic.subscribers[subscriber] = list.New()
	}
	topic.subCount += 1
	topic.Unlock()
}

func (topic *Topic) Unsubscribe(subscriber string) {
	topic.Lock()
	defer topic.Unlock()

	if _, ok := topic.subscribers[subscriber]; !ok {
		return
	}

	delete(topic.subscribers, subscriber)

	topic.subCount -= 1
	for id, counter := range topic.unreadCount {
		counter -= 1
		if counter <= 0 {
			delete(topic.messages, id)
			delete(topic.unreadCount, id)
		}
		topic.unreadCount[id] = 0
	}
}
