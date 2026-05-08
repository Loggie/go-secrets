package secrets_test

import (
	"testing"

	"github.com/Loggie/go-secrets"
	"github.com/stretchr/testify/require"
)

func TestAddEmptyCollectionUsesDefault(t *testing.T) {
	attrs := uniqueAttributes("helpers-add-default")
	want := "helper-default-secret"

	item, err := secrets.Add("", "go-secrets-helper-default-item", "text/plain", want, attrs)
	require.NoError(t, err)
	require.NotNil(t, item)

	got, err := secrets.Get(attrs)
	require.NoError(t, err)
	require.Equal(t, want, got)

	_, _ = item.Delete()
}

func TestAddHelperAndGet(t *testing.T) {
	attrs := uniqueAttributes("helpers-create-item")
	want := "helper-created-secret"

	item, err := secrets.Add("default", "go-secrets-helper-item", "text/plain", want, attrs)
	require.NoError(t, err)
	require.NotNil(t, item)

	got, err := secrets.Get(attrs)
	require.NoError(t, err)
	require.Equal(t, want, got)

	_, _ = item.Delete()
}

func TestGetNoMatchReturnsEmpty(t *testing.T) {
	got, err := secrets.Get(uniqueAttributes("helpers-get"))
	require.NoError(t, err)
	require.Empty(t, got)
}

func TestGetFoundSecret(t *testing.T) {
	svc := requireService(t)
	collection := requireAnyCollection(t, svc)
	require.NoError(t, svc.Unlock([]secrets.LockableObject{collection}))

	attrs := uniqueAttributes("helpers-get-found")
	want := "helpers-get-secret-value"
	item := requireCreateItem(t, svc, collection, attrs, want)

	got, err := secrets.Get(attrs)
	require.NoError(t, err)
	require.Equal(t, want, got)

	_, _ = item.Delete()
}
