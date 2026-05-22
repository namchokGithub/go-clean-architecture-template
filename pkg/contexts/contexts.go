package contexts

import "context"

// Key is the typed context key to prevent collisions with other packages.
type Key string

// Set stores a value under a typed key.
func Set(ctx context.Context, key Key, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

// Get retrieves a value and casts it to T. Returns zero value + false if missing or wrong type.
func Get[T any](ctx context.Context, key Key) (T, bool) {
	v, ok := ctx.Value(key).(T)
	return v, ok
}
