package secrets

import (
	"strings"
)

// Add creates or replaces an item in a collection without requiring callers
// to manage session creation or unlock flow.
func Add(
	collection string,
	name string,
	contentType string,
	value string,
	attrs Attributes,
) (*Item, error) {
	collection = strings.TrimSpace(collection)
	if collection == "" {
		collection = "default"
	}

	service, err := NewService()
	if err != nil {
		return nil, err
	}

	col, err := service.ReadAlias(collection)
	if err != nil || col.Path() == "/" {
		path := CollectionPath(collection)
		col = NewCollection(service, path)
	}

	if errUnlock := service.Unlock([]LockableObject{col}); errUnlock != nil {
		return nil, errUnlock
	}

	session, err := service.OpenSession(
		SecretServiceSessionAlgorithmPlain,
		"",
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = session.Close()
	}()

	return col.CreateItem(
		name,
		attrs,
		Secret{
			Path:        session.Path(),
			Value:       []byte(value),
			ContentType: contentType,
		},
		true,
	)
}

// Get retrieves a secret from the secret service based on the provided attributes.
// It returns the secret value as a string, or an error if the retrieval fails.
func Get(
	attrs Attributes,
) (string, error) {
	service, err := NewService()
	if err != nil {
		return "", err
	}

	session, err := service.OpenSession(
		SecretServiceSessionAlgorithmPlain,
		"",
	)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = session.Close()
	}()

	items, locked, err := service.SearchItems(
		attrs,
	)
	if err != nil {
		return "", err
	}

	items = append(items, locked...)

	if len(items) == 0 {
		return "", nil
	}

	objects := make([]LockableObject, 0, len(items))
	for _, item := range items {
		objects = append(objects, item)
	}

	err = service.Unlock(objects)
	if err != nil {
		return "", err
	}

	secrets, err := service.GetSecrets(items, session)
	if err != nil {
		return "", err
	}

	for _, secret := range secrets {
		return string(secret.Value), nil
	}

	return "", ErrSecretNotFound
}
