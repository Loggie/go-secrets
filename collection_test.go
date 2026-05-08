package secrets_test

import (
	"testing"

	"github.com/Loggie/go-secrets"
	"github.com/stretchr/testify/require"
)

func TestCollectionNew(t *testing.T) {
	svc := requireService(t)
	col := secrets.NewCollection(svc, secrets.CollectionPath("testCollection"))
	require.NotNil(t, col)
	require.Equal(t, secrets.CollectionPath("testCollection"), col.Path())
}

func TestCollectionItems(t *testing.T) {
	svc := requireService(t)
	col := requireAnyCollection(t, svc)

	items, err := col.Items()
	require.NoError(t, err)
	require.NotNil(t, items)
}

func TestCollectionItemsInvalidPath(t *testing.T) {
	svc := requireService(t)
	bad := secrets.NewCollection(svc, secrets.CollectionPath("invalid"))

	_, err := bad.Items()
	require.Error(t, err)
}

func TestCollectionLabel(t *testing.T) {
	svc := requireService(t)
	col := requireAnyCollection(t, svc)

	label, err := col.Label()
	require.NoError(t, err)
	require.NotNil(t, []byte(label))
}

func TestCollectionLabelInvalidPath(t *testing.T) {
	svc := requireService(t)
	bad := secrets.NewCollection(svc, secrets.CollectionPath("invalid"))

	_, err := bad.Label()
	require.Error(t, err)
}

func TestCollectionSetLabel(t *testing.T) {
	svc := requireService(t)
	col := requireAnyCollection(t, svc)
	require.NoError(t, svc.Unlock([]secrets.LockableObject{col}))

	label, err := col.Label()
	require.NoError(t, err)

	err = col.SetLabel(label)
	require.NoError(t, err)
}

func TestCollectionLocked(t *testing.T) {
	svc := requireService(t)
	col := requireAnyCollection(t, svc)

	_, err := col.Locked()
	require.NoError(t, err)
}

func TestCollectionLockedInvalidPath(t *testing.T) {
	svc := requireService(t)
	bad := secrets.NewCollection(svc, secrets.CollectionPath("invalid"))

	_, err := bad.Locked()
	require.Error(t, err)
}

func TestCollectionCreated(t *testing.T) {
	svc := requireService(t)
	col := requireAnyCollection(t, svc)

	created, err := col.Created()
	require.NoError(t, err)
	require.False(t, created.IsZero())
}

func TestCollectionCreatedInvalidPath(t *testing.T) {
	svc := requireService(t)
	bad := secrets.NewCollection(svc, secrets.CollectionPath("invalid"))

	_, err := bad.Created()
	require.Error(t, err)
}

func TestCollectionModified(t *testing.T) {
	svc := requireService(t)
	col := requireAnyCollection(t, svc)

	modified, err := col.Modified()
	require.NoError(t, err)
	require.False(t, modified.IsZero())
}

func TestCollectionModifiedInvalidPath(t *testing.T) {
	svc := requireService(t)
	bad := secrets.NewCollection(svc, secrets.CollectionPath("invalid"))

	_, err := bad.Modified()
	require.Error(t, err)
}

func TestCollectionDeleteInvalidPath(t *testing.T) {
	svc := requireService(t)
	bad := secrets.NewCollection(svc, secrets.CollectionPath("invalid"))

	_, err := bad.Delete()
	require.Error(t, err)
}

func TestCollectionSearchItems(t *testing.T) {
	svc := requireService(t)
	col := requireAnyCollection(t, svc)

	items, err := col.SearchItems(uniqueAttributes("collection-search"))
	require.NoError(t, err)
	require.Empty(t, items)
}

func TestCollectionSearchItemsInvalidPath(t *testing.T) {
	svc := requireService(t)
	bad := secrets.NewCollection(svc, secrets.CollectionPath("invalid"))

	_, err := bad.SearchItems(uniqueAttributes("collection-search-invalid"))
	require.Error(t, err)
}

func TestCollectionCreateItemInvalidPath(t *testing.T) {
	svc := requireService(t)
	bad := secrets.NewCollection(svc, secrets.CollectionPath("invalid"))

	_, err := bad.CreateItem(
		"go-secrets-temp",
		secrets.Attributes{},
		secrets.Secret{Value: []byte("value"), ContentType: "text/plain"},
		false,
	)
	require.Error(t, err)
}

func TestCollectionCreateItem(t *testing.T) {
	svc := requireService(t)
	col := requireAnyCollection(t, svc)
	require.NoError(t, svc.Unlock([]secrets.LockableObject{col}))

	session := requireOpenSession(t, svc)
	defer func() {
		_ = session.Close()
	}()

	item, err := col.CreateItem(
		"go-secrets-collection-create-item",
		secrets.Attributes{"go-secrets-create-item": "true"},
		secrets.Secret{
			Path:        session.Path(),
			Value:       []byte("create-item-value"),
			ContentType: "text/plain",
		},
		true,
	)
	require.NoError(t, err)
	require.NotNil(t, item)
	require.NotEmpty(t, item.Path())
}
