package schema

import (
	"time"
)

// Meow schema
type Meow struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
