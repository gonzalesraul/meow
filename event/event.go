package event

import (
	"github.com/gonzalesraul/meow/schema"
)

// EventStore .
type EventStore interface {
	Close()
	PublishMeowCreated(meow schema.Meow) error
	SubscribeMeowCreated() (<-chan MeowCreatedMessage, error)
	OnMeowCreated(f func(MeowCreatedMessage)) error
}

var impl EventStore

// SetEventStore .
func SetEventStore(es EventStore) {
	impl = es
}

// Close .
func Close() {
	impl.Close()
}

// PublishMeowCreated - Push MeowCreated event to the EventStore
func PublishMeowCreated(meow schema.Meow) error {
	return impl.PublishMeowCreated(meow)
}

// SubscribeMeowCreated - Receives MeowCreated events
func SubscribeMeowCreated() (<-chan MeowCreatedMessage, error) {
	return impl.SubscribeMeowCreated()
}

// OnMeowCreated - Trigger on MeowCreated event
func OnMeowCreated(f func(MeowCreatedMessage)) error {
	return impl.OnMeowCreated(f)
}
