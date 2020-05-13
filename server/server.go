package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sheb-gregor/polly-demo/mq"
)

type Message struct {
	Topic string          `json:"topic"`
	Data  json.RawMessage `json:"data,omitempty"`
}

func (msg Message) Validate() error {
	if msg.Topic == "" {
		return errors.New("topic should not be empty")
	}

	if msg.Data == nil {
		return errors.New("data should not be empty")
	}
	return nil
}

type PollReq struct {
	Topic      string `json:"topic"`
	Subscriber string `json:"subscriber"`
}

func (msg PollReq) Validate() error {
	if msg.Topic == "" {
		return errors.New("topic should not be empty")
	}

	if msg.Subscriber == "" {
		return errors.New("subscriber should not be empty")
	}
	return nil
}

type StatusMsg struct {
	Message string `json:"message"`
}

func GetServer(broker *mq.Broker) http.Handler {
	mux := chi.NewMux()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/poll", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		req := PollReq{
			Topic:      query.Get("topic"),
			Subscriber: query.Get("subscriber"),
		}

		if err := req.Validate(); err != nil {
			writeError(w, err)
			return
		}

		msg, subscribed := broker.Poll(req.Topic, req.Subscriber)
		if !subscribed {
			writeData(w, http.StatusNotFound, StatusMsg{Message: http.StatusText(http.StatusNotFound)})
			return
		}
		writeSuccess(w, Message{Topic: req.Topic, Data: msg})
	})

	mux.Post("/publish", func(w http.ResponseWriter, r *http.Request) {
		req := Message{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeError(w, err)
			return
		}
		if err := req.Validate(); err != nil {
			writeError(w, err)
			return
		}

		broker.HandleNewMessage(req.Topic, req.Data)
		writeSuccess(w, StatusMsg{Message: http.StatusText(http.StatusOK)})
	})

	mux.Post("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		req := PollReq{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeError(w, err)
			return
		}
		if err := req.Validate(); err != nil {
			writeError(w, err)
			return
		}

		broker.Subscribe(req.Topic, req.Subscriber)
		writeSuccess(w, StatusMsg{Message: http.StatusText(http.StatusOK)})
	})

	mux.Post("/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		req := PollReq{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeError(w, err)
			return
		}
		if err := req.Validate(); err != nil {
			writeError(w, err)
			return
		}

		broker.Unsubscribe(req.Topic, req.Subscriber)
		writeSuccess(w, StatusMsg{Message: http.StatusText(http.StatusOK)})
	})

	return mux
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	writeData(w, http.StatusOK, data)
}

func writeError(w http.ResponseWriter, err error) {
	log.Println("ERROR: processing error", err.Error())

	writeData(w, http.StatusBadRequest, StatusMsg{Message: err.Error()})
}

func writeData(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("ERROR: filed to write response ", err.Error())
	}
}
