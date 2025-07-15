package metacontext

import (
	"context"
	"fmt"
)

// GetValueFromContext retrieves a value from the context using the provided key.
// It returns the value and an error if the key does not exist in the context.
func GetValueFromContext(ctx context.Context, key string) (interface{}, error) {
	value := ctx.Value(key)
	if value == nil {
		return nil, fmt.Errorf("key %s not found in context", key)
	}
	return value, nil
}
