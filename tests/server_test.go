package tests

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/lancer-kit/uwe/v2/presets/api"
	"github.com/sheb-gregor/polly-demo/client"
	"github.com/sheb-gregor/polly-demo/mq"
	"github.com/sheb-gregor/polly-demo/server"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	name := "bob"
	topic := "test_topic"
	message := json.RawMessage(`{"my_key":"my_message"}`)
	broker := mq.NewBroker()

	ctx, cancel := context.WithCancel(context.Background())
	apiServer := api.NewServer(api.Config{Host: "0.0.0.0", Port: 8080}, server.GetServer(broker))
	go func() { _ = apiServer.Run(ctx) }()

	pClient, err := client.NewClient("http://0.0.0.0:8080")
	assert.NoError(t, err)

	err = pClient.Subscribe(topic, name)
	assert.NoError(t, err)

	err = pClient.Publish(topic, message)
	assert.NoError(t, err)

	msg, err := pClient.Poll(topic, name)
	assert.NoError(t, err)
	assert.Equal(t, message, msg)

	err = pClient.Unsubscribe(topic, name)
	assert.NoError(t, err)
	cancel()
}
