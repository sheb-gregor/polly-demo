package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type Message struct {
	Topic string          `json:"topic"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type PollReq struct {
	Topic      string `json:"topic"`
	Subscriber string `json:"subscriber"`
}

// PollyClient is a client for the Polly Pub/Sub Server.
type PollyClient interface {
	// Poll receiving the next unseen message or no message if everything is seen,
	// or an error if the subscription is not found.
	Poll(topic, subscriber string) (json.RawMessage, error)
	// Publish send a new message to the topic.
	Publish(topic string, data json.RawMessage) error
	// Subscribe add a subscriber subscription to a topic.
	Subscribe(topic, subscriber string) error
	// Unsubscribe remove the subscription from the topic.
	Unsubscribe(topic, subscriber string) error
}

type client struct {
	url  url.URL
	http http.Client
}

// NewClient create new instance of the PollyClient.
func NewClient(baseAddr string) (PollyClient, error) {
	parsed, err := url.Parse(baseAddr)
	if err != nil {
		return nil, err
	}

	return &client{url: *parsed}, nil
}

func (client *client) Publish(topic string, data json.RawMessage) error {
	return client.postData("publish", Message{Topic: topic, Data: data})
}

func (client *client) Subscribe(topic, subscriber string) error {
	return client.postData("subscribe", PollReq{Topic: topic, Subscriber: subscriber})
}

func (client *client) Unsubscribe(topic, subscriber string) error {
	return client.postData("unsubscribe", PollReq{Topic: topic, Subscriber: subscriber})
}

func (client *client) Poll(topic, subscriber string) (json.RawMessage, error) {
	query := url.Values{}
	query.Set("topic", topic)
	query.Set("subscriber", subscriber)

	client.url.Path = "poll"
	client.url.RawQuery = query.Encode()
	resp, err := client.http.Get(client.url.String())
	if err != nil {
		return nil, errors.Wrap(err, "unable to send request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("request failed with status: " + resp.Status)
	}

	data := Message{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode response")
	}

	return data.Data, nil
}

func (client *client) postData(path string, body interface{}) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "unable to encode request")
	}

	client.url.Path = path
	resp, err := client.http.Post(client.url.String(),
		"application/json", bytes.NewBuffer(raw))
	if err != nil {
		return errors.Wrap(err, "unable to send request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("request failed with status: " + resp.Status)
	}

	return nil
}
