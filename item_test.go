package secrets_test

import (
	"testing"

	"github.com/Loggie/go-secrets"
	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/require"
)

func firstItemOrSkip(t *testing.T, svc *secrets.Service) *secrets.Item {
	t.Helper()

	collections, err := svc.Collections()
	require.NoError(t, err)

	for _, col := range collections {
		items, errItems := col.Items()
		require.NoError(t, errItems)
		if len(items) > 0 {
			return items[0]
		}
	}

	require.FailNow(t, "no items available in collections")
	return nil
}

func TestNewItem(t *testing.T) {
	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		if errClose := conn.Close(); errClose != nil {
			t.Fatal(errClose)
		}
	}()

	path := secrets.ItemPath(secrets.CollectionPath("test"), "item1")
	item := secrets.NewItem(conn, path)
	require.Equal(t, path, item.Path())
}

func TestItemProperties(t *testing.T) {
	svc := requireService(t)
	collection := requireAnyCollection(t, svc)
	require.NoError(t, svc.Unlock([]secrets.LockableObject{collection}))

	attrs := uniqueAttributes("item-properties")
	item := requireCreateItem(t, svc, collection, attrs, "item-properties-secret")
	defer func() {
		_, _ = item.Delete()
	}()

	_, err := item.Locked()
	require.NoError(t, err)

	_, err = item.Attributes()
	require.NoError(t, err)

	_, err = item.Label()
	require.NoError(t, err)

	created, err := item.Created()
	require.NoError(t, err)
	require.False(t, created.IsZero())

	modified, err := item.Modified()
	require.NoError(t, err)
	require.False(t, modified.IsZero())
}

func TestItemPropertiesInvalidPath(t *testing.T) {
	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	bad := secrets.NewItem(conn, secrets.ItemPath(secrets.CollectionPath("nope"), "nope"))

	_, err = bad.Locked()
	require.Error(t, err)

	_, err = bad.Attributes()
	require.Error(t, err)

	_, err = bad.Label()
	require.Error(t, err)

	_, err = bad.Created()
	require.Error(t, err)

	_, err = bad.Modified()
	require.Error(t, err)
}

func TestItemGetSecret(t *testing.T) {
	svc := requireService(t)
	item := firstItemOrSkip(t, svc)
	session := requireOpenSession(t, svc)
	defer func() {
		_ = session.Close()
	}()

	secret, err := item.GetSecret(*session)
	require.NoError(t, err)
	require.NotNil(t, secret.Value)
}

func TestItemGetSecretLockedOrInvalidSession(t *testing.T) {
	svc := requireService(t)
	item := firstItemOrSkip(t, svc)

	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	badSession := secrets.NewSession(conn, secrets.SessionPath("does-not-exist"))
	_, err = item.GetSecret(*badSession)
	require.Error(t, err)
}

func TestItemDeleteInvalidPathReturnsError(t *testing.T) {
	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	bad := secrets.NewItem(conn, secrets.ItemPath(secrets.CollectionPath("nope"), "nope"))

	_, err = bad.Delete()
	require.Error(t, err)
}

func TestItemSetSecretInvalidPathReturnsError(t *testing.T) {
	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	bad := secrets.NewItem(conn, secrets.ItemPath(secrets.CollectionPath("nope"), "nope"))

	err = bad.SetSecret(secrets.Secret{Value: []byte("v"), ContentType: "text/plain"})
	require.Error(t, err)
}

func TestItemSetSecretSuccess(t *testing.T) {
	svc := requireService(t)
	collection := requireAnyCollection(t, svc)
	require.NoError(t, svc.Unlock([]secrets.LockableObject{collection}))

	attrs := uniqueAttributes("item-set-secret")
	item := requireCreateItem(t, svc, collection, attrs, "original")

	session := requireOpenSession(t, svc)
	defer func() {
		_ = session.Close()
	}()

	err := item.SetSecret(secrets.Secret{
		Path:        session.Path(),
		Value:       []byte("updated"),
		ContentType: "text/plain",
	})
	require.NoError(t, err)

	_, _ = item.Delete()
}
