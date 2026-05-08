package secrets

import "errors"

var (
	ErrSecretLocked      = errors.New("secret service is locked")
	ErrSecretNotFound    = errors.New("secret not found")
	ErrUnexpectedType    = errors.New("unexpected type")
	ErrAccessingProperty = errors.New("error accessing property")
)
