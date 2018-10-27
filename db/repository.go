package db

import (
	"context"

	"github.com/gonzalesraul/meow/schema"
)

// Repository interface for Command services
type Repository interface {
	Close()
	InsertMeow(ctx context.Context, meow schema.Meow) error
	ListMeows(ctx context.Context, skip uint64, take uint64) ([]schema.Meow, error)
}

var impl Repository

// SetRepository - set the repository connection
func SetRepository(repository Repository) {
	impl = repository
}

// Close - Closes repository connection
func Close() {
	impl.Close()
}

// InsertMeow - InsertMeow implementation
func InsertMeow(ctx context.Context, meow schema.Meow) error {
	return impl.InsertMeow(ctx, meow)
}

// ListMeow - InsertMeow implementation
func ListMeow(ctx context.Context, skip uint64, take uint64) ([]schema.Meow, error) {
	return impl.ListMeows(ctx, skip, take)
}
