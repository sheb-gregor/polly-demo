package mq

import (
	"container/list"
	"encoding/json"
	"sync"
)

type Topic struct {
	sync.Mutex

	subCount int64
	lastID   int64

	// subscribers - this map contains Unread Message Queues (FIFOs) for each user,
	// the value of the list item is the message identifier.
	subscribers map[string]*list.List
	// unreadCount map contains counters that show how many subscribers have not yet received each message.
	unreadCount map[int64]int64
	// messages map with messages, key is unique ID of message
	messages map[int64]json.RawMessage
}

// NewTopic creates and initialize new topic instance.
func NewTopic() *Topic {
	return &Topic{
		subscribers: map[string]*list.List{},
		unreadCount: map[int64]int64{},
		messages:    map[int64]json.RawMessage{},
	}
}

// Poll checks if the subscriber exists and retrieves the last unread message from the queue.
func (topic *Topic) Poll(subscriber string) (json.RawMessage, bool) {
	topic.Lock()
	defer topic.Unlock()

	queue, ok := topic.subscribers[subscriber]
	if !ok {
		return nil, false
	}

	if queue.Len() == 0 {
		return nil, true
	}

	el := queue.Front()
	queue.Remove(el)

	return topic.fetchMessage(el), true
}

// PutMessage adds a new message to this topic, increases the message lastID
// and sets this message as unread for all subscribers.
func (topic *Topic) PutMessage(data json.RawMessage) {
	topic.Lock()

	topic.lastID += 1
	topic.messages[topic.lastID] = data
	topic.unreadCount[topic.lastID] = topic.subCount

	for name := range topic.subscribers {
		topic.subscribers[name].PushBack(topic.lastID)
	}

	topic.Unlock()
}

// Poll adds a new subscriber to this topic, increases the counter of the total number of subscribers.
func (topic *Topic) Subscribe(subscriber string) {
	topic.Lock()
	if _, ok := topic.subscribers[subscriber]; !ok {
		topic.subscribers[subscriber] = list.New()
	}
	topic.subCount += 1
	topic.Unlock()
}

// Poll removes a subscriber from this topic, decreases the counter of the total number of subscribers.
// Also deletes all messages for which this subscriber was the last who did not receive.
func (topic *Topic) Unsubscribe(subscriber string) {
	topic.Lock()
	defer topic.Unlock()

	queue, ok := topic.subscribers[subscriber]
	if !ok {
		return
	}

	for {
		if queue.Len() == 0 {
			break
		}
		el := queue.Front()
		if el == nil {
			break
		}

		queue.Remove(el)

		topic.fetchMessage(el)
	}

	delete(topic.subscribers, subscriber)
	topic.subCount -= 1
}

func (topic *Topic) fetchMessage(el *list.Element) json.RawMessage {
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
