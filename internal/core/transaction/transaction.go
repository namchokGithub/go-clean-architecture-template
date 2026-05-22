package transaction

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// NewContext embeds a GORM transaction into the context.
func NewContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// FromContext extracts the transaction from context. Returns (nil, false) if none.
func FromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok && tx != nil
}
