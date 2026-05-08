package secrets_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Loggie/go-secrets"
	"github.com/stretchr/testify/require"
)

func requireService(t *testing.T) *secrets.Service {
	t.Helper()

	svc, err := secrets.NewService()
	require.NoError(t, err)

	return svc
}

func requireOpenSession(t *testing.T, svc *secrets.Service) *secrets.Session {
	t.Helper()

	session, err := svc.OpenSession(secrets.SecretServiceSessionAlgorithmPlain, "")
	require.NoError(t, err)
	require.NotNil(t, session)

	return session
}

func requireAnyCollection(t *testing.T, svc *secrets.Service) *secrets.Collection {
	t.Helper()

	collections, err := svc.Collections()
	require.NoError(t, err)
	require.NotEmpty(t, collections, "secret service has no collections")

	return collections[0]
}

func uniqueAttributes(prefix string) secrets.Attributes {
	return secrets.Attributes{
		"go-secrets-test-key": fmt.Sprintf("%s-%d", prefix, time.Now().Unix()),
	}
}

func requireCreateItem(
	t *testing.T,
	svc *secrets.Service,
	collection *secrets.Collection,
	attrs secrets.Attributes,
	value string,
) *secrets.Item {
	t.Helper()
	require.NoError(t, svc.Unlock([]secrets.LockableObject{collection}))

	session := requireOpenSession(t, svc)
	defer func() {
		_ = session.Close()
	}()

	item, err := collection.CreateItem(
		"go-secrets-test-item",
		attrs,
		secrets.Secret{
			Path:        session.Path(),
			Value:       []byte(value),
			ContentType: "text/plain",
		},
		true,
	)
	require.NoError(t, err)
	require.NotNil(t, item)

	return item
}
