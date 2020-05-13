package mq

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroker_Subscribe(t *testing.T) {
	broker := NewBroker()
	name := "bob"
	topics := []string{"test_1", "test_2", "test_3", "test_4", "test_5"}

	for _, topic := range topics {
		_, ok := broker.topics.Load(topic)
		assert.False(t, ok)

		broker.Subscribe(topic, name)
		topicObj, ok := broker.topics.Load(topic)
		assert.True(t, ok)
		assert.NotNil(t, topicObj)

		tp, ok := topicObj.(*Topic)
		assert.True(t, ok)

		list, ok := tp.subscribers[name]
		assert.True(t, ok)
		assert.NotNil(t, list)
	}
}

func TestBroker_Unsubscribe(t *testing.T) {
	broker := NewBroker()
	name := "bob"
	topics := []string{"test_1", "test_2", "test_3", "test_4", "test_5"}

	for _, topic := range topics {
		_, ok := broker.topics.Load(topic)
		assert.False(t, ok)

		broker.Unsubscribe(topic, name)
		_, ok = broker.topics.Load(topic)
		assert.False(t, ok)

		broker.Subscribe(topic, name)
		topicObj, ok := broker.topics.Load(topic)
		assert.True(t, ok)
		assert.NotNil(t, topicObj)

		broker.Unsubscribe(topic, name)
		_, ok = broker.topics.Load(topic)
		assert.False(t, ok)
	}
}

func TestBroker_HandleNewMessage(t *testing.T) {
	broker := NewBroker()
	name := "bob"
	message := json.RawMessage("test_message")
	topics := []string{"test_1", "test_2", "test_3", "test_4", "test_5"}

	for _, topic := range topics {
		broker.HandleNewMessage(topic, message)
		_, ok := broker.topics.Load(topic)
		assert.False(t, ok)

		broker.Subscribe(topic, name)
		topicObj, ok := broker.topics.Load(topic)
		assert.True(t, ok)
		assert.NotNil(t, topicObj)

		broker.HandleNewMessage(topic, message)
		topicObj, ok = broker.topics.Load(topic)
		assert.True(t, ok)
		assert.NotNil(t, topicObj)

		tp, ok := topicObj.(*Topic)
		assert.True(t, ok)
		assert.Equal(t, 1, len(tp.messages))
		assert.Equal(t, 1, len(tp.unreadCount))
	}
}

func TestBroker_Poll(t *testing.T) {
	broker := NewBroker()
	name := "bob"
	message := json.RawMessage("test_message")
	topics := []string{"test_1", "test_2", "test_3", "test_4", "test_5"}

	for _, topic := range topics {
		broker.HandleNewMessage(topic, message)
		_, ok := broker.topics.Load(topic)
		assert.False(t, ok)

		broker.Subscribe(topic, name)
		topicObj, ok := broker.topics.Load(topic)
		assert.True(t, ok)
		assert.NotNil(t, topicObj)

		broker.HandleNewMessage(topic, message)
		topicObj, ok = broker.topics.Load(topic)
		assert.True(t, ok)
		assert.NotNil(t, topicObj)

		tp, ok := topicObj.(*Topic)
		assert.True(t, ok)
		assert.Equal(t, 1, len(tp.messages))
		assert.Equal(t, 1, len(tp.unreadCount))

		msg, subscribed := broker.Poll(topic, name)
		assert.True(t, subscribed)
		assert.Equal(t, message, msg)

		topicObj, ok = broker.topics.Load(topic)
		assert.True(t, ok)
		assert.NotNil(t, topicObj)

		tp, ok = topicObj.(*Topic)
		assert.True(t, ok)
		assert.Equal(t, 0, len(tp.messages))
		assert.Equal(t, 0, len(tp.unreadCount))
	}
}
