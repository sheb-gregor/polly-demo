package server

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopic_New(t *testing.T) {
	topic := NewTopic()

	assert.NotNil(t, topic.subscribers)
	assert.NotNil(t, topic.messages)
	assert.NotNil(t, topic.unreadCount)

	assert.Equal(t, int64(0), topic.subCount)
	assert.Equal(t, int64(0), topic.lastId)
	assert.Equal(t, 0, len(topic.subscribers))
	assert.Equal(t, 0, len(topic.messages))
	assert.Equal(t, 0, len(topic.unreadCount))

}

func TestTopic_Subscribe(t *testing.T) {
	topic := NewTopic()

	names := []string{"alice", "bob", "dave", "flash", "kilo"}
	for i, name := range names {
		topic.Subscribe(name)
		assert.Equal(t, int64(i+1), topic.subCount)
		assert.Equal(t, i+1, len(topic.subscribers))
		assert.Equal(t, 0, len(topic.messages))
		assert.Equal(t, 0, len(topic.unreadCount))

		list, ok := topic.subscribers[name]
		assert.True(t, ok)
		assert.NotNil(t, list)
		assert.Equal(t, 0, list.Len())
	}

	assert.Equal(t, int64(len(names)), topic.subCount)
	assert.Equal(t, len(names), len(topic.subscribers))
	assert.Equal(t, int64(0), topic.lastId)

	unknownNames := []string{"alpha", "bravo", "eva", "travis"}
	for _, name := range unknownNames {
		list, ok := topic.subscribers[name]
		assert.False(t, ok)
		assert.Nil(t, list)
	}

}

func TestTopic_Unsubscribe(t *testing.T) {
	topic := NewTopic()
	names := []string{"alice", "bob", "dave", "flash", "kilo"}
	for _, name := range names {
		topic.Subscribe(name)
	}

	count := topic.subCount
	for i, name := range names {
		topic.Unsubscribe(name)

		assert.Equal(t, count-int64(i+1), topic.subCount)
		assert.Equal(t, count-int64(i+1), int64(len(topic.subscribers)))
		assert.Equal(t, 0, len(topic.messages))
		assert.Equal(t, 0, len(topic.unreadCount))

		list, ok := topic.subscribers[name]
		assert.False(t, ok)
		assert.Nil(t, list)
	}

	assert.Equal(t, int64(0), topic.subCount)
}

func TestTopic_PutMessage(t *testing.T) {
	topic := NewTopic()

	names := []string{"alice", "bob", "dave", "flash", "kilo"}
	for _, name := range names {
		topic.Subscribe(name)
	}

	count := topic.subCount

	messages := []json.RawMessage{
		json.RawMessage("test_1"),
		json.RawMessage("test_2"),
		json.RawMessage("test_3"),
		json.RawMessage("test_4"),
		json.RawMessage("test_5"),
		json.RawMessage("test_6"),
		json.RawMessage("test_7"),
	}

	for i, msg := range messages {
		topic.PutMessage(msg)
		id := int64(i) + 1

		assert.Equal(t, count, topic.subCount)
		assert.Equal(t, id, int64(len(topic.messages)))
		assert.Equal(t, id, int64(len(topic.unreadCount)))

		m, ok := topic.messages[id]
		assert.True(t, ok)
		assert.Equal(t, msg, m)

		count, ok := topic.unreadCount[id]
		assert.True(t, ok)
		assert.Equal(t, len(names), int(count))

		for _, name := range names {
			list, ok := topic.subscribers[name]
			assert.True(t, ok)
			assert.NotNil(t, list)
			assert.Equal(t, int(id), list.Len())
		}
	}

	assert.Equal(t, len(messages), int(topic.lastId))

	for i, message := range messages {
		m := topic.messages[int64(i+1)]
		assert.Equal(t, message, m)
	}
}

func TestTopic_Poll(t *testing.T) {
	topic := NewTopic()

	names := []string{"alice", "bob", "dave", "flash", "kilo"}
	messages := []json.RawMessage{
		json.RawMessage("test_1"),
		json.RawMessage("test_2"),
		json.RawMessage("test_3"),
		json.RawMessage("test_4"),
		json.RawMessage("test_5"),
		json.RawMessage("test_6"),
		json.RawMessage("test_7"),
	}

	for _, name := range names {
		topic.Subscribe(name)
	}
	subCount := topic.subCount

	for _, msg := range messages {
		topic.PutMessage(msg)
	}
	msgCount := int(topic.lastId)

	for subIndex, name := range names {
		list, ok := topic.subscribers[name]
		assert.True(t, ok)
		assert.NotNil(t, list)
		assert.Equal(t, msgCount, list.Len())

		for msgID := 0; msgID < msgCount; msgID++ {
			message := topic.Poll(name)

			list, ok := topic.subscribers[name]
			assert.True(t, ok)
			assert.NotNil(t, list)
			assert.Equal(t, msgCount-(msgID+1), list.Len())
			assert.Equal(t, messages[msgID], message)

			unreadCount := topic.unreadCount[int64(msgID+1)]
			assert.Equal(t, subCount-int64(subIndex+1), unreadCount)
		}
	}

	for id := range messages {
		_, ok := topic.messages[int64(id+1)]
		assert.False(t, ok)
		_, ok = topic.unreadCount[int64(id+1)]
		assert.False(t, ok)
	}
}
