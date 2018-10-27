package event

import "time"

//Message interface
type Message interface {
	Key() string
}

//MeowCreatedMessage Meow messag
type MeowCreatedMessage struct {
	ID        string
	Body      string
	CreatedAt time.Time
}

//Key gets message key
func (m *MeowCreatedMessage) Key() string {
	return "meow.created"
}
