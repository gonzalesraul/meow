package main

import (
	"time"
)

const (
	//KindMeowCreated - Incremental controller
	KindMeowCreated = iota + 1
)

//MeowCreatedMessage - message to push through websocket
type MeowCreatedMessage struct {
	Kind      uint32    `json:"kind"`
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func newMeowCreatedMessage(id string, body string, createdAt time.Time) *MeowCreatedMessage {
	return &MeowCreatedMessage{
		Kind:      KindMeowCreated,
		ID:        id,
		Body:      body,
		CreatedAt: createdAt,
	}
}
