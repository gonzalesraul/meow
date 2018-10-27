package search

import (
	"context"

	"github.com/gonzalesraul/meow/schema"
)

//Repository interface for Querying services
type Repository interface {
	Close()
	InsertMeow(ctx context.Context, meow schema.Meow) error
	InquiryMeow(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Meow, error)
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

// InquiryMeow - InquiryMeow implementation
func InquiryMeow(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Meow, error) {
	return impl.InquiryMeow(ctx, query, skip, take)
}
