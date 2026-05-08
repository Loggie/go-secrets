package secrets_test

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/Loggie/go-secrets"
	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/require"
)

type fakeLockable struct {
	path dbus.ObjectPath
}

func (f fakeLockable) Locked() (bool, error) {
	return false, nil
}

func (f fakeLockable) Path() dbus.ObjectPath {
	return f.path
}

func TestNewService(t *testing.T) {
	t.Parallel()

	_ = requireService(t)
}

func TestServiceCollections(t *testing.T) {
	t.Parallel()

	svc := requireService(t)

	collections, err := svc.Collections()
	require.NoError(t, err)
	require.NotNil(t, collections)
}

func TestServiceOpenSession(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	session := requireOpenSession(t, svc)
	require.NoError(t, session.Close())
}

func TestServiceCreateCollection(t *testing.T) {
	t.Parallel()

	svc := requireService(t)

	name := fmt.Sprintf("go-secrets-test-%d", time.Now().Unix())
	collection, err := svc.CreateCollection(name, "")
	require.NoError(t, err)

	require.NotNil(t, collection)
	require.NotEmpty(t, collection.Path())

	_, _ = collection.Delete()
}

func TestServiceCreateCollectionAndAddItem(t *testing.T) {
	t.Parallel()

	svc := requireService(t)

	name := fmt.Sprintf("go-secrets-e2e-%d", time.Now().Unix())
	collection, err := svc.CreateCollection(name, "")
	require.NoError(t, err)
	require.NotNil(t, collection)

	require.NoError(t, svc.Unlock([]secrets.LockableObject{collection}))

	attrs := uniqueAttributes("create-collection-item")
	value := fmt.Appendf(nil, "value-%d", time.Now().Unix())

	session := requireOpenSession(t, svc)
	defer func() {
		_ = session.Close()
	}()

	item, err := collection.CreateItem(
		"go-secrets-e2e-item",
		attrs,
		secrets.Secret{
			Path:        session.Path(),
			Value:       value,
			ContentType: "text/plain",
		},
		true,
	)
	require.NoError(t, err)
	require.NotNil(t, item)

	found, err := collection.SearchItems(attrs)
	require.NoError(t, err)
	require.NotEmpty(t, found)

	idx := slices.IndexFunc(found, func(i *secrets.Item) bool {
		return i.Path() == item.Path()
	})
	require.NotEqual(t, -1, idx)

	_, _ = item.Delete()
	_, _ = collection.Delete()
}

func TestServiceSearchItems(t *testing.T) {
	t.Parallel()

	svc := requireService(t)

	unlocked, locked, err := svc.SearchItems(uniqueAttributes("search"))
	require.NoError(t, err)
	require.Empty(t, unlocked)
	require.Empty(t, locked)
}

func TestServiceSearchItemsInvalidAttributesType(t *testing.T) {
	t.Parallel()

	svc := requireService(t)

	_, _, err := svc.SearchItems(secrets.Attributes{"": ""})
	require.NoError(t, err)
}

func TestServiceLockUnlock(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	collection := requireAnyCollection(t, svc)

	err := svc.Unlock([]secrets.LockableObject{collection})
	require.NoError(t, err)

	err = svc.Lock([]secrets.LockableObject{collection})
	require.NoError(t, err)

	err = svc.Unlock([]secrets.LockableObject{collection})
	require.NoError(t, err)
}

func TestServiceGetSecretsEmptyItems(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	session := requireOpenSession(t, svc)
	defer func() {
		_ = session.Close()
	}()

	secretsMap, err := svc.GetSecrets(nil, session)
	require.NoError(t, err)
	require.NotNil(t, secretsMap)
}

func TestServiceReadAlias(t *testing.T) {
	t.Parallel()

	svc := requireService(t)

	collection, err := svc.ReadAlias("default")
	require.NoError(t, err)
	require.NotNil(t, collection)
	require.NotEmpty(t, collection.Path())
}

func TestServiceSetAlias(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	collection := requireAnyCollection(t, svc)

	// Use existing collection to avoid creating extra state.
	err := svc.SetAlias(*collection, "default")
	require.NoError(t, err)
}

func TestServiceUnlockEmpty(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	require.NoError(t, svc.Unlock(nil))
}

func TestServiceUnlockInvalidPath(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	err := svc.Unlock(
		[]secrets.LockableObject{fakeLockable{path: dbus.ObjectPath(string(secrets.ServicePath()) + "/invalid")}},
	)
	require.NoError(t, err)
}

func TestServiceUnlockMalformedPathReturnsError(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	err := svc.Unlock([]secrets.LockableObject{fakeLockable{path: dbus.ObjectPath("not-an-object-path")}})
	require.Error(t, err)
}

func TestServiceLockInvalidPath(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	err := svc.Lock(
		[]secrets.LockableObject{fakeLockable{path: dbus.ObjectPath(string(secrets.ServicePath()) + "/invalid")}},
	)
	require.NoError(t, err)
}

func TestServiceLockMalformedPathReturnsError(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	err := svc.Lock([]secrets.LockableObject{fakeLockable{path: dbus.ObjectPath("not-an-object-path")}})
	require.Error(t, err)
}

func TestServiceGetSecretsInvalidSessionReturnsError(t *testing.T) {
	t.Parallel()

	svc := requireService(t)

	conn, err := dbus.ConnectSessionBus()
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	invalidSession := secrets.NewSession(conn, secrets.SessionPath("does-not-exist"))

	_, err = svc.GetSecrets(nil, invalidSession)
	require.Error(t, err)
}

func TestServiceReadAliasUnknownAlias(t *testing.T) {
	t.Parallel()

	svc := requireService(t)

	collection, err := svc.ReadAlias("definitely-not-a-real-alias")
	require.NoError(t, err)
	require.NotNil(t, collection)
}

func TestServiceSetAliasInvalidPathReturnsError(t *testing.T) {
	t.Parallel()

	svc := requireService(t)
	bad := secrets.NewCollection(svc, secrets.CollectionPath("invalid"))

	err := svc.SetAlias(*bad, "default")
	require.Error(t, err)
}
