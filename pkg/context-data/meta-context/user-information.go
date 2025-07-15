package metacontext

import (
	"context"
)

// This struct defines the UserInformationMeta struct
//
//	It can be used to store metadata about the request
type UserInformationMeta struct {
	UserID   int64
	Username string
	Email    string
	Roles    []string
}

// This struct defines the UserInformationMetaKeyType struct
//
//	It is used as a key for storing and retrieving UserInformationMeta from the context
type UserInformationMetaKeyType struct{}

// Define a key for storing UserInformationMeta in the context
var userInformationMetaKey = UserInformationMetaKeyType{}

// InjectUserInformationMeta injects the UserInformationMeta into the context.
// This function is used to add metadata to the context for later retrieval
func InjectUserInformationMeta(ctx context.Context, meta UserInformationMeta) context.Context {
	return context.WithValue(ctx, userInformationMetaKey, meta)
}

// ExtractUserInformationMeta retrieves the UserInformationMeta from the context.
// This function is used to access the metadata stored in the context
func ExtractUserInformationMeta(ctx context.Context) (UserInformationMeta, bool) {
	meta, ok := ctx.Value(userInformationMetaKey).(UserInformationMeta)
	return meta, ok
}
