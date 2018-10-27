package event

import (
	"bytes"
	"encoding/gob"

	"github.com/gonzalesraul/meow/schema"
	nats "github.com/nats-io/go-nats"
)

// NatEventStore NATS implementation for EventStore interface
type NatEventStore struct {
	nc                      *nats.Conn
	meowCreatedSubscription *nats.Subscription
	meowCreatedChan         chan MeowCreatedMessage
}

// NewNats creates stream connection with NATS
func NewNats(url string) (*NatEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatEventStore{nc: nc}, nil
}

// Close NATS connection
func (e *NatEventStore) Close() {
	if e.nc != nil {
		e.Close()
	}
	if e.meowCreatedSubscription != nil {
		e.meowCreatedSubscription.Unsubscribe()
	}
	close(e.meowCreatedChan)
}

// PublishMeowCreated inserts MeowCreated event
func (e *NatEventStore) PublishMeowCreated(meow schema.Meow) error {
	m := MeowCreatedMessage{meow.ID, meow.Body, meow.CreatedAt}
	data, err := e.writeMessage(&m)
	if err != nil {
		return err
	}
	return e.nc.Publish(m.Key(), data)
}

// SubscribeMeowCreated Decode events from NATS and routes to a channel with MeowCreated event
func (e *NatEventStore) SubscribeMeowCreated() (<-chan MeowCreatedMessage, error) {
	m := MeowCreatedMessage{}
	e.meowCreatedChan = make(chan MeowCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)

	var err error
	e.meowCreatedSubscription, err = e.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				e.readMessage(msg.Data, &m)
				e.meowCreatedChan <- m
			}
		}
	}()
	return (<-chan MeowCreatedMessage)(e.meowCreatedChan), nil
}

// OnMeowCreated attach a function to handle MeowCreated event
func (e *NatEventStore) OnMeowCreated(fn func(MeowCreatedMessage)) (err error) {
	m := MeowCreatedMessage{}
	e.meowCreatedSubscription, err = e.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		e.readMessage(msg.Data, &m)
		fn(m)
	})
	return
}

//writeMessage transforms the message into a bytes to be transported through NATs
func (e *NatEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//readMessage transforms the bytes from NATS to the message
func (e *NatEventStore) readMessage(data []byte, m *MeowCreatedMessage) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}
